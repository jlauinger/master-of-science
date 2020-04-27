package ast

import (
	"data-aquisition/lexical"
	"fmt"
)

func AnalyzeAst(offset, length int, dataDir string, skipProjects []string) {
	astFindingsFilename := fmt.Sprintf("%s/ast/ast_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	functionsFilename := fmt.Sprintf("%s/ast/functions_%d_%d.csv", dataDir, offset, offset + length - 1)
	statementsFilename := fmt.Sprintf("%s/ast/statements_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/ast/errors_ast_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := openAstFindingsFile(astFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := openFunctionsFile(functionsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := openStatementsFile(statementsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := openErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer closeFiles()

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Println("reading projects data...")
	projects, err := lexical.ReadProjects(projectsFilename)
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

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.Name)

		err := analyzeProject(project)

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
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

func analyzeProject(project *lexical.ProjectData) error {
	packages, err := lexical.GetProjectPackages(project)
	if err != nil {
		return err
	}

	fmt.Println("  parsing files...")

	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)

			AnalyzeAstSingleFile("save", fullFilename, pkg)
		}
	}

	return nil
}