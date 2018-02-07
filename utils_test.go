package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TouchFile(t *testing.T, fname string) *os.File {
	f, err := os.Create(fname)
	assert.Nil(t, err)
	return f
}

func Test_FileExists_ReturnsTrue_IfFileExists(t *testing.T) {
	assert.True(t, FileExists(os.Getenv("HOME")))
}

func Test_FileExists_ReturnsFalse_IfFileDoesNotExist(t *testing.T) {
	assert.False(t, FileExists("/foo"))
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
