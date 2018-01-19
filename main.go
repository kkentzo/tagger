package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{})

	// parse command line args
	configFilePath := flag.String("c", substTilde("~/.tagger.yml"), "Path to config file")
	debug := flag.Bool("d", false, "Activate debug logging level")
	x := flag.Bool("x", false, "Generate tags in current directory and exit")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var config *Config
	if *x {
		indexer := DefaultIndexer()
		config = &Config{
			Indexer:  indexer,
			Projects: []*Project{DefaultProject(indexer)},
		}
	} else {
		config = NewConfig(*configFilePath)
	}

	manager := &Manager{Projects: config.Projects}
	manager.Start()
}
