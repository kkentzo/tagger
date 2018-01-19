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

	var config *Config
	if *x {
		config = &Config{
			Indexer:  DefaultIndexer(),
			Projects: []struct{ Path string }{{"."}},
		}
	} else {
		config = NewConfig(*configFilePath)
	}

	projects := make([]*Project, len(config.Projects))
	for idx, p := range config.Projects {
		projects[idx] = NewProject(p.Path, config.Indexer)
	}

	manager := &Manager{Projects: projects}
	manager.Start()
}
