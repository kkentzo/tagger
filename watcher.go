package main

import (
	"context"
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

type FsEventHandlerFunc func(FsWatchable, fsnotify.Event) bool

type Watcher struct {
	Root         string
	MaxFrequency time.Duration
	HandlerFunc  FsEventHandlerFunc
	fsWatcher    FsWatchable
	events       chan struct{}
}

func NewWatcher(root string, exclusions []string, maxFrequency time.Duration) *Watcher {
	return &Watcher{
		Root:         root,
		MaxFrequency: maxFrequency,
		HandlerFunc:  handle,
		fsWatcher:    NewFsWatcher(exclusions),
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
	// add project files
	// TODO: handle error
	watcher.fsWatcher.Add(watcher.Root)

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
		case event := <-watcher.fsWatcher.Events():
			// TODO: make TAGS a parameter
			if filepath.Base(event.Name) == "TAGS" {
				continue
			}
			mustReindex = mustReindex ||
				watcher.HandlerFunc(watcher.fsWatcher, event)
		case err := <-watcher.fsWatcher.Errors():
			log.Error(err.Error())
		}
	}

}

func handle(fsWatcher FsWatchable, event fsnotify.Event) bool {
	return fsWatcher.Handle(event)
}
