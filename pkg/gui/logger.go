package gui

import (
	"fmt"

	"github.com/fatih/color"
)

// PleaseHold is a logging helper for indicating active processes
func PleaseHold(msg string, extra interface{}) {
	if extra != nil {
		fmt.Println(fmt.Sprintf("%s %s [%v]", color.HiYellowString("<pls hold>"), msg, color.HiBlueString("%v", extra)))
	} else {
		fmt.Println("%s %s", color.HiYellowString("<pls hold>"), msg)
	}
}

// OhNo is a logging helper for an uh-oh-esque message
func OhNo(msg string) {
	fmt.Println(fmt.Sprintf("%s %s %s", color.HiRedString("<plsdontcry>"), msg, color.HiYellowString(":'(")))
}
