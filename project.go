package main

import (
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

func (project *Project) Monitor() {
	watcher := NewWatcher(project.Path, project.Indexer.Exclude, project.Indexer.MaxFrequency)
	indexEvents := make(chan struct{})
	go watcher.Watch(indexEvents)
	// perform an initial indexing
	project.Index()
	for range indexEvents {
		go project.Index()
	}
}

func (project *Project) Index() {
	log.Info("Indexing ", project.Path)
	project.Indexer.Index(project.Path)
}
