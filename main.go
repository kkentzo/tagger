package main

import (
	"flag"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {

	// parse command line args
	configFilePath := flag.String("c", substTilde("~/.tagger.yml"), "Path to config file")
	debug := flag.Bool("d", false, "Activate debug logging level")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	config := NewConfig(*configFilePath)

	var wg sync.WaitGroup
	for _, project := range config.Projects {
		project.Initialize(&config.Indexer)
		go project.Monitor()
		wg.Add(1)
	}
	wg.Wait()
}
