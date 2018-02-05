package main

import (
	"context"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ProjectWithContext struct {
	Project *Project
	Cancel  context.CancelFunc
}

type Manager struct {
	indexer  Indexer
	projects map[string]*ProjectWithContext
	pg       sync.WaitGroup
}

func NewManager(indexer Indexer, projects []struct{ Path string }) *Manager {
	manager := &Manager{
		indexer:  indexer,
		projects: make(map[string]*ProjectWithContext),
	}
	for _, p := range projects {
		manager.Add(p.Path)
	}
	return manager
}

func (manager *Manager) Add(path string) {
	path = Canonicalize(path)
	// skip non-existent path
	if !fileExists(path) {
		log.Debugf("Path %s does not exist in filesystem", path)
		return
	}
	if manager.Exists(path) {
		log.Debugf("Path %s already monitored", path)
		return
	}
	if _, ok := manager.projects[path]; !ok {
		project := &Project{
			Path:    path,
			Indexer: manager.indexer,
			Watcher: manager.indexer.CreateWatcher(path),
		}
		ctx, cancel := context.WithCancel(context.Background())
		manager.projects[path] = &ProjectWithContext{
			Project: project,
			Cancel:  cancel,
		}
		manager.pg.Add(1)
		go project.Monitor(ctx)
	}
}

func (manager *Manager) Remove(path string) {
	path = Canonicalize(path)
	// what happens if path does not exist?
	// This is legit in case the project root is deleted from the fs
	if project, ok := manager.projects[path]; ok {
		// Send cancellation signal to project
		project.Cancel()
		// remove project from registry
		delete(manager.projects, path)
		manager.pg.Done()
	}
}

func (manager *Manager) Exists(path string) bool {
	_, ok := manager.projects[Canonicalize(path)]
	return ok
}

func (manager *Manager) Start() {
	manager.pg.Wait()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Canonicalize(path string) string {
	if strings.Contains(path, "~") {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}
