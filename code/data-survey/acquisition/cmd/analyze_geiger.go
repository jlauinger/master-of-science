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
		geiger.Run(dataDir, offset, length, skipProjects)
	},
}

func init() {
	analyzeCmd.AddCommand(GeigerCmd)
}
