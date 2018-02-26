package watchers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

var msg struct{}

type Watchable interface {
	Watch(context.Context)
	Events() chan Event
	Close()
}

type Event struct {
	IsSpecial bool
}

type Watcher struct {
	Root        string
	MaxPeriod   time.Duration
	SpecialFile string
	fsWatcher   FsWatchable
	events      chan Event
}

func NewWatcher(root string, exclusions []string, tagFilePrefix string, maxFrequency time.Duration) *Watcher {
	return &Watcher{
		Root:      root,
		MaxPeriod: maxFrequency,
		fsWatcher: NewFsWatcher(exclusions, tagFilePrefix),
		events:    make(chan Event),
	}
}

func (watcher *Watcher) Events() chan Event {
	return watcher.events
}

func (watcher *Watcher) Close() {
	close(watcher.events)
	watcher.fsWatcher.Close()
}

func (watcher *Watcher) Watch(ctx context.Context) {
	// add project files
	watcher.fsWatcher.Add(watcher.Root)

	log.Info("Watching ", watcher.Root)
	// start monitoring
	mustReindex := false
	isSpecial := false

	ticker := time.NewTicker(watcher.MaxPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if mustReindex {
				watcher.events <- Event{IsSpecial: isSpecial}
				mustReindex = false
				isSpecial = false
			}
		case event := <-watcher.fsWatcher.Events():
			mustReindex = mustReindex ||
				watcher.fsWatcher.Handle(event)
			if mustReindex && event.Name == watcher.SpecialFile {
				isSpecial = true
			}
		case err := <-watcher.fsWatcher.Errors():
			log.Error(err.Error())
		}
	}

}
