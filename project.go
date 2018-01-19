package main

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	watcher *RecursiveWatcher
	Indexer *Indexer
	Path    string
}

func NewProject(path string, indexer *Indexer) *Project {
	return &Project{
		Path:    path,
		Indexer: indexer,
		watcher: NewRecursiveWatcher(path, NewPathSet(indexer.Exclude)),
	}
}

func DefaultProject(indexer *Indexer) *Project {
	return NewProject(".", indexer)
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
	project.Indexer.Index(project.watcher.Root)
	log.Info("Indexing ", project.Path)
}
