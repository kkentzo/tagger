package main

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var DefaultConfigFilePath string = filepath.Join(os.Getenv("HOME"), ".tagger.yml")

func main() {

	log.SetFormatter(&log.TextFormatter{})

	// parse command line args
	configFilePath := flag.String("c", DefaultConfigFilePath, "Path to config file")
	debug := flag.Bool("d", false, "Activate debug logging level")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// parse config
	config := NewConfig(*configFilePath)
	// create project manager
	manager := NewManager(config.Indexer, config.Projects)
	// create server
	server := &Server{Manager: manager, Port: config.Port}
	go server.Listen()
	// start monitoring
	manager.Start()
}
