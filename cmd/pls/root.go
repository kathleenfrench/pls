package pls

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
)

var rootCmd = &cobra.Command{
	Use:   "pls",
	Short: "a helpful little CLI that does things for you when you ask nice...",
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
