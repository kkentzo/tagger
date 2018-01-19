package main

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	Root         string
	Exclusions   *PathSet
	maxFrequency time.Duration
	trigger      chan struct{}
}

func NewRecursiveWatcher(root string, exclusions *PathSet) *RecursiveWatcher {
	return &RecursiveWatcher{
		Root:       root,
		Exclusions: exclusions,
		trigger:    make(chan struct{}),
	}
}

func (rw *RecursiveWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer watcher.Close()

	defer close(rw.trigger) // TODO: initialize trigger in method (somehow)??

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
				rw.trigger <- idxMsg
				mustReindex = false
			}
		case event := <-watcher.Events:
			// TODO: make tags a parameter
			if filepath.Base(event.Name) == "TAGS" {
				continue
			}
			log.Debug("Event %s on %s", event.Op, event.Name)
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			} else if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				remove(event.Name, watcher) // this is non-recursive...
				mustReindex = true
			} else if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				fileInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Error(err.Error()) // stat error
				} else if fileInfo.IsDir() {
					add(event.Name, watcher, rw.Exclusions)
					mustReindex = true
				}
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
