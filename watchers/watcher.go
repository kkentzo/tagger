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

type Watcher struct {
	Root      string
	MaxPeriod time.Duration
	fsWatcher FsWatchable
	events    chan Event
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

	ticker := time.NewTicker(watcher.MaxPeriod)
	defer ticker.Stop()

	// TODO: Change this to pointer
	event := NewEvent()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if mustReindex {
				watcher.events <- event
				mustReindex = false
				event = NewEvent()
			}
		case fsEvent := <-watcher.fsWatcher.Events():
			shouldReindex := watcher.fsWatcher.Handle(fsEvent)
			if shouldReindex {
				event.Names.Add(fsEvent.Name)
			}
			mustReindex = mustReindex || shouldReindex
		case err := <-watcher.fsWatcher.Errors():
			log.Error(err.Error())
		}
	}

}
