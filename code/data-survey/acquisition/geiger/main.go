package geiger

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

/**
 * entry point for the go-geiger analysis operation
 */
func Run(dataDir string, offset, length int, skipProjects []string) {
	// build the filenames for the files that will be written
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	geigerFilename := fmt.Sprintf("%s/geiger/geiger_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/geiger/errors_geiger_%d_%d.csv", dataDir, offset, offset + length - 1)

	// open the files and close them later
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

	// run the analysis by calling the generic project analysis function with a callback for go-geiger
	base.AnalyzeProjects(dataDir, offset, length, skipProjects, callbackGeiger, true, false)
}

/**
 * callback handling the go-geiger analysis coordination. This is called for each project
 */
func callbackGeiger(project *base.ProjectData, packages []*base.PackageData, _ map[string]*base.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	geigerPackages(project, packages, fileToLineCountMap, fileToByteCountMap)
}