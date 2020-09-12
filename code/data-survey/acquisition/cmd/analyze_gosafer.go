package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/linters"
)

var analyzeGosaferCmd = &cobra.Command{
	Use:   "linter",
	Short: "Runs linter on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the go-safer analysis operation
		linters.AnalyzeGosafer(offset, length, dataDir, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(analyzeGosaferCmd)
}
