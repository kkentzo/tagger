package main

import (
	"flag"

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

	// TODO: do this for all projects
	config.Projects[0].Initialize(&config.Indexer)
	config.Projects[0].Monitor()
}
