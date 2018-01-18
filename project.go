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
	Libs    []string
	// TODO: add exclusions here
	// TODO: add file types (a regex??) (inclusions)
}

func (project *Project) Initialize(config *Config) {
	project.indexer = &config.Indexer
	project.watcher = NewRecursiveWatcher(project.Path, NewPathSet(config.Indexer.Exclude))
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
