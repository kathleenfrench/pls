package pls

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version is the current version of pls
	Version = "master"
	// Commit is the current commit of pls
	Commit = "none"
	// Date is the date at compile
	Date = "unknown"
	// Builder is the user who compiled the pls binary
	Builder = "unknown"
	// Verbose is whether to return a verbose output
	Verbose bool

	// flags
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "pls",
	Short: "a helpful little CLI for the lazy ones...",
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "print the current version of pls",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pls version: %s\n", Version)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// config flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pls.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	// persistent flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
