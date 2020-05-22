package cmd

import (
	"data-acquisition/ast"
	"github.com/spf13/cobra"
)

var analyzeAstCmd = &cobra.Command{
	Use:   "ast",
	Short: "Saves AST analysis results",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ast.AnalyzeAst(offset, length, dataDir, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeAstCmd)
}
