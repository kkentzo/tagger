package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	root    string   // project filesystem root (absolute)
	args    string   // args to ctags binary
	include []string // directories within root to monitor
	exclude []string // directories to exclude from tagging
	// TODO: add file types (a regex??)
}

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
			// check event type: Create, Write, Remove, Rename
			log.Info(event)
			// if event type: remove / rename
			// ==> if directory: remove from project.include
			// if event type: create / write
			// ==> if directory: add to project.include
			fi, err := os.Stat(event.Name)
			if err != nil {
				continue
			}
			if fi.IsDir() {
				// add the directory to project files
				fmt.Println("[directory]")
			} else {
				fmt.Println("[file]")
			}

		case err = <-watcher.Errors:
			log.Error("error:", err)
		}

	}
}

// return a slice with all project directories (recursive)
func (project *Project) files() []string {
	return []string{project.root}
}
