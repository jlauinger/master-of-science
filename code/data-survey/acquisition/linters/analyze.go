package linters

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

/**
 * entry point for the grep analysis operation
 */
func AnalyzeGrep(offset, length int, dataDir string, skipProjects []string) {
	// build the filenames for files that will be written based on the configuration
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	grepFindingsFilename := fmt.Sprintf("%s/lexical/grep_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_grep_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the files and later close them
	if err := base.OpenPackagesFile(packagesFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenGrepFindingsFile(grepFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	// run the analysis by calling the generic project analysis function with a callback for grep
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackGrep, false, true)
}

/**
 * entry point for the go vet analysis operation
 */
func AnalyzeVet(offset, length int, dataDir string, skipProjects []string) {
	// build the filenames for files that will be written based on the configuration
	vetFindingsFilename := fmt.Sprintf("%s/linters/vet_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/linters/errors_vet_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the files and later close them
	if err := base.OpenVetFindingsFile(vetFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	// run the analysis by calling the generic project analysis function with a callback for go vet
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackVet, false, true)
}

/**
 * entry point for the gosec analysis operation
 */
func AnalyzeGosec(offset, length int, dataDir string, skipProjects []string) {
	// build the filenames for files that will be written based on the configuration
	gosecFindingsFilename := fmt.Sprintf("%s/linters/gosec_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/linters/errors_gosec_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the files and later close them
	if err := base.OpenGosecFindingsFile(gosecFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	// run the analysis by calling the generic project analysis function with a callback for gosec
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackGosec, false, true)
}

/**
 * entry point for the go-safer analysis operation
 */
func AnalyzeGosafer(offset, length int, dataDir string, skipProjects []string) {
	// build the filenames for files that will be written based on the configuration
	gosaferFindingsFilename := fmt.Sprintf("%s/linters/gosafer_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/linters/errors_gosafer_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the files and later close them
	if err := base.OpenGosaferFindingsFile(gosaferFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := base.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer base.CloseFiles()

	// run the analysis by calling the generic project analysis function with a callback for go-safer
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackGosafer, false, true)
}
