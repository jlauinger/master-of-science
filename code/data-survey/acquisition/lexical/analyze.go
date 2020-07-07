package lexical

import (
	"fmt"
)

func AnalyzeGrep(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string) {
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	grepFindingsFilename := fmt.Sprintf("%s/lexical/grep_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_grep_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := OpenPackagesFile(packagesFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := openGrepFindingsFile(grepFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer CloseFiles()

	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, true, operatorGrepAnalysis)
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
	vetFindingsFilename := fmt.Sprintf("%s/lexical/vet_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_vet_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := openVetFindingsFile(vetFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer CloseFiles()

	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, false, operatorVetAnalysis)
}
func operatorVetAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	vetFindings := runVet(project, packages)
	return analyzeVetFindings(vetFindings, fileToPackageMap, fileToLineCountMap, fileToByteCountMap, project)
}


func AnalyzeGosec(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string) {
	gosecFindingsFilename := fmt.Sprintf("%s/lexical/gosec_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_gosec_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := openGosecFindingsFile(gosecFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer CloseFiles()

	commonAnalysis(offset, length, dataDir, skipProjects, doCopy, copyDestination, false, operatorGosecAnalysis)
}
func operatorGosecAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	gosecFindings, _ := runGosec(project, packages)
	return analyzeGosecFindings(gosecFindings, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}


func AnalyzeLinter(offset, length int, dataDir string, skipProjects []string) {
	linterFindingsFilename := fmt.Sprintf("%s/lexical/linter_findings_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_linter_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := openLinterFindingsFile(linterFindingsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer CloseFiles()

	commonAnalysis(offset, length, dataDir, skipProjects, false, "", false, operatorLinterAnalysis)
}
func operatorLinterAnalysis(project *ProjectData, packages []*PackageData, fileToPackageMap map[string]*PackageData,
	_, _ map[string]int) map[string]string {

	linterFindings := runLinter(project, packages)
	return analyzeLinterFindings(linterFindings, fileToPackageMap, project)
}


func commonAnalysis(offset, length int, dataDir string, skipProjects []string, doCopy bool, copyDestination string, writePackagesToFile bool,
	operator func(*ProjectData, []*PackageData, map[string]*PackageData, map[string]int, map[string]int) map[string]string) {

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Println("reading projects data...")
	projects, err := ReadProjects(projectsFilename)
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

		filesToCopy, err := analyzeProject(project, writePackagesToFile, operator)

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