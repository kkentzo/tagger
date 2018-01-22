package main

import (
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Path    string
	Indexer *Indexer
	Watcher *Watcher
}

func DefaultProject(indexer *Indexer) *Project {
	return &Project{
		Path:    ".",
		Indexer: indexer,
		Watcher: &Watcher{
			Root:       ".",
			Exclusions: NewPathSet(indexer.Exclude),
		},
	}
}

func (project *Project) Monitor() {
	indexEvents := make(chan struct{})
	go project.Watcher.Watch(indexEvents)
	// perform an initial indexing
	project.Index()
	for range indexEvents {
		go project.Index()
	}
}

func (project *Project) Index() {
	project.Indexer.Index(project.Watcher.Root)
	log.Info("Indexing ", project.Path)
}
