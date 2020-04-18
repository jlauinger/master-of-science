package cmd

import (
	"data-aquisition/analysis"
	"github.com/spf13/cobra"
)

var analyzeGrepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Extracts unsafe code fragments using grep",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		analysis.AnalyzeGrep(offset, length, dataDir, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGrepCmd)
}
