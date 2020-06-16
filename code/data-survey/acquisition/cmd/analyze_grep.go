package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/data-acquisition/lexical"
)

var analyzeGrepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Extracts unsafe code fragments using grep",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lexical.AnalyzeGrep(offset, length, dataDir, skipProjects, doCopy, copyDestination)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGrepCmd)
}
