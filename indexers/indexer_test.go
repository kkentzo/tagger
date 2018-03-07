package indexers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func CheckGenericArguments(t *testing.T, args []string) {
	assert.Contains(t, args, "-R")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "--exclude=.git")
}

func Test_Indexer_Deserialization(t *testing.T) {
	t.Skip("TODO")
}

func Test_Indexer_DefaultIndexer(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "ctags", indexer.Program)
	assert.Contains(t, indexer.Args, "-R")
	assert.Contains(t, indexer.Args, "-e")
	assert.Equal(t, "TAGS", indexer.TagFileName)
	assert.Equal(t, Generic, indexer.Type)
	assert.Contains(t, indexer.ExcludeDirs, ".git")
}

// TODO: This breaks in Travis CI
func Test_Indexer_Index_ShouldTriggerCommand(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := DefaultIndexer()

	indexer.Index(path, watchers.Event{})
	assert.True(t, utils.FileExists(filepath.Join(path, indexer.TagFileName)))
}

func Test_Indexer_CreateWatcher_ShouldReturnAWatcher(t *testing.T) {
	indexer := &Indexer{
		MaxPeriod: 2 * time.Second,
		Type:      Rvm,
	}
	watcher := indexer.CreateWatcher("foo").(*watchers.Watcher)
	defer watcher.Close()

	assert.Equal(t, "foo", watcher.Root)
	assert.Equal(t, 2*time.Second, watcher.MaxPeriod)
}

func Test_Indexer_GetGenericArguments(t *testing.T) {
	indexer := DefaultIndexer()
	args := indexer.GetGenericArguments("foo")
	CheckGenericArguments(t, args)
	assert.Equal(t, 3, len(args))
}

func Test_Indexer_GetProjectArguments(t *testing.T) {
	indexer := DefaultIndexer()
	args := indexer.GetProjectArguments("foo")
	CheckGenericArguments(t, args)
	assert.Contains(t, args, "-f TAGS")
	assert.Equal(t, ".", args[len(args)-1])
}
