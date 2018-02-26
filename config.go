package main

import (
	"io/ioutil"

	"github.com/kkentzo/tagger/indexers"
	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Port     int
	Indexer  *indexers.Indexer
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
