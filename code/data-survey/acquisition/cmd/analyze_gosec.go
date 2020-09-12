package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/linters"
)

var analyzeGosecCmd = &cobra.Command{
	Use:   "gosec",
	Short: "Runs gosec on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the gosec analysis step
		linters.AnalyzeGosec(offset, length, dataDir, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(analyzeGosecCmd)
}
