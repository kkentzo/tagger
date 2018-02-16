package main

import (
	"os"
	"os/exec"
	"strings"
)

func Canonicalize(path string) string {
	if strings.Contains(path, "~") {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ExecInPath(cmd string, args []string, path string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Dir = path
	out, err := command.CombinedOutput()
	return string(out), err
}
