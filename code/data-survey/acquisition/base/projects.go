package base

import (
	"fmt"
)

/**
 * type definition for a callback function handling the specific analysis operation
 */
type AnalysisCallback func(*ProjectData, []*PackageData, map[string]*PackageData, map[string]int, map[string]int)

/*
 * magic number to indicate no length given by the user
 */
const NoLengthGiven = 99999

/**
 * runs an analysis callback for all requested projects. This is a common function to all analysis operations
 */
func AnalyzeProjects(dataDir string, offset, length int, skipProjects []string, callback AnalysisCallback, doWritePackages, ignoreModules bool) {
	// build the filename for the projects CSV data
	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	// load the project data
	fmt.Println("reading projects data...")
	projects, err := ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	// transform the list of projects that should be skipped into a hash set for efficient lookup later on
	skipProjectMap := make(map[string]struct{}, len(skipProjects))
	for _, skipProject := range skipProjects {
		skipProjectMap[skipProject] = struct{}{}
	}

	// set the correct length if none is given
	if length == NoLengthGiven {
		length = len(projects)
	}

	// go through all projects in the slice determined by the start / end indices
	for projectIdx, project := range projects[offset:offset+length] {
		// check if the project should be skipped and do accordingly
		if _, ok := skipProjectMap[project.Name]; ok {
			fmt.Printf("%d/%d (#%d): Skipping %s as requested\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		// check if the project uses modules, and skip it if it doesn't unless modules should be ignored. The go-geiger
		// analysis needs to only analyze projects with module support, but the other analysis operations do not care
		if !project.UsesModules && !ignoreModules {
			fmt.Printf("%d/%d (#%d): Skipping %s because it does not use modules\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.Name)

		// analyze the project with the callback function and save any possible error
		err := analyzeProject(project, doWritePackages, callback)
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

/**
 * analyzes a single project with a generic callback function handling the specific analysis
 */
func analyzeProject(project *ProjectData, doWritePackages bool, callback AnalysisCallback) error {
	// identify packages in the dependency tree of this project
	packages, err := getProjectPackages(project)
	if err != nil {
		return err
	}

	// initialize a list of filenames for all package files and a hash map that can be used to map a specific file to
	// its package
	fullFilenames := make([]string, 0, 500)
	fileToPackageMap := map[string]*PackageData{}

	// to fill them, go through all packages
	for _, pkg := range packages {
		// then go through all go files in the package
		for _, file := range pkg.GoFiles {
			// build the full filename for the source file
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			// and insert it into the data structures
			fullFilenames = append(fullFilenames, fullFilename)
			fileToPackageMap[fullFilename] = pkg
		}
	}

	// count lines and bytes using wc for all files and weave the LOC information into the packages list
	fileToLineCountMap, err := countLines(fullFilenames)
	if err != nil {
		return err
	}
	fileToByteCountMap, err := countBytes(fullFilenames)
	if err != nil {
		return err
	}
	fillPackageLOC(packages, fileToLineCountMap, fileToByteCountMap)

	// fill the dependency depth / hop count information in the packages list
	analyzeHopCounts(packages)

	// do the specific analysis that was requested by calling the callback
	callback(project, packages, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)

	// if the packages file should be written, do accordingly
	if doWritePackages {
		writePackages(packages)
	}

	// if we make it until here, there have been no errors
	return nil
}