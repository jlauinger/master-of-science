package cmd

import (
	"fmt"
	"geiger/rewrite"
	"github.com/spf13/cobra"
	"os"
)

var maxIndent int
var shortenSeenPackages, showStandardPackages bool

var RootCmd = &cobra.Command{
	Use:   "geiger",
	Short: "Counts unsafe usages in dependencies",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code`,
	Args: cobra.RangeArgs(0, 1000),
	Run: func(cmd *cobra.Command, args []string) {
		rewrite.Run(rewrite.Config{
			MaxIndent:            maxIndent,
			ShortenSeenPackages:  shortenSeenPackages,
			ShowStandardPackages: showStandardPackages,
		}, args...)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVar(&maxIndent, "level", 10, "Maximum indent level")
	RootCmd.PersistentFlags().BoolVar(&shortenSeenPackages, "dnr", true, "Do not repeat packages")
	RootCmd.PersistentFlags().BoolVar(&showStandardPackages, "show-std", false, "Show Goland stdlib packages")
}