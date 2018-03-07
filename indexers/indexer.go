package indexers

import (
	"fmt"
	"time"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	log "github.com/sirupsen/logrus"
)

const TagFileName string = "TAGS"

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
	Program     string
	Args        []string
	TagFileName string `yaml:"tag_file"`
	Type        IndexerType
	ExcludeDirs []string      `yaml:"exclude"`
	MaxPeriod   time.Duration `yaml:"max_period"`
}

func DefaultIndexer() *Indexer {
	return &Indexer{
		Program:     "ctags",
		Args:        []string{"-R", "-e"},
		TagFileName: TagFileName,
		Type:        Generic,
		ExcludeDirs: []string{".git"},
	}
}

func (indexer *Indexer) Index(root string, event watchers.Event) {
	indexer.indexProject(root)
}

func (indexer *Indexer) Create() Indexable {
	if indexer.Type == Rvm {
		return &RvmIndexer{
			Indexer:    indexer,
			RvmHandler: DefaultRvmHandler(),
		}
	} else {
		return indexer
	}
}

func (indexer *Indexer) CreateWatcher(root string) watchers.Watchable {
	return watchers.NewWatcher(root, indexer.ExcludeDirs,
		indexer.TagFileName, indexer.MaxPeriod)
}

func (indexer *Indexer) indexProject(root string) {
	args := indexer.GetProjectArguments(root)
	out, err := utils.ExecInPath(indexer.Program, args, root)
	if err != nil {
		log.Error(string(out), err.Error())
	}
}

func (indexer *Indexer) GetProjectArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s", indexer.TagFileName))
	args = append(args, ".")
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
