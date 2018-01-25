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

func Test_Watcher_handle_ReturnsTrue_OnFileRename(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	fname := filepath.Join(path, "test_file")
	TouchFile(t, fname).Close()

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	// TODO: the following cleanup statement hangs on macOS (test linux)
	// see fsnotify: kqueue.go#Close()
	//defer fsWatcher.Close()
	err = fsWatcher.Add(path)
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	new_fname := filepath.Join(path, "test_file_new")
	err = os.Rename(fname, new_fname)
	assert.Nil(t, err)

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_handle_ReturnTrue_OnDirectoryCreation(t *testing.T) {
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
	err = os.Mkdir(filepath.Join(path, "test_dir"), os.ModePerm)
	assert.Nil(t, err)

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_handle_ReturnTrue_OnDirectoryDeletion(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	dirName := filepath.Join(path, "test_dir")
	err = os.Mkdir(dirName, os.ModePerm)
	assert.Nil(t, err)

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	defer fsWatcher.Close()
	err = fsWatcher.Add(path)
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	err = os.RemoveAll(dirName)
	assert.Nil(t, err)

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_handle_ReturnTrue_OnDirectoryRename(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	dirName := filepath.Join(path, "test_dir")
	err = os.Mkdir(dirName, os.ModePerm)
	assert.Nil(t, err)

	// create and setup the filesystem watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	err = fsWatcher.Add(path)
	// TODO: the following cleanup statement hangs on macOS (test linux)
	// see fsnotify: kqueue.go#Close()
	//defer fsWatcher.Close()
	assert.Nil(t, err)
	// create the PathSet
	pathSet := NewPathSet([]string{})

	// fire!
	newDirName := dirName + "_new"
	err = os.Rename(dirName, newDirName)
	assert.Nil(t, err)
	os.RemoveAll(newDirName)

	event := <-fsWatcher.Events
	assert.True(t, handle(event, fsWatcher, pathSet))
}

func Test_Watcher_discover_IncludesAllDirectoriesUnderRoot(t *testing.T) {
}

func Test_Watcher_discover_DoesNotIncludeExcludedDirectories(t *testing.T) {
}
