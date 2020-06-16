package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var dataDir string

var RootCmd = &cobra.Command{
	Use:   "data-aquisition",
	Short: "Data aquisition tool",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hello World\n")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "directory for CSV data files")
}