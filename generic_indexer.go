package main

import (
	"fmt"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

type GenericIndexer struct {
	Program      string
	Args         []string
	TagFile      string `yaml:"tag_file"`
	Type         IndexerType
	Exclude      []string
	MaxFrequency time.Duration `yaml:"max_frequency"`
}

func DefaultGenericIndexer() *GenericIndexer {
	return &GenericIndexer{
		Program: "ctags",
		Args:    []string{"-R", "-e"},
		TagFile: "TAGS",
		Type:    Generic,
		Exclude: []string{".git"},
	}
}

func (indexer *GenericIndexer) CreateProjectIndexer(other *GenericIndexer) Indexer {
	// merge indexer into other and return the other
	return nil
}

func (indexer *GenericIndexer) Index(root string) {
	args := indexer.GetArguments(root)
	cmd := exec.Command(indexer.Program, args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(out, err.Error())
	}
}

func (indexer *GenericIndexer) CreateWatcher(root string) *Watcher {
	return NewWatcher(root, indexer.Exclude, indexer.MaxFrequency)
}

func (indexer *GenericIndexer) GetArguments(root string) []string {
	args := []string{fmt.Sprintf("-f %s", indexer.TagFile)}
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
		paths = append(paths, rvmGemsetPath(root))
	}
	args = append(args, paths...)
	return args
}
