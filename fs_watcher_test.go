package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFsWatcher struct {
	mock.Mock
}

func (w *MockFsWatcher) Add(path string) error {
	args := w.Called(path)
	return args.Get(0).(error)
}

func (w *MockFsWatcher) Remove(path string) error {
	args := w.Called(path)
	return args.Get(0).(error)
}

func (w *MockFsWatcher) Close() error {
	args := w.Called()
	return args.Get(0).(error)
}

func (w *MockFsWatcher) Events() chan fsnotify.Event {
	args := w.Called()
	return args.Get(0).(chan fsnotify.Event)
}

func (w *MockFsWatcher) Errors() chan error {
	args := w.Called()
	return args.Get(0).(chan error)
}

// These functions test the functionality of fsnotify.Watcher

func Test_FsWatcher_Fires_OnFileCreation(t *testing.T) {
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

	// fire!
	TouchFile(t, filepath.Join(path, "test_file")).Close()

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnFileChange(t *testing.T) {
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

	// fire!
	file.WriteString("hello")

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnFileDeletion(t *testing.T) {
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

	// fire!
	os.Remove(fname)

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnFileRename(t *testing.T) {
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

	// fire!
	new_fname := filepath.Join(path, "test_file_new")
	err = os.Rename(fname, new_fname)
	assert.Nil(t, err)

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnDirectoryCreation(t *testing.T) {
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

	// fire!
	err = os.Mkdir(filepath.Join(path, "test_dir"), os.ModePerm)
	assert.Nil(t, err)

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnDirectoryDeletion(t *testing.T) {
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

	// fire!
	err = os.RemoveAll(dirName)
	assert.Nil(t, err)

	assert.NotNil(t, <-fsWatcher.Events)
}

func Test_FsWatcher_Fires_OnDirectoryRename(t *testing.T) {
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

	// fire!
	newDirName := dirName + "_new"
	err = os.Rename(dirName, newDirName)
	assert.Nil(t, err)
	os.RemoveAll(newDirName)

	assert.NotNil(t, <-fsWatcher.Events)
}
