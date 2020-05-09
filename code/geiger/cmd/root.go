package cmd

import (
	"fmt"
	"geiger/facts"
	"geiger/prettyprint"
	"geiger/tools/go/analysis/singlechecker"
	"geiger/unsafecountpass"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var maxIndent int
var shortenSeen bool

var RootCmd = &cobra.Command{
	Use:   "geiger",
	Short: "Counts unsafe usages in dependencies",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code`,
	Run: func(cmd *cobra.Command, args []string) {
		facts.Init()
		singlechecker.Run(unsafecountpass.Analyzer)
		results := facts.GetAll()
		prettyprint.Print(results, maxIndent, shortenSeen)
	},
}

func Execute() {
	maxIndentValue, err := strconv.Atoi(os.Getenv("GEIGER_LEVEL"))
	if err == nil {
		maxIndent = maxIndentValue
	} else {
		maxIndent = 2
	}

	shortenSeenValue, err := strconv.ParseBool(os.Getenv("GEIGER_SHORTEN_SEEN"))
	if err == nil {
		shortenSeen = shortenSeenValue
	} else {
		shortenSeen = true
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}