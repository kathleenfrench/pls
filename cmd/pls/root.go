package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "master"
	commit  = "none"
	date    = "unknown"
	builder = "unknown"
)

var root = &cobra.Command{
	Use:   "pls",
	Short: "pls is a helpful little cli",
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
