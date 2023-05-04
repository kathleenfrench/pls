package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// FileExists checks for whether a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks for whether a directory exists
func DirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDir creates a directory if it does not exist
func CreateDir(path string) error {
	if DirExists(path) {
		return nil
	}

	color.HiBlue(fmt.Sprintf("creating directory at %s...", path))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// CreateFile creates a file if it does not exist
func CreateFile(path string) error {
	if !FileExists(path) {
		color.HiBlue(fmt.Sprintf("creating file at %s...", path))
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	return nil
}
