package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{})

	// parse command line args
	configFilePath := flag.String("c", Canonicalize("~/.tagger.yml"), "Path to config file")
	debug := flag.Bool("d", false, "Activate debug logging level")
	x := flag.Bool("x", false, "Generate tags in current directory and exit")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// parse config
	var config *Config
	if *x {
		config = &Config{
			Indexer:  DefaultGenericIndexer(),
			Projects: []struct{ Path string }{{"."}},
		}
	} else {
		config = NewConfig(*configFilePath)
	}

	// create project manager
	manager := NewManager(config)
	// create server
	server := &Server{Manager: manager, Port: config.Port}
	go server.Listen()
	// start monitoring
	manager.Start()
}
