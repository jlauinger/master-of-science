package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/data-acquisition/projects"
)

var download, createForks bool
var downloadDir, accessToken string

var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Gets projects from Github and populates projects.csv",
	Long:  `Can also download the repositories itself, as well as fork them into a new organization`,
	Run: func(cmd *cobra.Command, args []string) {
		projects.GetProjects(dataDir, downloadDir, download, createForks, accessToken)
	},
}

func init() {
	RootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.Flags().BoolVar(&download, "download", false, "Download repositories")
	getProjectsCmd.Flags().StringVar(&downloadDir, "destination", "", "Download destination")
	getProjectsCmd.Flags().BoolVar(&createForks, "fork", false, "Fork repositories into a user account")
	getProjectsCmd.Flags().StringVar(&accessToken, "access-token", "", "Github access token for fork target user")
}
