package indexers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Rvm_isRuby_ReturnsTrue_WhenGemfileExists(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rvm := &RvmHandler{}

	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	assert.True(t, rvm.IsRuby(path))
}

func Test_Rvm_isRuby_ReturnsFalse_WhenGemfileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rvm := &RvmHandler{}

	assert.False(t, rvm.IsRuby(path))
}
