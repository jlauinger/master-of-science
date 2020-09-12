package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/geiger"
)

var GeigerCmd = &cobra.Command{
	Use:   "geiger",
	Short: "run evaluation with go-geiger implementation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// run the go-geiger analysis operation
		geiger.Run(dataDir, offset, length, skipProjects)
	},
}

func init() {
	// register the command
	analyzeCmd.AddCommand(GeigerCmd)
}
