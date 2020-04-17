package analysis

import (
	"fmt"
)

func AnalyzeGrep(offset, length int, dataDir string) {
	commonAnalysis(offset, length, dataDir, operatorGrepAnalysis)
}
func operatorGrepAnalysis(project *ProjectData, modules []ModuleData, fileToModuleMap map[string]ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {

	parsedGrepLines, err := grepForUnsafe(modules)
	if err != nil {
		_ = WriteErrorCondition(ErrorConditionData{
			Stage:            "grep-parse",
			ProjectName:      project.ProjectName,
			ModuleImportPath: "",
			FileName:         "",
			Message:          err.Error(),
		})
		fmt.Println("SAVING ERROR!")
	}
	analyzeGrepLines(parsedGrepLines, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)
}


func AnalyzeVet(offset, length int, dataDir string) {
	commonAnalysis(offset, length, dataDir, operatorVetAnalysis)
}
func operatorVetAnalysis(project *ProjectData, modules []ModuleData, fileToModuleMap map[string]ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {

	vetFindings := runVet(project, modules)
	analyzeVetFindings(vetFindings, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)
}


func AnalyzeGosec(offset, length int, dataDir string) {
	commonAnalysis(offset, length, dataDir, operatorGosecAnalysis)
}
func operatorGosecAnalysis(project *ProjectData, modules []ModuleData, fileToModuleMap map[string]ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {

	gosecFindings, _ := runGosec(project, modules)
	analyzeGosecFindings(gosecFindings, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)
}


func commonAnalysis(offset, length int, dataDir string,
	operator func(*ProjectData, []ModuleData, map[string]ModuleData, map[string]int, map[string]int)) {

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)
	modulesFilename := fmt.Sprintf("%s/modules_%d_%d.csv", dataDir, offset, offset + length - 1)
	matchesFilename := fmt.Sprintf("%s/unsafe_matches_%d_%d.csv", dataDir, offset, offset + length - 1)
	vetResultsFilename := fmt.Sprintf("%s/vet_results_%d_%d.csv", dataDir, offset, offset + length - 1)
	gosecResultsFilename := fmt.Sprintf("%s/gosec_results_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/errors_%d_%d.csv", dataDir, offset, offset + length - 1)

	defer closeFiles()
	if err := openFiles(modulesFilename, matchesFilename, vetResultsFilename,
		gosecResultsFilename, errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	projects, err := ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	for projectIdx, project := range projects[offset:offset+length] {
		if !goModExists(project) {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "go.mod",
				ProjectName:      project.ProjectName,
				ModuleImportPath: "",
				FileName:         "",
				Message:          "go.mod does not exist",
			})
			fmt.Printf("%d/%d (#%d): Skipping %s\n", projectIdx+1-offset, length, projectIdx+1, project.ProjectName)
			continue
		}

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.ProjectName)

		err := analyzeProject(project, operator)

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "project",
				ProjectName:      project.ProjectName,
				ModuleImportPath: "",
				FileName:         "",
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}