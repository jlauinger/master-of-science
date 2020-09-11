package geiger

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

func Run(dataDir string, offset, length int, skipProjects []string) {
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	geigerFilename := fmt.Sprintf("%s/geiger/geiger_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/geiger/errors_geiger_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := base.OpenPackagesFile(packagesFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenGeigerFindingsFile(geigerFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Println("reading projects data...")
	projects, err := base.ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	skipProjectMap := make(map[string]struct{}, len(skipProjects))
	for _, skipProject := range skipProjects {
		skipProjectMap[skipProject] = struct{}{}
	}

	for projectIdx, project := range projects[offset:offset+length] {
		if _, ok := skipProjectMap[project.Name]; ok {
			fmt.Printf("%d/%d (#%d): Skipping %s as requested\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		if !project.UsesModules {
			fmt.Printf("%d/%d (#%d): Skipping %s because it does not use modules\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.Name)

		err := analyzeProject(project)
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "project",
				ProjectName:       project.Name,
				PackageImportPath: "",
				FileName:          "",
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}