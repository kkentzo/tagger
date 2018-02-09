package main

import (
	"fmt"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

type Indexable interface {
	Index(string)
	CreateWatcher(string) Watchable
}

type IndexerType string

const (
	Generic IndexerType = "generic"
	Rvm                 = "rvm"
)

type Indexer struct {
	Program      string
	Args         []string
	TagFile      string `yaml:"tag_file"`
	Type         IndexerType
	Exclude      []string
	MaxFrequency time.Duration `yaml:"max_frequency"`
}

func DefaultIndexer() *Indexer {
	return &Indexer{
		Program: "ctags",
		Args:    []string{"-R", "-e"},
		TagFile: "TAGS",
		Type:    Generic,
		Exclude: []string{".git"},
	}
}

func (indexer *Indexer) Index(root string) {
	args := indexer.GetArguments(root)
	cmd := exec.Command(indexer.Program, args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(out, err.Error())
	}
}

func (indexer *Indexer) CreateWatcher(root string) Watchable {
	return NewWatcher(root, indexer.Exclude, indexer.MaxFrequency)
}

func (indexer *Indexer) GetArguments(root string) []string {
	var args []string
	if indexer.TagFile != "" {
		args = append(args, fmt.Sprintf("-f %s", indexer.TagFile))
	}
	// add user-requested arguments
	args = append(args, indexer.Args...)
	// add excluded paths
	exclusions := []string{}
	for _, excl := range indexer.Exclude {
		exclusions = append(exclusions, fmt.Sprintf("--exclude=%s", excl))
	}
	args = append(args, exclusions...)
	// add paths to be indexed
	paths := []string{
		".",
	}
	// is this an RVM project?
	if indexer.Type == Rvm && isRvm(root) {
		if gemsetPath, err := rvmGemsetPath(root); err != nil {
			log.Error("Can not determine gemset path for rvm project at ", root)
		} else {
			paths = append(paths, gemsetPath)
			// TODO: should we append --languages=ruby here?
		}
	}
	args = append(args, paths...)
	return args
}
