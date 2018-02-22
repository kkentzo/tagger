package main

import (
	"fmt"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type Indexable interface {
	Index(string, bool)
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
	ExcludeDirs  []string      `yaml:"exclude"`
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

func (indexer *Indexer) Index(root string, isSpecial bool) {
	indexer.indexProject(root)
	if isSpecial || !indexer.GemsetTagFileExists(root) {
		indexer.indexGemset(root)
	}
	tagFiles := []string{
		indexer.GetTagFileNameForProject(root),
		indexer.GetTagFileNameForGemset(root),
	}
	// TODO: This should be robust if one of the files does not exist
	err := ConcatFiles(filepath.Join(root, indexer.TagFile), tagFiles, root)
	if err != nil {
		log.Error("concat", tagFiles, err.Error())
	}
}

func (indexer *Indexer) CreateWatcher(root string) Watchable {
	w := NewWatcher(root, indexer.ExcludeDirs, indexer.MaxFrequency)
	if indexer.Type == Rvm {
		w.SpecialFile = "Gemfile.lock"
	}
	return w
}

func (indexer *Indexer) indexProject(root string) {
	args := indexer.GetProjectArguments(root)
	out, err := ExecInPath(indexer.Program, args, root)
	if err != nil {
		log.Error(out, err.Error())
	}
}

func (indexer *Indexer) indexGemset(root string) {
	if indexer.Type == Rvm && isRvm(root) {
		args := indexer.GetGemsetArguments(root)
		if len(args) == 0 {
			return
		}
		out, err := ExecInPath(indexer.Program, args, root)
		if err != nil {
			log.Error(out, err.Error())
		}
	}
}

func (indexer *Indexer) GetProjectArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s.project", indexer.TagFile))
	args = append(args, ".")
	return args
}

func (indexer *Indexer) GetGemsetArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s.gemset", indexer.TagFile))
	if gemsetPath, err := rvmGemsetPath(root); err != nil {
		log.Error("Can not determine gemset path for rvm project at ", root)
		return []string{}
	} else {
		args = append(args, gemsetPath)
	}
	return args
}

func (indexer *Indexer) GetGenericArguments(root string) []string {
	var args []string
	// add user-requested arguments
	args = append(args, indexer.Args...)
	// add excluded paths
	exclusions := []string{}
	for _, excl := range indexer.ExcludeDirs {
		exclusions = append(exclusions, fmt.Sprintf("--exclude=%s", excl))
	}
	args = append(args, exclusions...)
	return args
}

func (indexer *Indexer) GetTagFileNameForGemset(root string) string {
	return filepath.Join(root, fmt.Sprintf("%s.%s", indexer.TagFile, "gemset"))
}

func (indexer *Indexer) GemsetTagFileExists(root string) bool {
	return FileExists(indexer.GetTagFileNameForGemset(root))
}

func (indexer *Indexer) GetTagFileNameForProject(root string) string {
	return filepath.Join(root, fmt.Sprintf("%s.%s", indexer.TagFile, "project"))
}

func (indexer *Indexer) ProjectTagFileExists(root string) bool {
	return FileExists(indexer.GetTagFileNameForProject(root))
}
