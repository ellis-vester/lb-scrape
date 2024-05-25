package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Version string = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:     "lbs",
	Short:   "A tool for scraping and aggregating data from letterboxd.com.",
	Long:    `.`,
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
