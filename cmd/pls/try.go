package pls

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

// CURRENTLY UNUSED - USED FOR DEBUGGING TINGZ RN

var tryCmd = &cobra.Command{
	Use:     "try",
	Short:   "try to do something",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("TEST: %s", emoji.Sprint(":police_car_light:")))
	},
	Hidden: true,
}
