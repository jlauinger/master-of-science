package analysis

import (
	"bytes"
	"data-aquisition/base"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runGosec(project *base.ProjectData, modules []base.ModuleData) ([]base.GosecIssueOutput, error) {
	packagePaths := make([]string, 0, 1000)

	for _, module := range modules {
		packagePaths = append(packagePaths, module.PackageDir)
	}

	args := []string{"-quiet", "-no-fail", "-fmt=json"}
	args = append(args, packagePaths...)

	cmd := exec.Command("gosec", args...)
	cmd.Dir = project.ProjectCheckoutPath

	gosecOutput, _ := cmd.CombinedOutput()

	dec := json.NewDecoder(bytes.NewReader(gosecOutput))
	var gosecResult base.GosecOutput

	err := dec.Decode(&gosecResult)
	if err != nil {
		_ = WriteErrorCondition(base.ErrorConditionData{
			Stage:            "gosec-parse",
			ProjectName:      project.ProjectName,
			ModuleImportPath: "",
			FileName:         "",
			Message:          err.Error(),
		})
		fmt.Println("SAVING ERROR!")
		return nil, err
	}

	return gosecResult.Issues, nil
}

func analyzeGosecFindings(gosecFindings []base.GosecIssueOutput, fileToModuleMap map[string]base.ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {
	for _, line := range gosecFindings {
		module := fileToModuleMap[line.File]
		shortFilename := line.File[len(module.PackageDir)+1:]

		var lineNumberText string
		if strings.Contains(line.Line, "-") {
			lineNumberText = strings.Split(line.Line, "-")[0]
		} else {
			lineNumberText = line.Line
		}
		lineNumber, err := strconv.Atoi(lineNumberText)
		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "gosec-parse-linenumber",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         shortFilename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		var columnText string
		if strings.Contains(line.Column, "-") {
			columnText = strings.Split(line.Column, "-")[0]
		} else {
			columnText = line.Column
		}
		column, err := strconv.Atoi(columnText)
		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "gosec-parse-column",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         shortFilename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		err = WriteGosecFinding(base.GosecFindingData{
			ProjectName:          module.ProjectName,
			ModuleImportPath:     module.ModuleImportPath,
			ModuleRegistry:       module.ModuleRegistry,
			ModuleVersion:        module.ModuleVersion,
			ModuleNumberGoFiles:  module.ModuleNumberGoFiles,
			ModuleCheckoutFolder: module.PackageDir,
			FileName:             shortFilename,
			FileSizeBytes:        fileToByteCountMap[line.File],
			FileSizeLines:        fileToLineCountMap[line.File],
			FileImportsUnsafePkg: false,
			LineNumber:           lineNumber,
			Column:               column,
			Message:              line.Details,
			Text:                 line.Code,
			Confidence:           line.Confidence,
			Severity:             line.Severity,
			CweId:                line.Cwe.Id,
		})

		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
				Stage:            "gosec-write",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         shortFilename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}