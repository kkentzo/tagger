package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Project struct {
	root    string   // project filesystem root (absolute)
	args    string   // args to ctags binary
	exclude []string // directories to exclude from tagging
}

// need to recursively figure out all project directories

func (project *Project) Monitor(watcher *fsnotify.Watcher) {
	var err error

	// add all project files to the watcher
	for _, file := range project.files() {
		err = watcher.Add(file)
		if err != nil {
			log.Println("error:", err)
		}
	}

	// do the monitoring
	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err = <-watcher.Errors:
			log.Println("error:", err)
		}

	}
}

// return a slice with all project directories (recursive)
func (project *Project) files() []string {
	return []string{project.root}
}
