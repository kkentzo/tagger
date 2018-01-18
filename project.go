package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	watcher *RecursiveWatcher
	indexer *Indexer
	Path    string
}

func (project *Project) Initialize(indexer *Indexer) {
	project.indexer = indexer
	project.watcher = NewRecursiveWatcher(project.Path, NewPathSet(indexer.Exclude))
}

func (project *Project) Monitor(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // TODO: Use cancel appropriately (how?)
	go project.watcher.Watch(ctx)
	// perform an initial indexing
	project.Index()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("CANCELLED!")
			return
		case <-project.watcher.trigger:
			go project.Index()
		}
	}
}

func (project *Project) Index() {
	project.indexer.Index(project.watcher.Root)
	log.Info("Reindexing!")
}
