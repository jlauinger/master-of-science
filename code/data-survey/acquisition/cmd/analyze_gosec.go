package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
)

var analyzeGosecCmd = &cobra.Command{
	Use:   "gosec",
	Short: "Runs gosec on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		linters.AnalyzeGosec(offset, length, dataDir, skipProjects, doCopy, copyDestination)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGosecCmd)
}
