package indexers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	log "github.com/sirupsen/logrus"
)

const TagFilePrefix string = "TAGS"

type Indexable interface {
	Create() Indexable
	Index(string, watchers.Event)
	CreateWatcher(string) watchers.Watchable
}

type IndexerType string

const (
	Generic IndexerType = "generic"
	Rvm                 = "rvm"
)

type Indexer struct {
	Program       string
	Args          []string
	TagFilePrefix string `yaml:"tag_file"`
	Type          IndexerType
	ExcludeDirs   []string      `yaml:"exclude"`
	MaxPeriod     time.Duration `yaml:"max_period"`
}

func DefaultIndexer() *Indexer {
	return &Indexer{
		Program:       "ctags",
		Args:          []string{"-R", "-e"},
		TagFilePrefix: TagFilePrefix,
		Type:          Generic,
		ExcludeDirs:   []string{".git"},
	}
}

func (indexer *Indexer) Index(root string, event watchers.Event) {
	indexer.indexProject(root)
	if event.IsSpecial || !indexer.GemsetTagFileExists(root) {
		indexer.indexGemset(root)
	}
	tagFiles := []string{
		indexer.GetTagFileNameForProject(root),
		indexer.GetTagFileNameForGemset(root),
	}
	// TODO: This should be robust if one of the files does not exist
	err := utils.ConcatFiles(filepath.Join(root, indexer.TagFilePrefix), tagFiles, root)
	if err != nil {
		log.Error("concat", tagFiles, err.Error())
	}
}

func (indexer *Indexer) Create() Indexable {
	return indexer
}

func (indexer *Indexer) CreateWatcher(root string) watchers.Watchable {
	w := watchers.NewWatcher(root, indexer.ExcludeDirs,
		indexer.TagFilePrefix, indexer.MaxPeriod)
	if indexer.Type == Rvm {
		w.SpecialFile = "Gemfile.lock"
	}
	return w
}

func (indexer *Indexer) indexProject(root string) {
	args := indexer.GetProjectArguments(root)
	out, err := utils.ExecInPath(indexer.Program, args, root)
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
		out, err := utils.ExecInPath(indexer.Program, args, root)
		if err != nil {
			log.Error(out, err.Error())
		}
	}
}

func (indexer *Indexer) GetProjectArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s.project", indexer.TagFilePrefix))
	args = append(args, ".")
	return args
}

func (indexer *Indexer) GetGemsetArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s.gemset", indexer.TagFilePrefix))
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
	return filepath.Join(root, fmt.Sprintf("%s.%s", indexer.TagFilePrefix, "gemset"))
}

func (indexer *Indexer) GemsetTagFileExists(root string) bool {
	return utils.FileExists(indexer.GetTagFileNameForGemset(root))
}

func (indexer *Indexer) GetTagFileNameForProject(root string) string {
	return filepath.Join(root, fmt.Sprintf("%s.%s", indexer.TagFilePrefix, "project"))
}

func (indexer *Indexer) ProjectTagFileExists(root string) bool {
	return utils.FileExists(indexer.GetTagFileNameForProject(root))
}
