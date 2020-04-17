package cmd

import (
	"github.com/spf13/cobra"

	"data-aquisition/projects"
)

var download bool
var downloadDir string

var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Gets projects from Github and populates projects.csv",
	Long:  `Can also download the repositories itself`,
	Run: func(cmd *cobra.Command, args []string) {
		projects.GetProjects(dataDir, downloadDir, download)
	},
}

func init() {
	RootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.Flags().BoolVar(&download, "download", false, "Download repositories")
	getProjectsCmd.Flags().StringVar(&downloadDir, "destination", "", "Download destination")
}
