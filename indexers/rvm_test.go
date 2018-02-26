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

	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	assert.True(t, isRuby(path))
}

func Test_Rvm_isRuby_ReturnsFalse_WhenGemfileDoesNotExist(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	assert.False(t, isRuby(path))
}

func Test_Rvm_isRvm_ReturnsTrue_WhenAllConditionsAreMet(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	TouchFile(t, filepath.Join(path, "Gemfile")).Close()
	TouchFile(t, filepath.Join(path, ".ruby-version")).Close()
	TouchFile(t, filepath.Join(path, ".ruby-gemset")).Close()

	assert.True(t, isRvm(path))
}

func Test_Rvm_rubyVersion_ReturnsTheVersion(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	f := TouchFile(t, filepath.Join(path, ".ruby-version"))
	_, err = f.Write([]byte("2.1.3"))
	f.Close()
	assert.Nil(t, err)

	rv, err := rubyVersion(path)
	assert.Equal(t, "2.1.3", rv)
	assert.Nil(t, err)
}

func Test_Rvm_rubyVersion_ReturnsError_WhenNotFound(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rv, err := rubyVersion(path)
	assert.Empty(t, rv)
	assert.NotNil(t, err)
}

func Test_Rvm_rubyGemset_ReturnsGemset(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	f := TouchFile(t, filepath.Join(path, ".ruby-gemset"))
	_, err = f.Write([]byte("gem"))
	f.Close()
	assert.Nil(t, err)

	rg, err := rubyGemset(path)
	assert.Equal(t, "gem", rg)
	assert.Nil(t, err)
}

func Test_Rvm_rubyGemsetPath_ReturnsGemset(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	f := TouchFile(t, filepath.Join(path, ".ruby-gemset"))
	_, err = f.Write([]byte("gem"))
	f.Close()
	assert.Nil(t, err)

	f = TouchFile(t, filepath.Join(path, ".ruby-version"))
	_, err = f.Write([]byte("ruby-2.1.3"))
	f.Close()
	assert.Nil(t, err)

	gsPath, err := rvmGemsetPath(path)
	expPath := filepath.Join(os.Getenv("HOME"), ".rvm/gems/ruby-2.1.3@gem/gems")
	assert.Equal(t, expPath, gsPath)
	assert.Nil(t, err)
}

func Test_Rvm_rvmGemsetPath_ReturnsError_WhenNotFound(t *testing.T) {
	path, err := ioutil.TempDir("", "tagger-tests")
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	rg, err := rubyGemset(path)
	assert.Empty(t, rg)
	assert.NotNil(t, err)
}
