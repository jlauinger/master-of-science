package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// CLI parameter variables for all analyze commands
var offset, length int
var skipProjects []string

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Runs the analysis",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use geiger, grep, vet, gosec, gosafer, or ast")
	},
}

func init() {
	// register the command
	RootCmd.AddCommand(analyzeCmd)

	// register available CLI parameters
	analyzeCmd.PersistentFlags().IntVar(&offset, "offset", 0, "parallelization: projects slicing offset")
	analyzeCmd.PersistentFlags().IntVar(&length, "length", 500, "parallelization: projects slicing length")
	analyzeCmd.PersistentFlags().StringArrayVar(&skipProjects, "skip", []string{}, "skip these project names, e.g golang/go")
}
