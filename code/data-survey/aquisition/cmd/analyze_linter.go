package cmd

import (
	"data-aquisition/lexical"
	"github.com/spf13/cobra"
)

var analyzeLinterCmd = &cobra.Command{
	Use:   "linter",
	Short: "Runs linter on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lexical.AnalyzeLinter(offset, length, dataDir, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeLinterCmd)
}
