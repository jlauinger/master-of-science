package eval2

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
)

func analyzeProject(project *lexical.ProjectData) error {
	packages, err := GetProjectPackages(project)
	if err != nil {
		return err
	}

	fullFilenames := make([]string, 0, 500)
	fileToPackageMap := map[string]*lexical.PackageData{}

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

	analyzeDepTree(packages)

	geigerPackages(project, packages)

	projectUnsafeSum := 0
	for _, pkg := range packages {
		projectUnsafeSum += pkg.UnsafeSum
	}
	fmt.Printf("total project unsafe sum: %d\n", projectUnsafeSum)

	writePackages(packages)

	return nil
}