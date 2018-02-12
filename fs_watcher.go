package main

// encapsulates fsnotify.Watcher

import (
	"github.com/fsnotify/fsnotify"

	log "github.com/sirupsen/logrus"
)

type FsWatchable interface {
	Add(string) error
	Remove(string) error
	Events() chan fsnotify.Event
	Errors() chan error
	Close() error
}

type FsWatcher struct {
	watcher *fsnotify.Watcher
}

func NewFsWatcher() *FsWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Failed to initialize filesystem watcher")
	}
	return &FsWatcher{watcher: w}
}

func (watcher *FsWatcher) Add(path string) error {
	return watcher.Add(path)
}

func (watcher *FsWatcher) Remove(path string) error {
	return watcher.Remove(path)
}

func (watcher *FsWatcher) Close() error {
	return watcher.Close()
}

func (watcher *FsWatcher) Events() chan fsnotify.Event {
	return watcher.watcher.Events
}

func (watcher *FsWatcher) Errors() chan error {
	return watcher.watcher.Errors
}
