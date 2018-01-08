package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	project := &Project{
		root: "/tmp/foo",
	}

	project.Monitor(watcher)

}
