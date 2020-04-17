package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var analyzeGosecCmd = &cobra.Command{
	Use:   "gosec",
	Short: "Runs gosec on the projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("goseccing like never before now! From %d to %d", offset, length)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGosecCmd)
}
