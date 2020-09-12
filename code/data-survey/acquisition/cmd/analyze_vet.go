package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/linters"
)

var analyzeVetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Runs go vet on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the go vet analysis operation
		linters.AnalyzeVet(offset, length, dataDir, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(analyzeVetCmd)
}
