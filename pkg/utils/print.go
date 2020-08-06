package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

// PrintError outputs a standard error message to the user's console
func PrintError(msg interface{}) {
	color.HiRed(fmt.Sprintf("[ERROR %s]: %v", emoji.Sprint(":skull:"), msg))
}
