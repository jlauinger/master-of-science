package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runGosec(project *ProjectData, packages []*PackageData) ([]GosecIssueOutput, error) {
	fmt.Println("  running gosec...")

	packagePaths := make([]string, 0, 1000)

	for _, pkg := range packages {
		packagePaths = append(packagePaths, pkg.Dir)
	}

	args := []string{"-quiet", "-no-fail", "-fmt=json"}
	args = append(args, packagePaths...)

	cmd := exec.Command("gosec", args...)
	cmd.Dir = project.CheckoutPath

	gosecOutput, _ := cmd.CombinedOutput()

	dec := json.NewDecoder(bytes.NewReader(gosecOutput))
	var gosecResult GosecOutput

	err := dec.Decode(&gosecResult)
	if err != nil {
		_ = WriteErrorCondition(ErrorConditionData{
			Stage:             "gosec-parse",
			ProjectName:       project.Name,
			PackageImportPath: "",
			FileName:          "",
			Message:           err.Error(),
		})
		fmt.Println("SAVING ERROR!")
		return nil, err
	}

	return gosecResult.Issues, nil
}

func analyzeGosecFindings(gosecFindings []GosecIssueOutput, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	fmt.Println("  analyzing gosec output")

	var filesToCopy = make(map[string]string, 500)

	for _, line := range gosecFindings {
		pkg := fileToPackageMap[line.File]
		shortFilename := line.File[len(pkg.Dir)+1:]

		var lineNumberText string
		if strings.Contains(line.Line, "-") {
			lineNumberText = strings.Split(line.Line, "-")[0]
		} else {
			lineNumberText = line.Line
		}
		lineNumber, err := strconv.Atoi(lineNumberText)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "gosec-parse-linenumber",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          shortFilename,
				Message:           err.Error(),
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
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "gosec-parse-column",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          shortFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		copyDestination := fmt.Sprintf("%s/%s", pkg.ImportPath, shortFilename)

		err = WriteGosecFinding(GosecFindingData{
			Message:           line.Details,
			Context:           line.Code,
			Confidence:        line.Confidence,
			Severity:          line.Severity,
			CweId:             line.Cwe.Id,
			RuleId:            line.RuleId,
			LineNumber:        lineNumber,
			Column:            column,
			FileName:          shortFilename,
			FileLoc:           fileToLineCountMap[line.File],
			FileByteSize:      fileToByteCountMap[line.File],
			PackageImportPath: pkg.ImportPath,
			ModulePath:        pkg.ModulePath,
			ModuleVersion:     pkg.ModuleVersion,
			ProjectName:       pkg.ProjectName,
			FileCopyPath:      copyDestination,
		})

		filesToCopy[line.File] = copyDestination

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "gosec-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          shortFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}

	return filesToCopy
}