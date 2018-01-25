package main

import "sync"

type Manager struct {
	indexer  *Indexer
	projects map[string]*Project
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
	path = ExpandHomeDir(path)
	if _, ok := manager.projects[path]; !ok {
		manager.projects[path] = &Project{
			Path:    path,
			Indexer: manager.indexer,
		}
	}
}

func (manager *Manager) Start() {
	var wg sync.WaitGroup
	for _, project := range manager.projects {
		go project.Monitor()
		wg.Add(1)
	}
	wg.Wait()
}
