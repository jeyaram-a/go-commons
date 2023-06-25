package gocommons

import (
	"errors"
	"fmt"
	"os"
)

func FolderExists(path string) bool {
	fs, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fs.IsDir()
}

func IsFolderEmpty(path string) (bool, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !fs.IsDir() {
		return false, errors.New(fmt.Sprintf("%s not a folder", path))
	}
	contents, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(contents) == 0, nil
}
