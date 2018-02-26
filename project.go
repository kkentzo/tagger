package main

import (
	"context"

	"github.com/kkentzo/tagger/indexers"
	"github.com/kkentzo/tagger/watchers"
	log "github.com/sirupsen/logrus"
)

type Monitorable interface {
	Monitor(ctx context.Context)
	Index(bool)
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
	go project.Index(false)
	defer project.Watcher.Close()
	wctx, cancel := context.WithCancel(ctx)
	go project.Watcher.Watch(wctx)
	for {
		select {
		// TODO: Consume event (Special File)
		case e := <-project.Watcher.Events():
			go project.Index(e.IsSpecial)
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

func (project *Project) Index(isSpecial bool) {
	log.Info("Indexing ", project.Path)
	project.Indexer.Index(project.Path, isSpecial)
}
