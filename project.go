package main

import (
	"context"

	"github.com/kkentzo/tagger/indexers"
	"github.com/kkentzo/tagger/watchers"
	log "github.com/sirupsen/logrus"
)

type Monitorable interface {
	Monitor(context.Context)
	// TODO: change arg to pointer
	Index(watchers.Event)
}

type Project struct {
	Path    string
	Indexer indexers.Indexable
	Watcher watchers.Watchable
}

func DefaultProject(indexer indexers.Indexable, watcher watchers.Watchable) *Project {
	return &Project{
		Path:    ".",
		Indexer: indexer,
		Watcher: watcher,
	}
}

func (project *Project) Monitor(ctx context.Context) {
	// perform an initial indexing
	go project.Index(watchers.Event{})
	defer project.Watcher.Close()
	wctx, cancel := context.WithCancel(ctx)
	go project.Watcher.Watch(wctx)
	for {
		select {
		case e := <-project.Watcher.Events():
			// TODO: is this indexing goroutine thread-safe here?
			go project.Index(e)
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

func (project *Project) Index(event watchers.Event) {
	log.Info("Indexing ", project.Path)
	project.Indexer.Index(project.Path, event)
}
