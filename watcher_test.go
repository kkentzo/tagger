package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Watcher_Fires_OnFileCreation(t *testing.T) {
	// create the project directory
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)
	// create the Watcher
	watcher := &Watcher{
		Root: path,
	}
	indexEvents := make(chan struct{})
	go watcher.Watch(indexEvents)
	// fire
	os.OpenFile(filepath.Join(path, "test_file"), os.O_RDONLY|os.O_CREATE, 0666)
	// grab the event
	<-indexEvents
}
