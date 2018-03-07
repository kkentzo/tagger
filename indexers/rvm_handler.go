package indexers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kkentzo/tagger/utils"
)

type RvmHandleable interface {
	IsRuby(path string) bool
	GemsetPath(path string) (string, error)
}

type RvmHandler struct {
	Command string
	Args    []string
}

func DefaultRvmHandler() *RvmHandler {
	return &RvmHandler{
		Command: "/bin/bash",
		Args: []string{
			"-c",
			"source \"$HOME/.rvm/scripts/rvm\"; cd .; rvm gemset gemdir"},
	}
}

func (rvm *RvmHandler) IsRuby(path string) bool {
	return utils.FileExists(filepath.Join(path, "Gemfile"))
}

func (rvm *RvmHandler) GemsetPath(path string) (string, error) {
	cmd := "/bin/bash"
	args := []string{
		"-c",
		"source \"$HOME/.rvm/scripts/rvm\"; cd .; rvm gemset gemdir"}
	out, err := utils.ExecInPath(cmd, args, path)
	if err != nil {
		return "", errors.New(fmt.Sprint(string(out), err.Error()))
	} else {
		gemset := strings.TrimSpace(string(out))
		return filepath.Join(gemset, "gems"), nil
	}
}

// OLD FUNCTIONS
// func rubyVersion(root string) (string, error) {
// 	rv, err := ioutil.ReadFile(filepath.Join(root, ".ruby-version"))
// 	if err != nil {
// 		return "", err
// 	} else {
// 		return strings.TrimSpace(string(rv)), nil
// 	}
// }

// func rubyGemset(root string) (string, error) {
// 	rg, err := ioutil.ReadFile(filepath.Join(root, ".ruby-gemset"))
// 	if err != nil {
// 		return "", err
// 	} else {
// 		return strings.TrimSpace(string(rg)), nil
// 	}
// }

// func rvmGemsetPathFromFiles(root string) (string, error) {
// 	rv, err := rubyVersion(root)
// 	if err != nil {
// 		return "", err
// 	}
// 	rg, err := rubyGemset(root)
// 	if err != nil {
// 		return "", err
// 	}
// 	path := filepath.Join(os.Getenv("HOME"),
// 		fmt.Sprintf(".rvm/gems/%s@%s/gems", rv, rg))
// 	return path, nil

// }
