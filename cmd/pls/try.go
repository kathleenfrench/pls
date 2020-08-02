package pls

import (
	"fmt"

	"github.com/fatih/color"
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
		s, err := config.Parse(viper.GetViper())
		if err != nil {
			utils.ExitWithError(err)
		}

		s.UpdateSettings(viper.GetViper())

		color.HiGreen(fmt.Sprintf("parsed config: %v", s))
	},
}
