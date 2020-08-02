package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "update various resources like configs and pls itself",
}

var updateCfgSubCmd = &cobra.Command{
	Use:     "configs",
	Short:   "update your pls configs",
	Aliases: []string{"cnfs"},
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("TODO: INTERACTIVE DROPDOWN")
		config.UpdatePrompt(viper.AllSettings())
	},
}

var updateSingleCfgSubCmd = &cobra.Command{
	Use:     "config",
	Short:   "update a single config value",
	Args:    cobra.ExactArgs(2),
	Aliases: []string{"cnf"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		val := args[1]

		color.HiBlue(fmt.Sprintf("setting %s to %s from %s", key, val, viper.Get(key)))
		config.Set(key, val)

		err := viper.WriteConfig()
		if err != nil {
			utils.ExitWithError(err)
		}

		color.HiGreen(fmt.Sprintf("successfully updated %s to %s!", key, val))
	},
}

var updateSelfSubCmd = &cobra.Command{
	Use:     "yourself",
	Aliases: []string{"yrself"},
	Short:   "check if pls has an available update and install it if so",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("checking if i have any updates...")
	},
}

func init() {
	updateCmd.AddCommand(updateCfgSubCmd)
	updateCmd.AddCommand(updateSelfSubCmd)
	updateCmd.AddCommand(updateSingleCfgSubCmd)
}
