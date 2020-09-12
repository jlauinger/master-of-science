package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/ast"
)

var analyzeAstCmd = &cobra.Command{
	Use:   "ast",
	Short: "Saves AST analysis results",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the AST analysis operation
		ast.AnalyzeAst(dataDir, offset, length, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(analyzeAstCmd)
}
