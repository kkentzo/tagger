package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	watcher *RecursiveWatcher
	indexer *Indexer
	program string
	args    string
	// TODO: add file types (a regex??) (inclusions)
}

func NewProject(root string, indexer *Indexer, exclude []string) *Project {
	return &Project{
		watcher: NewRecursiveWatcher(root, NewPathSet(exclude)),
		indexer: indexer,
	}
}

func (project *Project) Monitor(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // TODO: Use cancel appropriately (how?)
	go project.watcher.Watch(ctx)
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
