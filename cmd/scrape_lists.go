package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ellis-vester/lb-scrape/internal/tui"
)

var ListsPath string
var OutputDir string
var PollInterval int

var scrapeListsCmd = &cobra.Command{
	Use:   "scrape-lists",
	Short: "Scrape the provided Letterboxd lists.",
	Long: `Scrape and aggregate the provided letterboxd lists, 
			outputting the results to two CSV files.`,
	Run: func(cmd *cobra.Command, args []string) {

		model := tui.NewScrapeListsModel(ListsPath, OutputDir, PollInterval)

		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Oh no!", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeListsCmd)

	scrapeListsCmd.PersistentFlags().StringVarP(
		&ListsPath,
		"lists-path",
		"l",
		"./lists.txt",
		"The path to the file containing the lists to scrape.")

	scrapeListsCmd.PersistentFlags().StringVarP(
		&OutputDir,
		"output-dir",
		"o",
		".",
		"The directory to output the CSV file to.")

	scrapeListsCmd.PersistentFlags().IntVarP(
		&PollInterval,
		"poll-interval",
		"p",
		5,
		"The seconds to wait between requests to Letterboxd.")
}
