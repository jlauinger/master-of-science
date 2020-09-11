package ast

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

/**
 * this is the entry point for the AST analysis
 */
func AnalyzeAst(offset, length int, dataDir string, skipProjects []string) {
	// build the result files filenames based on the configured project indices
	astFindingsFilename := fmt.Sprintf("%s/ast/ast_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	functionsFilename := fmt.Sprintf("%s/ast/functions_%d_%d.csv", dataDir, offset, offset + length - 1)
	statementsFilename := fmt.Sprintf("%s/ast/statements_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/ast/errors_ast_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the needed result files and defer closing them again
	if err := base.OpenAstFindingsFile(astFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenAstFunctionsFile(functionsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenAstStatementsFile(statementsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	// build the projects data CSV filename
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)
	// and load the projects data
	fmt.Println("reading projects data...")
	projects, err := base.ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	// transform the list of projects to skip into a hash set for faster identification of projects later on
	skipProjectMap := make(map[string]struct{}, len(skipProjects))
	for _, skipProject := range skipProjects {
		skipProjectMap[skipProject] = struct{}{}
	}

	// go through the projects specified by the configured from / to indices
	for projectIdx, project := range projects[offset:offset+length] {
		// check if the project is contained in the skip map and if so, skip it
		if _, ok := skipProjectMap[project.Name]; ok {
			fmt.Printf("%d/%d (#%d): Skipping %s as requested\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.Name)

		// analyze the project and save potential errors
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

/**
 * analyzes a single project using the AST analysis
 */
func analyzeProject(project *base.ProjectData) error {
	// identify all relevant packages used by this project
	packages, err := base.GetProjectPackages(project)
	if err != nil {
		return err
	}

	fmt.Println("  parsing files...")

	// go through all packages
	for _, pkg := range packages {
		// for each package, go through all the Go files
		for _, file := range pkg.GoFiles {
			// build the effective file path
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			// and analyze it using the AST analysis
			AnalyzeAstSingleFile("save", fullFilename, pkg)
		}
	}

	// if we come here, there have been no errors
	return nil
}