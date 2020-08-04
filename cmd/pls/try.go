package pls

import (
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tryCmd = &cobra.Command{
	Use:     "try",
	Short:   "try to do something",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ParseAndUpdate(viper.GetViper())
		if err != nil {
			utils.ExitWithError(err)
		}
	},
}
