package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func Test_Watcher_Start_Fires_OnFileCreation(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create and setup the Watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	fsWatcher.Add(path)
	watcher := &Watcher{
		Root: path,
		//MaxFrequency: 100 * time.Millisecond,
		fsWatcher: fsWatcher,
	}
	indexEvents := make(chan struct{})
	go watcher.Watch(indexEvents)

	//time.Sleep(2 * time.Second)

	// fire!
	TouchFile(t, filepath.Join(path, "test_file")).Close()

	// grab the event
	// !! if we read from a closed channel, this test is gonna pass no matter what
	<-indexEvents
}

func Test_Watcher_Fires_OnFileChange(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// create a file
	file := TouchFile(t, filepath.Join(path, "test_file"))

	// create and setup the Watcher
	fsWatcher, err := fsnotify.NewWatcher()
	assert.Nil(t, err)
	fsWatcher.Add(path)
	watcher := &Watcher{
		Root:         path,
		MaxFrequency: 100 * time.Millisecond,
		fsWatcher:    fsWatcher,
	}
	indexEvents := make(chan struct{})
	go watcher.Watch(indexEvents)

	// fire!
	file.WriteString("hello")
	file.Close()

	// grab the event
	<-indexEvents
}
