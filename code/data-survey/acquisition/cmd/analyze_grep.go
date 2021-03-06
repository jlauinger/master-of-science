package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/linters"
)

var analyzeGrepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Extracts unsafe code fragments using grep",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the grep analysis operation
		linters.AnalyzeGrep(offset, length, dataDir, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(analyzeGrepCmd)
}
