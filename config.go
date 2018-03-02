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
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", configFilePath, err.Error())
	}
	return config
}
