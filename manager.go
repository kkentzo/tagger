package main

import (
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Manager struct {
	indexer  *Indexer
	projects map[string]*Project
	pg       sync.WaitGroup
}

func NewManager(config *Config) *Manager {
	manager := &Manager{
		indexer:  config.Indexer,
		projects: make(map[string]*Project),
	}
	for _, p := range config.Projects {
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
		}
		manager.projects[path] = project
		manager.pg.Add(1)
		// TODO: Create cancellation context and pass to Monitor
		go project.Monitor()
	}
}

func (manager *Manager) Remove(path string) {
	path = Canonicalize(path)
	// what happens if path does not exist?
	// This is legit in case the project root is deleted from the fs
	if _, ok := manager.projects[path]; ok {
		// TODO: Send cancellation signal to project
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
