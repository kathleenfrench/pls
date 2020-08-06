package gui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

// Spin is a spinner used to indicate a pending process
var Spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

// Exit says by and exits the program
func Exit() {
	color.HiCyan("\n%sbye", emoji.Sprint(":wave:"))
	os.Exit(0)
}
