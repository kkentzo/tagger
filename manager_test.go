package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewManager(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	config := &Config{
		Projects: []struct{ Path string }{
			{Path: path},
		},
		Indexer: &GenericIndexer{MaxFrequency: 2 * time.Second},
	}
	manager := NewManager(config)
	assert.Equal(t, config.Indexer, manager.indexer)
	assert.Contains(t, manager.projects, path)
}

func Test_Manager_Add_WillNotAddProject_WhenPathDoesNotExist(t *testing.T) {
	config := &Config{
		Projects: []struct{ Path string }{},
		Indexer:  &GenericIndexer{MaxFrequency: 2 * time.Second},
	}
	manager := NewManager(config)

	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	os.RemoveAll(path)

	manager.Add(path)
	assert.NotContains(t, manager.projects, path)
}

func Test_Manager_Add_WillAddProject_WhenPathExists(t *testing.T) {
	config := &Config{
		Projects: []struct{ Path string }{},
		Indexer:  &GenericIndexer{MaxFrequency: 2 * time.Second},
	}
	manager := NewManager(config)

	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	manager.Add(path)
	assert.Contains(t, manager.projects, path)
}

func Test_Manager_Add_WillExpandTildeToHomeDir(t *testing.T) {

}

func Test_Manager_Remove_WillRemoveProjectFromManager(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	config := &Config{
		Projects: []struct{ Path string }{
			{Path: path},
		},
		Indexer: &GenericIndexer{MaxFrequency: 2 * time.Second},
	}
	manager := NewManager(config)
	assert.Contains(t, manager.projects, path)
	manager.Remove(path)
	assert.NotContains(t, manager.projects, path)
}

func Test_fileExists_ReturnsTrue_IfFileExists(t *testing.T) {
	assert.True(t, fileExists("/tmp"))
}

func Test_fileExists_ReturnsFalse_IfFileDoesNotExist(t *testing.T) {
	assert.False(t, fileExists("/foo"))
}

func Test_Canonicalize(t *testing.T) {
	home := os.Getenv("HOME")
	assert.NotEmpty(t, home)
	var testCases = []struct {
		path         string
		expandedPath string
	}{
		{"~", home},
		{"~/foo/bar", fmt.Sprintf("%s/foo/bar", home)},
		{"/foo/bar", "/foo/bar"},
		{"", ""},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.expandedPath, Canonicalize(testCase.path))
	}
}
