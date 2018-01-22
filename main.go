package main

import (
	"flag"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ExpandHomeDir(path string) string {
	if strings.Contains(path, "~") {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}

func main() {

	log.SetFormatter(&log.TextFormatter{})

	// parse command line args
	configFilePath := flag.String("c", ExpandHomeDir("~/.tagger.yml"), "Path to config file")
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
			Indexer:  DefaultIndexer(),
			Projects: []struct{ Path string }{{"."}},
		}
	} else {
		config = NewConfig(*configFilePath)
	}

	// create projects
	manager := &Manager{}
	for _, p := range config.Projects {
		path := ExpandHomeDir(p.Path)
		project := &Project{
			Path:    path,
			Indexer: config.Indexer,
			Watcher: &Watcher{
				Root:       path,
				Exclusions: NewPathSet(config.Indexer.Exclude),
			},
		}
		manager.Add(project)
	}
	manager.Start()
}
