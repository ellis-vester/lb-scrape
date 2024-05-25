package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool
var Debug bool

var rootCmd = &cobra.Command{
	Use:   "lbs",
	Short: "A tool for scraping and aggregating data on letterboxd.com.",
	Long:  `.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP(
		"toggle",
		"t",
		false,
		"Help message for toggle")

	rootCmd.PersistentFlags().BoolVarP(
		&Verbose,
		"verbose",
		"v",
		false,
		"Display more verbose output in console output.")
	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup(("verbose")))
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().BoolVarP(
		&Debug,
		"debug",
		"d",
		false,
		"Display debug output in console output.")
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup(("debug")))
	if err != nil {
		panic(err)
	}
}
