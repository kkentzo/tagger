package indexers

import (
	"fmt"
	"path/filepath"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	log "github.com/sirupsen/logrus"
)

type RvmIndexer struct {
	*Indexer
	RvmHandler RvmHandleable
}

func (indexer *RvmIndexer) Create(root string) Indexable {
	return indexer
}

func (indexer *RvmIndexer) Index(root string, event watchers.Event) {
	// Index the gemset (if necessary)
	if event.Names.Has("Gemfile.lock") || !indexer.GemsetTagFileExists(root) {
		indexer.indexGemset(root)
		event.Names.Remove("Gemfile.lock")
	}
	// Index the project
	indexer.Indexer.Index(root, event)
	// Join the tag files
	tagFiles := []string{
		indexer.TagFileName,
		indexer.GetTagFileNameForGemset(root),
	}
	// TODO: should this be more aggressive if one of the files does not exist?
	err := utils.ConcatFiles(filepath.Join(root, indexer.TagFileName), tagFiles, root)
	if err != nil {
		log.Error("concat:", tagFiles, err.Error())
	}
}

func (indexer *RvmIndexer) indexGemset(root string) {
	if indexer.RvmHandler.IsRuby(root) {
		args := indexer.GetGemsetArguments(root)
		if len(args) == 0 {
			return
		}
		out, err := utils.ExecInPath(indexer.Program, args, root)
		if err != nil {
			log.Error(string(out), err.Error())
		}
	}
}

func (indexer *RvmIndexer) GetGemsetArguments(root string) []string {
	args := indexer.GetGenericArguments(root)
	args = append(args, fmt.Sprintf("-f %s.gemset", indexer.TagFileName))
	if gemsetPath, err := indexer.RvmHandler.GemsetPath(root); err != nil {
		log.Error("Can not determine gemset path for rvm project at ", root)
		return []string{}
	} else {
		args = append(args, gemsetPath)
	}
	return args
}

func (indexer *RvmIndexer) GetTagFileNameForGemset(root string) string {
	return filepath.Join(root, fmt.Sprintf("%s.%s", indexer.TagFileName, "gemset"))
}

func (indexer *RvmIndexer) GemsetTagFileExists(root string) bool {
	return utils.FileExists(indexer.GetTagFileNameForGemset(root))
}
