package gui

import (
	"os"

	"github.com/fatih/color"
)

// Exit says by and exits the program
func Exit() {
	color.HiGreen("bye!")
	os.Exit(0)
}
