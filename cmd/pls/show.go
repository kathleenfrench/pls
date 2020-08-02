package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show values for various resources",
}

var showConfigsCmd = &cobra.Command{
	Use:   "configs",
	Short: "show config values",
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			color.HiBlue(fmt.Sprintf("%s: %v", k, v))
		}
	},
}

func init() {
	showCmd.AddCommand(showConfigsCmd)
}
