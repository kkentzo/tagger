package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWatcher struct {
	mock.Mock
}

func CreateMockWatcher() *MockWatcher {
	watcher := &MockWatcher{}
	watcher.On("Watch", mock.AnythingOfType("*context.cancelCtx"))
	watcher.On("Close")
	watcher.On("Events")
	return watcher
}

func (watcher *MockWatcher) Watch(ctx context.Context) {
	watcher.Called(ctx)
}

func (watcher *MockWatcher) Events() chan struct{} {
	args := watcher.Called()
	if len(args) > 0 {
		return args.Get(0).(chan struct{})
	} else {
		return make(chan struct{})
	}
}

func (watcher *MockWatcher) Close() {
	watcher.Called()
}

func Test_NewWatcher(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	defer watcher.Close()
	assert.Equal(t, "foo", watcher.Root)
	assert.Contains(t, watcher.Exclusions, "excl")
	assert.Equal(t, 2*time.Second, watcher.MaxFrequency)
	assert.IsType(t, &FsWatcher{}, watcher.fsWatcher)
	assert.IsType(t, make(chan struct{}), watcher.events)
}

func Test_Watcher_Events_ReturnsTheChannel(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	defer watcher.Close()
	var s struct{}
	go func(w *Watcher) { watcher.Events() <- s }(watcher)
	assert.Equal(t, s, <-watcher.Events())
}

func Test_Watcher_Close_ClosesTheChannels(t *testing.T) {
	watcher := NewWatcher("foo", []string{"excl"}, 2*time.Second)
	watcher.Close()
	_, open := <-watcher.Events()
	assert.False(t, open)
	_, open = <-watcher.fsWatcher.Events()
	assert.False(t, open)
}

func Test_Watcher_Watch_ShouldCallHandlerFunc_OnFsNotify_Event(t *testing.T) {
	// create and setup the filesystem watcher
	events := make(chan fsnotify.Event)
	fsWatcher := &MockFsWatcher{}
	fsWatcher.On("Events").Return(events)
	// substitute the HandlerFunc of the MockWatcher

}

func Test_Watcher_Watch_ShouldReindex_WhenTickerTicks(t *testing.T) {
	t.Skip("Need to stub fswatcher")
}

func Test_Watcher_Watch_ShouldNotReindex_WhenTagFileChanges(t *testing.T) {
	t.Skip("Need to stub fswatcher")
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

func Test_isDirectory_ReturnsTrue_WhenPathIsDirectory(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	result, err := isDirectory(path)
	assert.True(t, result)
	assert.Nil(t, err)
}

func Test_isDirectory_ReturnsFalse_WhenPathIsFile(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	fname := filepath.Join(path, "test_file")
	TouchFile(t, fname).Close()

	result, err := isDirectory(fname)
	assert.False(t, result)
	assert.Nil(t, err)
}
