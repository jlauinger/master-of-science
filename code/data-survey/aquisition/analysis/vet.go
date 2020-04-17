package analysis

import (
	"data-aquisition/base"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runVet(project *base.ProjectData, modules []base.ModuleData) []base.VetFindingLine {
	packagePaths := make([]string, len(modules))

	for i, module := range modules {
		packagePaths[i] = module.ModuleImportPath
	}

	args := []string{"vet", "-c=0"}
	args = append(args, packagePaths...)

	cmd := exec.Command("go", args...)
	cmd.Dir = project.ProjectCheckoutPath

	vetOutput, _ := cmd.CombinedOutput()

	vetLines := strings.Split(string(vetOutput), "\n")
	vetFindings := make([]base.VetFindingLine, 0)

	for i := 0; i < len(vetLines); i++ {
		messageLine := vetLines[i]

		if len(messageLine) <= 0 || messageLine[0] == '#' {
			continue
		}

		var contextLines []string

		for {
			contextLine := vetLines[i+1]

			components := strings.Split(contextLine, "\t")

			if len(components) <= 1 {
				break
			}

			_, err := strconv.Atoi(components[0])
			if err != nil {
				break
			}

			contextLines = append(contextLines, contextLine)
			i++
		}

		vetFindings = append(vetFindings, base.VetFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return vetFindings
}

func analyzeVetFindings(vetFindings []base.VetFindingLine, fileToModuleMap map[string]base.ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {
	for _, line := range vetFindings {
		components := strings.Split(line.Message, ":")

		var fullFilename string
		var lineNumber int
		var column int
		var message string

		if components[0] == "vet" {
			components = components[1:]
		}

		fullFilename = strings.Trim(components[0], " ")
		module := fileToModuleMap[fullFilename]
		filename := fullFilename[len(module.PackageDir)+1:]

		if strings.Contains(filename, "test") {
			continue
		}

		if len(components) < 2 {
			fmt.Println(line.Message)
			os.Exit(1)
		}

		lineNumber, err := strconv.Atoi(components[1])
		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "vet-parse-linenumber",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         filename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		column, err = strconv.Atoi(components[2])
		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "vet-parse-column",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         filename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		message = strings.Trim(components[3], " ")

		err = WriteVetFinding(base.VetFindingData{
			ProjectName:          module.ProjectName,
			ModuleImportPath:     module.ModuleImportPath,
			ModuleRegistry:       module.ModuleRegistry,
			ModuleVersion:        module.ModuleVersion,
			ModuleNumberGoFiles:  module.ModuleNumberGoFiles,
			ModuleCheckoutFolder: module.PackageDir,
			FileName:             filename,
			FileSizeBytes:        fileToByteCountMap[fullFilename],
			FileSizeLines:        fileToLineCountMap[fullFilename],
			FileImportsUnsafePkg: false,
			FileGoVetOutput:      "",
			LineNumber:           lineNumber,
			Column:               column,
			Message:              message,
			RawOutput:            line.Message,
		})

		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "vet-write",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         filename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}