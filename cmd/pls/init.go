package pls

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags
var (
	// Verbose is whether to return a verbose output
	Verbose bool

	cfgFile string
)

func initGlobalFlags() {
	// config flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pls/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	// persistent flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func addTopLevelSubcommands() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(tryCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(updateCmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	initGlobalFlags()
	addTopLevelSubcommands()
}
