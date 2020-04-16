package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	offset, _ := strconv.Atoi(os.Args[1])
	length, _ := strconv.Atoi(os.Args[2])
	pathPrefix := os.Args[3]

	projectsFilename := fmt.Sprintf("%s/projects.csv", pathPrefix)
	modulesFilename := fmt.Sprintf("%s/modules_%d_%d.csv", pathPrefix, offset, offset + length - 1)
	matchesFilename := fmt.Sprintf("%s/unsafe_matches_%d_%d.csv", pathPrefix, offset, offset + length - 1)
	vetResultsFilename := fmt.Sprintf("%s/vet_results_%d_%d.csv", pathPrefix, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/errors_%d_%d.csv", pathPrefix, offset, offset + length - 1)

	defer closeFiles()
	if err := openFiles(modulesFilename, matchesFilename, vetResultsFilename, errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	projects, err := ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	for _, project := range projects[offset:offset+length] {
		if !goModExists(project) {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "go.mod",
				ProjectName:      project.ProjectName,
				ModuleImportPath: "",
				FileName:         "",
				Message:          "go.mod does not exist",
			})
		}

		err := analyzeProject(project)

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "project",
				ProjectName:      project.ProjectName,
				ModuleImportPath: "",
				FileName:         "",
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}