package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/projects"
)

// CLI configuration variables
var datasetSize int
var download, createForks bool
var downloadDir, accessToken string

var GetProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Gets projects from Github and populates projects.csv",
	Long:  `Can also download the repositories itself, as well as fork them into a new organization`,
	Run: func(cmd *cobra.Command, args []string) {
		// call the get projects operation with the respective CLI parameters
		projects.GetProjects(datasetSize, dataDir, downloadDir, download, createForks, accessToken)
	},
}

func init() {
	// register the command
	RootCmd.AddCommand(GetProjectsCmd)

	// register available CLI parameters
	GetProjectsCmd.Flags().IntVar(&datasetSize, "size", 500, "Number of projects")
	GetProjectsCmd.Flags().BoolVar(&download, "download", false, "Download repositories")
	GetProjectsCmd.Flags().StringVar(&downloadDir, "destination", "", "Download destination")
	GetProjectsCmd.Flags().BoolVar(&createForks, "fork", false, "Fork repositories into a user account")
	GetProjectsCmd.Flags().StringVar(&accessToken, "access-token", "", "Github access token for fork target user")
}
