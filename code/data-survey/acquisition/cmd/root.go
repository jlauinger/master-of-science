package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// this represents the data directory for all data CSV files used by the tool (write/read)
var dataDir string

var RootCmd = &cobra.Command{
	Use:   "data-acquisition",
	Short: "Data acquisition tool",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Data Acquisition Tool: Hello World!\n")
	},
}

/**
 * this is the main entry point called by main.go
 */
func Execute() {
	// execute the root command using Cobra
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// add available command line options
	RootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "directory for CSV data files")
}