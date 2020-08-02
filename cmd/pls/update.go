package pls

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var noUI bool

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "update various resources like configs and pls itself",
}

var updateCfgSubCmd = &cobra.Command{
	Use:     "configs",
	Short:   "update your pls configs",
	Aliases: []string{"config", "cnf", "cnfs"},
	Run: func(cmd *cobra.Command, args []string) {
		if noUI {
			if len(args) != 2 {
				utils.ExitWithError("you must have two arguments, the key and value - `pls update configs [key] [value]`")
			} else {
				key := args[0]
				val := args[1]

				color.HiBlue(fmt.Sprintf("setting %s to %s from %s", key, val, viper.Get(key)))
				config.Set(key, val)

				err := config.ParseAndUpdate(viper.GetViper())
				if err != nil {
					utils.ExitWithError(err)
				}

				color.HiGreen(fmt.Sprintf("successfully updated %s to %s!", key, val))
			}
		} else {
			config.UpdatePrompt(viper.AllSettings())
		}
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
	updateCfgSubCmd.Flags().BoolVarP(&noUI, "raw", "r", false, "input as key value, skip the dropdown GUI")
	updateCmd.AddCommand(updateCfgSubCmd)
	updateCmd.AddCommand(updateSelfSubCmd)
}
