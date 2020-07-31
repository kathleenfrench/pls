package utils

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintError outputs a standard error message to the user's console
func PrintError(msg interface{}) {
	color.HiRed(fmt.Sprintf("[ERROR]: %v", msg))
}
