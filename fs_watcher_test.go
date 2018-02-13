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

func (w *MockFsWatcher) Handle(e fsnotify.Event) bool {
	args := w.Called(e)
	return args.Get(0).(bool)
}

func (w *MockFsWatcher) Add(path string) error {
	args := w.Called(path)
	obj := args.Get(0)
	if obj != nil {
		return obj.(error)
	} else {
		return nil
	}
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

func Test_discover_IncludesAllDirectoriesUnderRoot(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create two top-level directories and one underneath the first
	dirNameA := filepath.Join(path, "dirA")
	err = os.Mkdir(dirNameA, os.ModePerm)
	assert.Nil(t, err)
	dirNameB := filepath.Join(path, "dirB")
	err = os.Mkdir(dirNameB, os.ModePerm)
	assert.Nil(t, err)
	dirNameAA := filepath.Join(dirNameA, "dirAA")
	err = os.Mkdir(dirNameAA, os.ModePerm)
	assert.Nil(t, err)

	dirs, err := discover(path, NewPathSet([]string{}))
	assert.Equal(t, 4, len(dirs))
	assert.Contains(t, dirs, path)
	assert.Contains(t, dirs, dirNameA)
	assert.Contains(t, dirs, dirNameB)
	assert.Contains(t, dirs, dirNameAA)
}

func Test_discover_DoesNotIncludeFiles(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create a top-level directory with a file inside
	dirName := filepath.Join(path, "dirA")
	err = os.Mkdir(dirName, os.ModePerm)
	assert.Nil(t, err)
	TouchFile(t, filepath.Join(dirName, "test_file"))

	dirs, err := discover(path, NewPathSet([]string{}))
	assert.Equal(t, 2, len(dirs))
	assert.Contains(t, dirs, path)
	assert.Contains(t, dirs, dirName)
}

func Test_discover_DoesNotIncludeExcludedDirectories(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create a top-level directory and one inside to be ignored
	dirName := filepath.Join(path, "dirA")
	err = os.Mkdir(dirName, os.ModePerm)
	assert.Nil(t, err)
	ignoredDir := filepath.Join(dirName, "log")
	err = os.Mkdir(ignoredDir, os.ModePerm)
	assert.Nil(t, err)

	dirs, err := discover(path, NewPathSet([]string{"log"}))
	assert.Equal(t, 2, len(dirs))
	assert.Contains(t, dirs, path)
	assert.Contains(t, dirs, dirName)
}
