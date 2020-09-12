package ast

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

/**
 * this is the entry point for the AST analysis
 */
func AnalyzeAst(dataDir string, offset, length int, skipProjects []string) {
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

	// run the analysis by calling the generic project analysis function with a callback for the AST analysis
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackAst, false, true)
}

/**
 * callback handling the AST analysis coordination. This is called for each project
 */
func callbackAst(_ *base.ProjectData, packages []*base.PackageData, _ map[string]*base.PackageData,	_, _ map[string]int) {

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
}