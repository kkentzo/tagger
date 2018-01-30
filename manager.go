package main

import (
	"fmt"
	"net/http"
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
	path = ExpandHomeDir(path)
	// skip non-existent path
	if !fileExists(path) {
		return
	}
	if _, ok := manager.projects[path]; !ok {
		// TODO: Create cancellation context
		manager.projects[path] = &Project{
			Path:    path,
			Indexer: manager.indexer,
		}
		manager.pg.Add(1)
	}
}

func (manager *Manager) Remove(path string) {
	// TODO: what happens if path does not exist? => 404
	// This is legit in case the project root is deleted from the fs
	if _, ok := manager.projects[path]; ok {
		// TODO: Send cancellation signal to project
		// remove project from registry
		delete(manager.projects, path)
		manager.pg.Done()
	}
}

func (manager *Manager) Start() {
	for _, project := range manager.projects {
		go project.Monitor()
	}
	manager.pg.Wait()
}

func (manager *Manager) Listen(port int) {
	// register handlers
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, manager)
	})
	// launch server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request, m *Manager) {
	switch r.Method {
	case "GET":
		fmt.Println(r.Method)
	case "POST":
		fmt.Println(r.Method)
	case "DELETE":
		fmt.Println(r.Method)
	default:
		fmt.Println(r.Method)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ExpandHomeDir(path string) string {
	if strings.Contains(path, "~") {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}
