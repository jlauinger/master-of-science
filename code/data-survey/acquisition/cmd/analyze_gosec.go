package cmd

import (
	"data-acquisition/lexical"
	"github.com/spf13/cobra"
)

var analyzeGosecCmd = &cobra.Command{
	Use:   "gosec",
	Short: "Runs gosec on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lexical.AnalyzeGosec(offset, length, dataDir, skipProjects, doCopy, copyDestination)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGosecCmd)
}
