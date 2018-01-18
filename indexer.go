package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TODO: Make Indexer an interface; create RubyIndexer and GenericIndexer
type Indexer struct {
	Program string
	Args    []string
	TagFile string `yaml:"tag_file"`
	Exclude []string
}

func (indexer *Indexer) Index(root string) {
	// TODO: implement exclusions and out file (-f)
	// TODO: Does ctags binary exist?
	cmd := exec.Command("ctags", "-R", "--languages=ruby", "-f TAGS", "--exclude=.git", ".", rvmGemsetPath(root))
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(out, err.Error())
	}
}

func (indexer *Indexer) args() string {
	// Add exclusions string, language, File
	return strings.Join(indexer.Args, " ")
}

func isRuby(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "Gemfile")); os.IsNotExist(err) {
		return false
	}
	return true
}

func rubyVersion(root string) string {
	rv, _ := ioutil.ReadFile(filepath.Join(root, ".ruby-version"))
	// TODO: Deal with error
	return strings.TrimSpace(string(rv))
}

func rubyGemset(root string) string {
	rg, _ := ioutil.ReadFile(filepath.Join(root, ".ruby-gemset"))
	// TODO: Catch errors!
	return strings.TrimSpace(string(rg))
}

func rvmGemsetPath(root string) string {
	// TODO: support other stuff besides rvm
	// cmd := exec.Command("bash", "-l", "-c", "rvm gemset use GEMSET_HERE; bundle list --paths")
	return filepath.Join(os.Getenv("HOME"),
		fmt.Sprintf(".rvm/gems/%s@%s/gems", rubyVersion(root), rubyGemset(root)))
}
