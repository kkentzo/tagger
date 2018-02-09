package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	Path    string
	Indexer Indexable
	Watcher Watcher
}

func DefaultProject(indexer Indexable, watcher Watcher) *Project {
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
