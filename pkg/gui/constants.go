package gui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Spin is a spinner used to indicate a pending process
var Spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

// Exit says by and exits the program
func Exit() {
	color.HiGreen("bye!")
	os.Exit(0)
}
