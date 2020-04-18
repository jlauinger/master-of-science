package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var offset, length int
var skipProjects []string

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Runs the analysis",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use grep, vet, or gosec")
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)

	analyzeCmd.PersistentFlags().IntVar(&offset, "offset", 0, "parallelization: projects slicing offset")
	analyzeCmd.PersistentFlags().IntVar(&length, "length", 500, "parallelization: projects slicing length")
	analyzeCmd.PersistentFlags().StringArrayVar(&skipProjects, "skip", []string{}, "skip these project names, e.g golang/go")
}
