package pls

import (
	"github.com/kathleenfrench/pls/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Verbose is whether to return a verbose output
	Verbose bool
	cfgFile string
)

func initGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pls/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "verbose output")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "see the current version of pls")
}

func addTopLevelSubcommands() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(tryCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(addSubCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(getCmd)
}

func setPlsStyling() {
	rootCmd.SetUsageTemplate(stylizePls())
}

func init() {
	cobra.OnInitialize(config.Initialize)
	initGlobalFlags()
	addTopLevelSubcommands()
	setPlsStyling()
}
