package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Indexer  *Indexer
	Projects []struct{ Path string }
}

func NewConfig(configFilePath string) *Config {
	config := &Config{}
	// TODO: Make path (demo.yml) a flag // if not supplied read from ~/.tagger.yml
	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Config file not found: ", configFilePath)
	}
	yaml.Unmarshal(contents, config)
	return config
}

func prepareProjects(config *Config) {
	// process ~ (HOME)
	for _, project := range config.Projects {
		project.Path = ExpandHomeDir(project.Path)
		// TODO: Insert Indexer settings into project??
	}
	// TODO: check for ctags binary
}
