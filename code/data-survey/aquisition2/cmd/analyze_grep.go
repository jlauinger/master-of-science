package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var analyzeGrepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Extracts unsafe code fragments using grep",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("grepping for profit now! From %d to %d", offset, length)
	},
}

func init() {
	analyzeCmd.AddCommand(analyzeGrepCmd)
}
