package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/ast"
)

// CLI parameter variables
var mode, filename string

var astCmd = &cobra.Command{
	Use:   "ast",
	Short: "Analyze abstract syntax tree",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the AST single analysis operation with the given parameters
		ast.AnalyzeAstSingleFile(mode, filename,  nil)
	},
}

func init() {
	// register the command
	RootCmd.AddCommand(astCmd)

	// register available CLI parameters
	astCmd.PersistentFlags().StringVar(&mode, "mode", "", "print mode (tree,func,stmt,save)")
	astCmd.PersistentFlags().StringVar(&filename, "file", "", "file to analyze")
}

