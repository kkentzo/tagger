package indexers

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkentzo/tagger/utils"
	"github.com/kkentzo/tagger/watchers"
	"github.com/stretchr/testify/assert"
)

func Test_RvmIndexer_Index_ShouldIndexGemset_WhenGemsetTagFile_DoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rvm := &MockRvmHandler{}
	indexer := RvmIndexer{
		Indexer:    DefaultIndexer(),
		RvmHandler: rvm,
	}
	rvm.On("GemsetPath", path).Return(path, nil)
	rvm.On("IsRuby", path).Return(true)
	indexer.Index(path, watchers.NewEvent())
	assert.True(t, utils.FileExists(filepath.Join(path, "TAGS.gemset")))
	assert.True(t, utils.FileExists(filepath.Join(path, "TAGS")))
}

func Test_RvmIndexer_Index_ShouldIndexGemset_WhenEventNames_ContainGemfileLock(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rvm := &MockRvmHandler{}
	indexer := RvmIndexer{
		Indexer:    DefaultIndexer(),
		RvmHandler: rvm,
	}
	rvm.On("GemsetPath", path).Return(path, nil)
	rvm.On("IsRuby", path).Return(true)

	event := watchers.NewEvent()
	event.Names.Add("Gemfile.lock")
	indexer.Index(path, event)
	assert.True(t, utils.FileExists(filepath.Join(path, "TAGS.gemset")))
	assert.True(t, utils.FileExists(filepath.Join(path, "TAGS")))
	assert.False(t, event.Names.Has("Gemfile.lock"))
}

func Test_RvmIndexer_Index_ShouldConcatTagFiles(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	f := TouchFile(t, filepath.Join(path, "hello.rb"))
	f.Write([]byte("def hello; end"))
	f.Close()

	rvm := &MockRvmHandler{}
	indexer := RvmIndexer{
		Indexer:    DefaultIndexer(),
		RvmHandler: rvm,
	}
	rvm.On("GemsetPath", path).Return(path, nil)
	rvm.On("IsRuby", path).Return(true)

	event := watchers.NewEvent()
	event.Names.Add("Gemfile.lock")
	indexer.Index(path, event)
	contents, _ := ioutil.ReadFile(filepath.Join(path, "TAGS"))
	assert.Equal(t, 2, strings.Count(string(contents), "hello.rb,24"))
}

func Test_RvmIndexer_GetGemsetArguments_WhenGemsetPathCanBeDetermined(t *testing.T) {
	rvm := &MockRvmHandler{}
	indexer := RvmIndexer{
		Indexer:    DefaultIndexer(),
		RvmHandler: rvm,
	}

	rvm.On("GemsetPath", "project_path").Return("gemset_path", nil)
	args := indexer.GetGemsetArguments("project_path")
	CheckGenericArguments(t, args)

	assert.Contains(t, args, "-f TAGS.gemset")
	assert.Equal(t, "gemset_path", args[len(args)-1])
}

func Test_RvmIndexer_GetGemsetArguments_WhenGemsetPathCanNotBeDetermined(t *testing.T) {
	rvm := &MockRvmHandler{}
	indexer := RvmIndexer{
		Indexer:    DefaultIndexer(),
		RvmHandler: rvm,
	}

	rvm.On("GemsetPath", "project_path").Return("", errors.New("Something went wrong"))
	args := indexer.GetGemsetArguments("project_path")
	assert.Empty(t, args)
}

func Test_RvmIndexer_GetTagFileNameForGemset(t *testing.T) {
	indexer := &RvmIndexer{Indexer: DefaultIndexer()}
	assert.Equal(t, "foo/TAGS.gemset", indexer.GetTagFileNameForGemset("foo"))
}

func Test_RvmIndexer_GemsetTagFileExists_ReturnsTrue_WhenTagFileExists(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "TAGS.gemset")).Close()
	indexer := &RvmIndexer{Indexer: DefaultIndexer()}
	assert.True(t, indexer.GemsetTagFileExists(path))
}

func Test_RvmIndexer_GemsetTagFileExists_ReturnsFalse_WhenTagFileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	indexer := &RvmIndexer{Indexer: DefaultIndexer()}
	assert.False(t, indexer.GemsetTagFileExists(path))
}
