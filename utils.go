package main

import (
	"errors"
	"fmt"
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

func ExecInPath(cmd string, args []string, path string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	command.Dir = path
	out, err := command.CombinedOutput()
	return out, err
}

func ConcatFiles(to string, files []string, path string) error {
	out, err := ExecInPath("cat", files, path)
	if err != nil {
		return errors.New(fmt.Sprint(string(out), err.Error()))
	}
	f, err := os.Create(to)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(out)
	return nil
}
