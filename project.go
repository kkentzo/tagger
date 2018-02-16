package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Monitorable interface {
	Monitor(ctx context.Context)
	Index()
}

type Project struct {
	Path    string
	Indexer Indexable
	Watcher Watchable
	// TODO: Pass to indexer
	TagFile string
	// TODO: do we need this if multi-project is killed?
	Notify func()
}

func DefaultProject(indexer Indexable, watcher Watchable) *Project {
	return &Project{
		Path:    ".",
		Indexer: indexer,
		Watcher: watcher,
	}
}

func (project *Project) Monitor(ctx context.Context) {
	// perform an initial indexing
	go project.Index()
	defer project.Watcher.Close()
	wctx, cancel := context.WithCancel(ctx)
	go project.Watcher.Watch(wctx)
	for {
		select {
		case <-project.Watcher.Events():
			go project.Index()
			if project.Notify != nil {
				go project.Notify()
			}
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

func (project *Project) Index() {
	log.Info("Indexing ", project.Path)
	project.Indexer.Index(project.Path)
}
