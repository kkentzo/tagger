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
	// perform an initial indexing
	go project.Index()

	watcher := NewWatcher(project.Path, project.Indexer.Exclude, project.Indexer.MaxFrequency)
	go watcher.Watch()
	for range watcher.Events() {
		go project.Index()
	}
	watcher.Close()
}

func (project *Project) Index() {
	log.Info("Indexing ", project.Path)
	project.Indexer.Index(project.Path)
}
