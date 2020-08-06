package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

// PleaseHold is a logging helper for indicating active processes
func PleaseHold(msg string, extra interface{}) {
	if extra != nil {
		fmt.Println(fmt.Sprintf("%s%s [%v]", emoji.Sprint(":popcorn:"), fmt.Sprintf("%s...", msg), color.HiBlueString("%v", extra)))
	} else {
		fmt.Println(fmt.Sprintf("%s%s", emoji.Sprint(":popcorn"), msg))
	}
}

// OhNo is a logging helper for an uh-oh-esque message
func OhNo(msg string) {
	fmt.Println(fmt.Sprintf("%s%s", emoji.Sprint(":disappointed:"), color.HiRedString(msg)))
}
