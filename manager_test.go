package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_NewManager(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)
	projects := []struct{ Path string }{
		{Path: path},
	}
	watcher := CreateMockWatcher()

	indexer := &MockIndexer{}
	indexer.On("Create", path).Return(indexer)
	indexer.On("CreateWatcher", path).Return(watcher)
	indexer.On("Index", path, mock.AnythingOfType("watchers.Event"))

	manager := NewManager(indexer, projects)

	assert.Equal(t, indexer, manager.indexer)
	assert.Contains(t, manager.projects, path)
}

func Test_Manager_Add_WillNotAddProject_WhenPathDoesNotExist(t *testing.T) {
	projects := []struct{ Path string }{}
	indexer := &MockIndexer{}
	manager := NewManager(indexer, projects)

	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	os.RemoveAll(path)

	indexer.On("Index", path, mock.AnythingOfType("watchers.Event"))

	manager.Add(path)
	assert.NotContains(t, manager.projects, path)
}

func Test_Manager_Add_WillAddProject_WhenPathExists(t *testing.T) {
	projects := []struct{ Path string }{}
	indexer := &MockIndexer{}
	manager := NewManager(indexer, projects)

	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	watcher := CreateMockWatcher()
	indexer.On("Create", path).Return(indexer)
	indexer.On("CreateWatcher", path).Return(watcher)
	indexer.On("Index", path, mock.AnythingOfType("watchers.Event"))
	manager.Add(path)

	assert.Contains(t, manager.projects, path)
}

func Test_Manager_Remove_WillRemoveProjectFromManager(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	projects := []struct{ Path string }{
		{Path: path},
	}
	watcher := CreateMockWatcher()
	indexer := &MockIndexer{}
	indexer.On("Create", path).Return(indexer)
	indexer.On("CreateWatcher", path).Return(watcher)
	indexer.On("Index", path, mock.AnythingOfType("watchers.Event"))
	manager := NewManager(indexer, projects)
	assert.Contains(t, manager.projects, path)

	manager.Remove(path)
	assert.NotContains(t, manager.projects, path)
}
