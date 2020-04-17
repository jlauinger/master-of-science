package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Gets projects from Github and populates projects.csv",
	Long:  `Can also download the repositories itself`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get projects called")
	},
}

func init() {
	RootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.Flags().BoolP("download", "", false, "Download repositories")
}
