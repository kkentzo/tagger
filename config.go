package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Indexer  Indexer
	Projects []*Project
}

func NewConfig() *Config {
	config := &Config{}
	// TODO: Make path (demo.yml) a flag // if not supplied read from ~/.tagger.yml
	contents, err := ioutil.ReadFile("demo.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	yaml.Unmarshal(contents, config)
	process(config)
	return config
}

func process(config *Config) {
	// process ~ (HOME)
	for _, project := range config.Projects {
		project.Path = substTilde(project.Path)
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
