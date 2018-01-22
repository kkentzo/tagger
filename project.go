package main

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Path    string
	Indexer *Indexer
	Watcher *RecursiveWatcher
}

func DefaultProject(indexer *Indexer) *Project {
	return &Project{
		Path:    ".",
		Indexer: indexer,
		Watcher: NewRecursiveWatcher(".", NewPathSet(indexer.Exclude)),
	}
}

func (project *Project) Monitor() {
	go project.Watcher.Watch()
	// perform an initial indexing
	project.Index()
	for range project.Watcher.trigger {
		go project.Index()
	}
}

func (project *Project) Index() {
	project.Indexer.Index(project.Watcher.Root)
	log.Info("Indexing ", project.Path)
}
