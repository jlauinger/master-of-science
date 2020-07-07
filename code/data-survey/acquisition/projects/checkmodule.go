package projects

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CheckModule(dataDir string) {
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Println("reading projects data...")
	projects, err := lexical.ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	headerWritten := false
	projectsFile, err := os.OpenFile(projectsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer projectsFile.Close()

	for _, project := range projects {
		cmd := exec.Command("go", "list", "-m")
		cmd.Dir = project.CheckoutPath

		output, err := cmd.Output()
		if err == nil {
			project.RootModule = strings.TrimSuffix(string(output), "\n")
			project.UsesModules = true
		} else {
			project.RootModule = ""
			project.UsesModules = false
		}

		fmt.Printf("%s: %s (%s)\n", project.Name, strconv.FormatBool(project.UsesModules), project.RootModule)

		if headerWritten {
			_ = gocsv.MarshalWithoutHeaders([]lexical.ProjectData{*project}, projectsFile)
		} else {
			headerWritten = true
			_ = gocsv.Marshal([]lexical.ProjectData{*project}, projectsFile)
		}
	}
}