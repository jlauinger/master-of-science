package cmd

import (
	"data-aquisition/ast"
	"github.com/spf13/cobra"
)

var analyzeAstCmd = &cobra.Command{
	Use:   "ast",
	Short: "Analyze abstract syntax tree",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ast.AnalyzeAst()
	},
}

func init() {
	RootCmd.AddCommand(analyzeAstCmd)
}

