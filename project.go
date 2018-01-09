package main

import log "github.com/sirupsen/logrus"

type Project struct {
	watcher *RecursiveWatcher
	args    string // args to ctags binary
	// TODO: add file types (a regex??) (inclusions)
}

func NewProject(root string, args string, exclude []string) *Project {
	return &Project{
		watcher: NewRecursiveWatcher(root, NewPathSet(exclude)),
		args:    args,
	}
}

func (project *Project) Monitor() {
	go project.watcher.Watch()
	// TODO: listen on events from watcher to trigger indexing
	for range project.watcher.trigger {
		project.Reindex()
	}

}

func (project *Project) Reindex() {
	log.Printf("Reindexing")
}
