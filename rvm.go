package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func isRuby(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "Gemfile")); os.IsNotExist(err) {
		return false
	}
	return true
}

func isRvm(root string) bool {
	return isRuby(root) &&
		fileExists(filepath.Join(root, ".ruby-version")) &&
		fileExists(filepath.Join(root, ".ruby-gemset"))
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
