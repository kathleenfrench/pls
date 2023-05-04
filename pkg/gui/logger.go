package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

// PleaseHold is a logging helper for indicating active processes
func PleaseHold(msg string, extra interface{}) {
	if extra != nil {
		fmt.Printf("%s %s [%v]\n", emoji.Sprint(":popcorn:"), msg, color.HiBlueString("%v", extra))
	} else {
		fmt.Printf("%s %s\n", emoji.Sprint(":popcorn:"), msg)
	}
}

// Log is a logging helper that allows custom input with the emoji
func Log(e interface{}, msg string, extra interface{}) {
	if extra != nil {
		fmt.Printf("%s %s [%v]\n", emoji.Sprint(e), msg, color.HiBlueString("%v", extra))
	} else {
		fmt.Printf("%s %s\n", emoji.Sprint(e), msg)
	}
}

// OhNo is a logging helper for an uh-oh-esque message
func OhNo(msg string) {
	fmt.Printf("%s%s\n", emoji.Sprint(":disappointed:"), color.HiRedString(msg))
}
