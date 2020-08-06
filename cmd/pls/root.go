package pls

import (
	"fmt"
	"os"

	"github.com/kathleenfrench/pls/internal/config"
	"github.com/kathleenfrench/pls/internal/style"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// variables injected during build
var (
	// Version is the current version of pls
	Version = "master"
	// Commit is the current commit of pls
	Commit = "none"
	// Date is the date at compile
	Date = "unknown"
	// Builder is the user who compiled the pls binary
	Builder = "unknown"

	versionFlag bool

	plsCfg config.Settings
)

var rootCmd = &cobra.Command{
	Use:   "pls",
	Short: "a helpful little CLI that does things for you when you ask nice...",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			style.PrintBanner()
			printVersion()
		} else {
			cmd.Usage()
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		s, err := config.Parse(viper.GetViper())
		if err != nil {
			utils.ExitWithError(err)
		}

		plsCfg = s
	},
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
