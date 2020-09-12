package projects

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"os/exec"
	"strconv"
	"strings"
)

/**
 * checks each project for module support, identifies the main module and updates the projects.csv file
 */
func CheckModule(dataDir string) {
	// build the projects filename from the configuration
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	// load the project CSV data
	fmt.Println("reading projects data...")
	projects, err := base.ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	// open the projects file for writing as it will be updated with the module information
	if err := base.OpenProjectsFile(projectsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer base.CloseFiles()

	// go through all projects
	for _, project := range projects {
		// build the go list -m command in the project checkout directory to identify module support
		cmd := exec.Command("go", "list", "-m")
		cmd.Dir = project.CheckoutPath

		// run the command and check both output and exit code
		output, err := cmd.Output()
		if err == nil {
			// if the exit code is 0 (no error), the project uses modules, and the output of the command is the main /
			// root module
			project.RootModule = strings.TrimSuffix(string(output), "\n")
			project.UsesModules = true
		} else {
			// otherwise, the project does not use modules and I use no-mod as a placeholder
			project.RootModule = "no-mod"
			project.UsesModules = false
		}

		fmt.Printf("%s: %s (%s)\n", project.Name, strconv.FormatBool(project.UsesModules), project.RootModule)

		// write the updated project back into the CSV file
		err = base.WriteProject(*project)
		if err != nil {
			panic(err)
		}
	}
}