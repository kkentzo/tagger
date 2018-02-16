package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIndexer struct {
	mock.Mock
}

func (indexer *MockIndexer) Index(root string) {
	indexer.Called(root)
}

func (indexer *MockIndexer) CreateWatcher(root string) Watchable {
	args := indexer.Called(root)
	return args.Get(0).(Watchable)
}

func Test_Indexer_DefaultIndexer(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "ctags", indexer.Program)
	assert.Contains(t, indexer.Args, "-R")
	assert.Contains(t, indexer.Args, "-e")
	assert.Equal(t, "TAGS", indexer.TagFile)
	assert.Equal(t, Generic, indexer.Type)
	assert.Contains(t, indexer.ExcludeDirs, ".git")
}

func Test_Indexer_Index_ShouldTriggerCommand(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := &Indexer{
		Program: "touch",
		Args:    []string{"aaa"},
		Type:    Generic,
	}

	indexer.Index(path)
	assert.True(t, FileExists(filepath.Join(path, "aaa")))
}

func Test_Indexer_CreateWatcher_ShouldReturnAWatcher(t *testing.T) {
	indexer := &Indexer{
		ExcludeDirs:  []string{".git"},
		MaxFrequency: 2 * time.Second,
	}
	watcher := indexer.CreateWatcher("foo").(*Watcher)
	defer watcher.Close()

	assert.Equal(t, "foo", watcher.Root)
	assert.Equal(t, 2*time.Second, watcher.MaxFrequency)
}

func Test_Indexer_GetArguments_WhenIndexerIsGeneric(t *testing.T) {
	indexer := DefaultIndexer()
	args := indexer.GetArguments("foo")
	assert.Contains(t, args, "-f TAGS")
	assert.Contains(t, args, "-R")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "--exclude=.git")
	assert.Contains(t, args, ".")
}

func Test_Indexer_GetArguments_WhenIndexerIsRvm(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	f := TouchFile(t, filepath.Join(path, ".ruby-version"))
	f.Write([]byte("2.1.3"))
	f.Close()
	f = TouchFile(t, filepath.Join(path, ".ruby-gemset"))
	f.Write([]byte("foo"))
	f.Close()

	indexer := DefaultIndexer()
	indexer.Type = Rvm
	args := indexer.GetArguments(path)
	assert.Contains(t, args, "-f TAGS")
	assert.Contains(t, args, "-R")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "--exclude=.git")
	assert.Contains(t, args, ".")
	gp, err := rvmGemsetPath(path)
	assert.Nil(t, err)
	assert.Contains(t, args, gp)
}

func Test_Indexer_GetArguments_WhenIndexerIsRvm_ButProjectIsNot(t *testing.T) {
	indexer := DefaultIndexer()
	indexer.Type = Rvm
	args := indexer.GetArguments("foo")
	assert.Contains(t, args, "-f TAGS")
	assert.Contains(t, args, "-R")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "--exclude=.git")
	assert.Contains(t, args, ".")
	assert.Equal(t, 5, len(args))
}
