package main

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Root         string
	Exclusions   []string
	MaxFrequency time.Duration
	fsWatcher    *fsnotify.Watcher
}

func NewWatcher(root string, exclusions []string, maxFrequency time.Duration) *Watcher {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Watcher{
		Root:         root,
		Exclusions:   exclusions,
		MaxFrequency: maxFrequency,
		fsWatcher:    fsWatcher,
	}
}

func (watcher *Watcher) Watch(indexEvents chan<- struct{}) {
	defer watcher.fsWatcher.Close() // TODO: Is this appropriate here?
	//defer close(indexEvents) => FIX: should we? this makes the test pass!!!

	// create set of excluded stuff
	exclusions := NewPathSet(watcher.Exclusions)

	// add project files
	add(watcher.Root, watcher.fsWatcher, exclusions)

	log.Info("Watching ", watcher.Root)
	// start monitoring
	mustReindex := false
	var idxMsg struct{}

	ticker := time.NewTicker(watcher.MaxFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if mustReindex {
				indexEvents <- idxMsg
				mustReindex = false
			}
		case event := <-watcher.fsWatcher.Events:
			// TODO: make tags a parameter
			if filepath.Base(event.Name) == "TAGS" {
				continue
			}
			log.Debugf("Event %s on %s", event.Op, event.Name)
			if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				remove(event.Name, watcher.fsWatcher) // this is non-recursive...
				mustReindex = true
			} else if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				fileInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Error(err.Error()) // stat error
					continue
				} else if fileInfo.IsDir() {
					add(event.Name, watcher.fsWatcher, exclusions)
				}
				// TODO: Consider file type here??
				mustReindex = true
			} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			} else {
				mustReindex = true
			}
		case err := <-watcher.fsWatcher.Errors:
			log.Error(err.Error())
		}
	}

}

func add(path string, fsWatcher *fsnotify.Watcher, exclusions *PathSet) error {
	directories, err := discover(path, exclusions)
	if err != nil {
		return err
	}
	for _, file := range directories {
		err := fsWatcher.Add(file)
		if err != nil {
			// TODO: This raises a "Too many files open" on MacOS
			log.Error(err.Error())
		}
		log.Debug("Adding", file)
	}
	return nil
}

func remove(path string, fsWatcher *fsnotify.Watcher) error {
	fsWatcher.Remove(path)
	log.Info("Removing", path)
	return nil
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
