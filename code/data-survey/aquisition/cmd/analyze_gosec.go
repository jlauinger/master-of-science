package cmd

import (
	"data-aquisition/analysis"
	"github.com/spf13/cobra"
)

var analyzeGosecCmd = &cobra.Command{
	Use:   "gosec",
	Short: "Runs gosec on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		analysis.AnalyzeGosec(offset, length, dataDir, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGosecCmd)
}
