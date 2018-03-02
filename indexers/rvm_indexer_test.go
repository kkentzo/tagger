package indexers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Indexer_GetGemsetArguments_WhenIndexerIsRvm(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	// Prepare rvm-specific files
	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	f := TouchFile(t, filepath.Join(path, ".ruby-version"))
	f.Write([]byte("2.1.3"))
	f.Close()
	f = TouchFile(t, filepath.Join(path, ".ruby-gemset"))
	f.Write([]byte("foo"))
	f.Close()

	indexer := DefaultIndexer()
	indexer.Type = Rvm
	args := indexer.GetGemsetArguments(path)
	CheckGenericArguments(t, args)
	gp, err := rvmGemsetPath(path)
	assert.Nil(t, err)
	assert.Contains(t, args, "-f TAGS.gemset")
	assert.Equal(t, gp, args[len(args)-1])
}

func Test_Indexer_GetTagFileNameForGemset(t *testing.T) {
	indexer := DefaultIndexer()
	assert.Equal(t, "aaa/TAGS.gemset", indexer.GetTagFileNameForGemset("aaa"))

}

func Test_Indexer_GemsetTagFileExists_ReturnsTrue_WhenTagFileExists(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "TAGS.gemset")).Close()

	indexer := DefaultIndexer()
	assert.True(t, indexer.GemsetTagFileExists(path))
}

func Test_Indexer_GemsetTagFileExists_ReturnsFalse_WhenTagFileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := DefaultIndexer()
	assert.False(t, indexer.GemsetTagFileExists(path))
}
