package main

import (
	"os"
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
	// TODO: check for err value notexists
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
