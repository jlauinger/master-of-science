package linters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"os/exec"
	"strconv"
	"strings"
)

/**
 * callback handling the gosec analysis coordination. This is called for each project
 */
func callbackGosec(project *base.ProjectData, packages []*base.PackageData, fileToPackageMap map[string]*base.PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	// run gosec, then analyze the findings and write them to disk
	gosecFindings, _ := runGosec(project, packages)
	analyzeGosecFindings(gosecFindings, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}

func runGosec(project *base.ProjectData, packages []*base.PackageData) ([]base.GosecIssueOutput, error) {
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
	var gosecResult base.GosecOutput

	err := dec.Decode(&gosecResult)
	if err != nil {
		_ = base.WriteErrorCondition(base.ErrorConditionData{
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

func analyzeGosecFindings(gosecFindings []base.GosecIssueOutput, fileToPackageMap map[string]*base.PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	fmt.Println("  analyzing gosec output")

	for _, line := range gosecFindings {
		pkg, ok := fileToPackageMap[line.File]
		if !ok {
			pkg = &base.PackageData{
				ImportPath: "unknown-vet-error",
			}
		}

		shortFilename := line.File[len(pkg.Dir)+1:]

		var lineNumberText string
		if strings.Contains(line.Line, "-") {
			lineNumberText = strings.Split(line.Line, "-")[0]
		} else {
			lineNumberText = line.Line
		}
		lineNumber, err := strconv.Atoi(lineNumberText)
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
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
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "gosec-parse-column",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          shortFilename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		err = base.WriteGosecFinding(base.GosecFindingData{
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
		})

		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
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
}