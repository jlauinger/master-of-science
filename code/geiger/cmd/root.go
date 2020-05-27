package cmd

import (
	"fmt"
	"geiger/counter"
	"github.com/spf13/cobra"
	"os"
)

var maxIndent int
var shortenSeenPackages, showStandardPackages, printLinkToPkgGoDev, printUnsafeLines bool

var RootCmd = &cobra.Command{
	Use:   "geiger",
	Short: "Counts unsafe usages in dependencies",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code`,
	Args: cobra.RangeArgs(0, 1000),
	Run: func(cmd *cobra.Command, args []string) {
		counter.Run(counter.Config{
			MaxIndent:            maxIndent,
			ShortenSeenPackages:  shortenSeenPackages,
			ShowStandardPackages: showStandardPackages,
			PrintLinkToPkgGoDev:  printLinkToPkgGoDev,
			PrintUnsafeLines:     printUnsafeLines,
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
	RootCmd.PersistentFlags().BoolVar(&printLinkToPkgGoDev, "link", false, "Print link to pkg.go.dev instead of package name")
	RootCmd.PersistentFlags().BoolVar(&printUnsafeLines, "show-code", false, "Print the code lines with unsafe usage")
}