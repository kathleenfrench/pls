package utils

import (
	"os"
)

// FileExists checks for whether a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return false
}

// DirExists checks for whether a directory exists
func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

// CreateDir creates a directory if it does not exist
func CreateDir(path string) error {
	if DirExists(path) {
		return nil
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// CreateFile creates a file if it does not exist
func CreateFile(path string) error {
	if !FileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	return nil
}
