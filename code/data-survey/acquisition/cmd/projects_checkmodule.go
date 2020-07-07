package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/projects"
)

var projectsCheckModuleCmd = &cobra.Command{
	Use:   "checkmodule",
	Short: "For each project in projects.csv, checks if the project uses modules and what is the main module",
	Long:  `Repopulates and rewrites projects.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		projects.CheckModule(dataDir)
	},
}

func init() {
	GetProjectsCmd.AddCommand(projectsCheckModuleCmd)
}
