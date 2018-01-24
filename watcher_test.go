package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func Test_Watcher_handle_ReturnsTrue_OnFileCreation(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	defer fsWatcher.Close()
	err = fsWatcher.Add(path)
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	TouchFile(t, filepath.Join(path, "test_file")).Close()

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_handle_ReturnsTrue_OnFileChange(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	file := TouchFile(t, filepath.Join(path, "test_file"))
	defer file.Close()

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	defer fsWatcher.Close()
	err = fsWatcher.Add(path)
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	file.WriteString("hello")

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_handle_ReturnsTrue_OnFileDeletion(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	fname := filepath.Join(path, "test_file")
	TouchFile(t, fname).Close()

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	defer fsWatcher.Close()
	err = fsWatcher.Add(path)
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	os.Remove(fname)

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}
