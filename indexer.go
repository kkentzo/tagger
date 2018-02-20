package main

import (
	"fmt"
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
	ExcludeDirs  []string
	MaxFrequency time.Duration `yaml:"max_frequency"`
}

func DefaultIndexer() *Indexer {
	return &Indexer{
		Program:     "ctags",
		Args:        []string{"-R", "-e"},
		TagFile:     TagFilePrefix,
		Type:        Generic,
		ExcludeDirs: []string{".git"},
	}
}

func (indexer *Indexer) Index(root string) {
	args := indexer.GetArguments(root)
	out, err := ExecInPath(indexer.Program, args, root)
	if err != nil {
		log.Error(out, err.Error())
	}
}

func (indexer *Indexer) CreateWatcher(root string) Watchable {
	w := NewWatcher(root, indexer.ExcludeDirs, indexer.MaxFrequency)
	if indexer.Type == Rvm {
		w.SpecialFile = "Gemfile.lock"
	}
	return w
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
	for _, excl := range indexer.ExcludeDirs {
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
