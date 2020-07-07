package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/eval2"
)

var Eval2Cmd = &cobra.Command{
	Use:   "eval2",
	Short: "run second wave evaluation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		eval2.Run(dataDir, offset, length)
	},
}

func init() {
	analyzeCmd.AddCommand(Eval2Cmd)
}
