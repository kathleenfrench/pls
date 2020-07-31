package utils

import "runtime"

// GetPlatform returns the OS platform of the current machine
func GetPlatform() string {
	return runtime.GOOS
}
