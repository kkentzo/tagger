package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

var msg struct{}

type Watchable interface {
	Watch(context.Context)
	Events() chan struct{}
	Close()
}

type FsEventHandlerFunc func(fsnotify.Event, *fsnotify.Watcher, *PathSet) bool

type Watcher struct {
	Root         string
	Exclusions   []string
	MaxFrequency time.Duration
	HandlerFunc  FsEventHandlerFunc
	fsWatcher    *fsnotify.Watcher
	events       chan struct{}
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
		HandlerFunc:  handle,
		fsWatcher:    fsWatcher,
		events:       make(chan struct{}),
	}
}

func (watcher *Watcher) Events() chan struct{} {
	return watcher.events
}

func (watcher *Watcher) Close() {
	close(watcher.events)
	watcher.fsWatcher.Close()
}

func (watcher *Watcher) Watch(ctx context.Context) {
	// create set of excluded stuff
	exclusions := NewPathSet(watcher.Exclusions)

	// add project files
	add(watcher.Root, watcher.fsWatcher, exclusions)

	log.Info("Watching ", watcher.Root)
	// start monitoring
	mustReindex := false

	ticker := time.NewTicker(watcher.MaxFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if mustReindex {
				watcher.events <- msg
				mustReindex = false
			}
		case event := <-watcher.fsWatcher.Events:
			// TODO: make TAGS a parameter
			if filepath.Base(event.Name) == "TAGS" {
				continue
			}
			mustReindex = mustReindex ||
				watcher.HandlerFunc(event, watcher.fsWatcher, exclusions)
		case err := <-watcher.fsWatcher.Errors:
			log.Error(err.Error())
		}
	}

}

func handle(event fsnotify.Event, fsWatcher *fsnotify.Watcher, excl *PathSet) bool {
	log.Debugf("Event %s on %s", event.Op, event.Name)
	if event.Op&fsnotify.Remove == fsnotify.Remove ||
		event.Op&fsnotify.Rename == fsnotify.Rename {
		remove(event.Name, fsWatcher) // this is non-recursive...
		return true
	} else if event.Op&fsnotify.Create == fsnotify.Create ||
		event.Op&fsnotify.Write == fsnotify.Write {
		if isDir, err := isDirectory(event.Name); err != nil {
			log.Error(err.Error())
			return false
		} else if isDir {
			add(event.Name, fsWatcher, excl)
		}
		return true
	} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return false
	} else {
		return true
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

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
