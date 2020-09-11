package geiger

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
)

func analyzeProject(project *base.ProjectData) error {
	packages, err := GetProjectPackages(project)
	if err != nil {
		return err
	}

	fullFilenames := make([]string, 0, 500)
	fileToPackageMap := map[string]*base.PackageData{}

	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			fullFilenames = append(fullFilenames, fullFilename)
			fileToPackageMap[fullFilename] = pkg
		}
	}

	fileToLineCountMap, err := countLines(fullFilenames)
	if err != nil {
		return err
	}
	fileToByteCountMap, err := countBytes(fullFilenames)
	if err != nil {
		return err
	}
	fillPackageLOC(packages, fileToLineCountMap, fileToByteCountMap)

	analyzeHopCounts(packages)

	geigerPackages(project, packages, fileToLineCountMap, fileToByteCountMap)

	writePackages(packages)

	projectSum := 0
	for _, pkg := range packages {
		projectSum += pkg.UnsafeSum
	}
	fmt.Printf("  unsafe sum: %d\n", projectSum)

	return nil
}