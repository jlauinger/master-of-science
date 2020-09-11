package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
)

var analyzeLinterCmd = &cobra.Command{
	Use:   "linter",
	Short: "Runs linter on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		linters.AnalyzeLinter(offset, length, dataDir, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeLinterCmd)
}
