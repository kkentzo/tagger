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
	Exclusions   *PathSet
	maxFrequency time.Duration
}

func (rw *Watcher) Watch(indexEvents chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer watcher.Close()

	defer close(indexEvents) // TODO: initialize Trigger in method (somehow)??

	// add project files
	add(rw.Root, watcher, rw.Exclusions)

	log.Info("Watching ", rw.Root)
	// start monitoring
	mustReindex := false
	var idxMsg struct{}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if mustReindex {
				indexEvents <- idxMsg
				mustReindex = false
			}
		case event := <-watcher.Events:
			// TODO: make tags a parameter
			if filepath.Base(event.Name) == "TAGS" {
				continue
			}
			log.Debugf("Event %s on %s", event.Op, event.Name)
			if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				remove(event.Name, watcher) // this is non-recursive...
				mustReindex = true
			} else if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				fileInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Error(err.Error()) // stat error
					continue
				} else if fileInfo.IsDir() {
					add(event.Name, watcher, rw.Exclusions)
				}
				// TODO: Consider file type here??
				mustReindex = true
			} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			} else {
				mustReindex = true
			}
		case err := <-watcher.Errors:
			log.Error(err.Error())
		}
	}

}

func add(path string, watcher *fsnotify.Watcher, exclusions *PathSet) error {
	directories, err := discover(path, exclusions)
	if err != nil {
		return err
	}
	for _, file := range directories {
		err := watcher.Add(file)
		if err != nil {
			// TODO: This raises a "Too many files open" on MacOS
			log.Error(err.Error())
		}
		log.Debug("Adding", file)
	}
	return nil
}

func remove(path string, watcher *fsnotify.Watcher) error {
	watcher.Remove(path)
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
