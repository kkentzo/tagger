package main

import "sync"

type Manager struct {
	Projects []*Project
}

func (manager *Manager) Add(project *Project) {
	manager.Projects = append(manager.Projects, project)
}

func (manager *Manager) Start() {
	var wg sync.WaitGroup
	for _, project := range manager.Projects {
		go project.Monitor()
		wg.Add(1)
	}
	wg.Wait()
}
