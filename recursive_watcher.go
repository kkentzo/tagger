package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	Root       string
	Exclusions *PathSet
	trigger    chan struct{}
}

func NewRecursiveWatcher(root string, exclusions *PathSet) *RecursiveWatcher {
	return &RecursiveWatcher{
		Root:       root,
		Exclusions: exclusions,
		trigger:    make(chan struct{}),
	}
}

func (rw *RecursiveWatcher) Trigger() <-chan struct{} {
	return rw.trigger
}

func (rw *RecursiveWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer watcher.Close()

	// add project files
	add(rw.Root, watcher, rw.Exclusions)
	// start monitoring
	// TODO: we need a way out of this for loop
	// TODO: when to trigger re-indexing??
	//var msg struct{}
	for {
		select {
		case event := <-watcher.Events:
			log.Printf("Event %s on %s", event.Op, event.Name)
			if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				remove(event.Name, watcher) // this is non-recursive...
			} else if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				fileInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Error(err.Error())
				} else if fileInfo.IsDir() {
					add(event.Name, watcher, rw.Exclusions)
				}
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
			log.Error(err.Error())
		}
		log.Info("Adding", file)
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
