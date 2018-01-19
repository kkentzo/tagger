package main

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	watcher *RecursiveWatcher
	indexer *Indexer
	Path    string
}

func DefaultProject(indexer *Indexer) *Project {
	project := &Project{Path: "."}
	project.Initialize(indexer)
	return project
}

func (project *Project) Initialize(indexer *Indexer) {
	project.indexer = indexer
	project.watcher = NewRecursiveWatcher(project.Path, NewPathSet(indexer.Exclude))
}

func (project *Project) Monitor() {
	go project.watcher.Watch()
	// perform an initial indexing
	project.Index()
	for range project.watcher.trigger {
		go project.Index()
	}
}

func (project *Project) Index() {
	project.indexer.Index(project.watcher.Root)
	log.Info("Reindexing!")
}
