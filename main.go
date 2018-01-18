package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

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
	fmt.Println(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// TODO: do this for all projects
	ctx, cancel := context.WithCancel(context.Background())
	config.Projects[0].Initialize(&config.Indexer)
	go config.Projects[0].Monitor(ctx)
	// wait for interrupt signal
	s := <-c
	fmt.Println("Got signal:", s)
	// send the cancel signal to all projects
	cancel()

}
