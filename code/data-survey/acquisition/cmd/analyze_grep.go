package cmd

import (
	"data-acquisition/lexical"
	"github.com/spf13/cobra"
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
