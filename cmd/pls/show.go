package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "show values for various resources",
	Aliases: []string{"s", "print", "list", "l"},
}

var showConfigsCmd = &cobra.Command{
	Use:     "configs",
	Aliases: []string{"c"},
	Short:   "show config values",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := utils.BashExec(fmt.Sprintf("cat %s", viper.ConfigFileUsed()))
		color.HiBlue(out)
	},
}

func init() {
	showCmd.AddCommand(showConfigsCmd)
}
