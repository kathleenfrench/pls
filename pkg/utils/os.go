package utils

import (
	"os"
	"runtime"
)

// GetPlatform returns the OS platform of the current machine
func GetPlatform() string {
	return runtime.GOOS
}

// ExitWithError prints a clear error and exits the program
func ExitWithError(msg interface{}) {
	PrintError(msg)
	os.Exit(1)
}
