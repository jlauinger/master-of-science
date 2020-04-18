package analysis

import (
	"fmt"
)

func AnalyzeGrep(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string) {
	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, operatorGrepAnalysis)
}
func operatorGrepAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	parsedGrepLines, err := grepForUnsafe(packages)
	if err != nil {
		_ = WriteErrorCondition(ErrorConditionData{
			Stage:             "grep-parse",
			ProjectName:       project.Name,
			PackageImportPath: "",
			FileName:          "",
			Message:           err.Error(),
		})
		fmt.Println("SAVING ERROR!")
	}
	return analyzeGrepLines(parsedGrepLines, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}


func AnalyzeVet(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string) {
	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, operatorVetAnalysis)
}
func operatorVetAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	vetFindings := runVet(project, packages)
	return analyzeVetFindings(vetFindings, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}


func AnalyzeGosec(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string) {
	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, operatorGosecAnalysis)
}
func operatorGosecAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	gosecFindings, _ := runGosec(project, packages)
	return analyzeGosecFindings(gosecFindings, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}


func commonAnalysis(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string,
	operator func(*ProjectData, []*PackageData, map[string]*PackageData, map[string]int, map[string]int) map[string]string) {

	// TODO: open only the ones needed for current analysis

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	grepFindingsFilename := fmt.Sprintf("%s/grep_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	vetFindingsFilename := fmt.Sprintf("%s/vet_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	gosecFindingsFilename := fmt.Sprintf("%s/gosec_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/errors_%d_%d.csv", dataDir, offset, offset + length - 1)

	defer closeFiles()
	if err := openFiles(packagesFilename, grepFindingsFilename, vetFindingsFilename,
		gosecFindingsFilename, errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	fmt.Println("reading projects data...")
	projects, err := readProjects(projectsFilename)
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

		filesToCopy, err := analyzeProject(project, operator)

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

		if doCopy {
			copyFiles(filesToCopy, copyDestination)
		}
	}
}