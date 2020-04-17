package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var analyzeVetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Runs go vet on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vetting now! From %d to %d", offset, length)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeVetCmd)
}
