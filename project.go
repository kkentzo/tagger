package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	Path    string
	Indexer *Indexer
}

func DefaultProject(indexer *Indexer) *Project {
	return &Project{
		Path:    ".",
		Indexer: indexer,
	}
}

func (project *Project) Monitor(ctx context.Context) {
	// perform an initial indexing
	go project.Index()

	watcher := NewWatcher(project.Path, project.Indexer.Exclude, project.Indexer.MaxFrequency)
	defer watcher.Close()
	wctx, cancel := context.WithCancel(ctx)
	go watcher.Watch(wctx)
	for {
		select {
		case <-watcher.Events():
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
