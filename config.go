package main

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Indexer  *Indexer
	Projects []*Project
}

func NewConfig(configFilePath string) *Config {
	config := &Config{}
	// TODO: Make path (demo.yml) a flag // if not supplied read from ~/.tagger.yml
	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Config file not found: ", configFilePath)
	}
	yaml.Unmarshal(contents, config)
	prepareProjects(config)
	return config
}

func prepareProjects(config *Config) {
	// process ~ (HOME)
	for _, project := range config.Projects {
		project.Path = substTilde(project.Path)
		project.Indexer = config.Indexer
		project.Initialize()
		// TODO: Insert Indexer settings into project??
	}
	// TODO: check for ctags binary
}

func substTilde(path string) string {
	if strings.Contains(path, "~") {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}
