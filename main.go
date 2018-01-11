package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	project := NewProject(
		"/Users/kkentzo/Workspace/agnostic_backend",
		&Indexer{
			command: "",
			args:    "",
		},
		[]string{".git", "coverage"})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go project.Monitor(ctx)
	// wait for interrupt signal
	s := <-c
	fmt.Println("Got signal:", s)
	// send the cancel signal to all projects
	cancel()

}
