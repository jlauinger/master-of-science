package projects

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
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

	for _, project := range projects {
		cmd := exec.Command("go", "list", "-m")
		cmd.Dir = project.CheckoutPath

		var rootModule string
		var usesModules bool

		output, err := cmd.Output()
		if err == nil {
			rootModule = strings.TrimSuffix(string(output), "\n")
			usesModules = true
		} else {
			rootModule = ""
			usesModules = false
		}

		fmt.Printf("%s: %s (%s)\n", project.Name, strconv.FormatBool(usesModules), rootModule)
	}
}