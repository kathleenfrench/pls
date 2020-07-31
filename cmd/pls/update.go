package pls

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "update various resources like configs and pls itself",
}

var updateCfgSubCmd = &cobra.Command{
	Use:   "config",
	Short: "update your pls configs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updating pls configs...")
	},
}

func init() {
	updateCmd.AddCommand(updateCfgSubCmd)
}
