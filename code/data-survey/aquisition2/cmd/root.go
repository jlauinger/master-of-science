package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

func init() {}