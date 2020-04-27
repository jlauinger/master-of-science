package cmd

import (
	"data-aquisition/lexical"
	"github.com/spf13/cobra"
)

var analyzeVetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Runs go vet on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lexical.AnalyzeVet(offset, length, dataDir, skipProjects, doCopy, copyDestination)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeVetCmd)
}
