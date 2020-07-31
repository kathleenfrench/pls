package cmd

import "fmt"

var (
	version = "master"
	commit  = "none"
	date    = "unknown"
	builder = "unknown"
)

func root = &cobra.Command{
	Use: "pls",
	Short: "pls is a helpful little cli",
}


func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}