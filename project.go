package main

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	watcher *RecursiveWatcher
	Indexer *Indexer
	Path    string
}

func DefaultProject(indexer *Indexer) *Project {
	project := &Project{Path: "."}
	return project
}

func (project *Project) Initialize() {
	project.watcher = NewRecursiveWatcher(project.Path, NewPathSet(project.Indexer.Exclude))
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
