package watchers

// encapsulates fsnotify.Watcher

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/kkentzo/tagger/utils"

	log "github.com/sirupsen/logrus"
)

type FsWatchable interface {
	Handle(fsnotify.Event) bool
	Add(string) error
	Remove(string) error
	Events() chan fsnotify.Event
	Errors() chan error
	Close() error
}

type FsWatcher struct {
	*fsnotify.Watcher
	exclusions    *PathSet
	tagFilePrefix string
}

func NewFsWatcher(exclusions []string, tagFilePrefix string) *FsWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Failed to initialize filesystem watcher")
	}
	return &FsWatcher{
		Watcher:       w,
		exclusions:    NewPathSet(exclusions),
		tagFilePrefix: tagFilePrefix,
	}
}

func (watcher *FsWatcher) Handle(event fsnotify.Event) bool {
	if strings.Contains(filepath.Base(event.Name), watcher.tagFilePrefix) {
		return false
	}
	log.Debugf("Event %s on %s", event.Op, event.Name)
	if event.Op&fsnotify.Remove == fsnotify.Remove ||
		event.Op&fsnotify.Rename == fsnotify.Rename {
		watcher.Remove(event.Name) // this is non-recursive...
		return true
	} else if event.Op&fsnotify.Create == fsnotify.Create ||
		event.Op&fsnotify.Write == fsnotify.Write {
		if isDir, err := utils.IsDirectory(event.Name); err != nil {
			log.Error(err.Error())
			return false
		} else if isDir {
			err = watcher.Add(event.Name)
			if err != nil {
				// TODO: This raises a "Too many files open" on MacOS
				log.Error(err.Error())
				return false
			}
		}
		return true
	} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return false
	} else {
		return true
	}
}

func (watcher *FsWatcher) Add(path string) error {
	directories, err := discover(path, watcher.exclusions)
	if err != nil {
		return err
	}
	for _, file := range directories {
		err := watcher.Watcher.Add(file)
		if err != nil {
			log.Error(err.Error())
		}
		log.Debug("Adding", file)
	}
	return nil
}

func (watcher *FsWatcher) Events() chan fsnotify.Event {
	return watcher.Watcher.Events
}

func (watcher *FsWatcher) Errors() chan error {
	return watcher.Watcher.Errors
}

// return a slice with all directories under root but the excluded ones
func discover(root string, exclusions *PathSet) ([]string, error) {
	var directories []string
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if exclusions.Has(info.Name()) {
					return filepath.SkipDir
				} else {
					directories = append(directories, path)
					return nil
				}
			}
			return nil
		})
	return directories, err
}
