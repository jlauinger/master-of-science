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

/**
 * runs gosec on the project and captures the results
 */
func runGosec(project *base.ProjectData, packages []*base.PackageData) ([]base.GosecIssueOutput, error) {
	fmt.Println("  running gosec...")

	// build a list with the import paths of all packages
	packagePaths := make([]string, 0, 1000)
	for _, pkg := range packages {
		packagePaths = append(packagePaths, pkg.Dir)
	}

	// build the gosec command to analyze the packages
	args := []string{"-quiet", "-no-fail", "-fmt=json"}
	args = append(args, packagePaths...)
	cmd := exec.Command("gosec", args...)
	cmd.Dir = project.CheckoutPath

	// run gosec and capture the output
	gosecOutput, _ := cmd.CombinedOutput()

	// initialize a new JSON decoder for the output
	dec := json.NewDecoder(bytes.NewReader(gosecOutput))
	var gosecResult base.GosecOutput

	// try decoding the output. This is not done in a loop because the gosec output is a single JSON object with the
	// findings contained in a JSON array field. If there is an error, log it
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

	// the identified warnings are contained in the Issues array, so we only need this to go forward.
	return gosecResult.Issues, nil
}

/**
 * analyzes the gosec findings and writes them to disk
 */
func analyzeGosecFindings(gosecFindings []base.GosecIssueOutput, fileToPackageMap map[string]*base.PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	fmt.Println("  analyzing gosec output")

	// go through all of the findings
	for _, line := range gosecFindings {
		// try to identify the package containing the file of this finding. It should be there, but if not use a dummy
		pkg, ok := fileToPackageMap[line.File]
		if !ok {
			pkg = &base.PackageData{
				ImportPath: "unknown-gosec-error",
			}
		}

		// build a short version of the filename without the package directory for storage in CSV
		shortFilename := line.File[len(pkg.Dir)+1:]

		// the line number might be a range of lines (e.g. 42-45). Check if the line number field contains a - character
		// and if so use only the first component, i.e. the start line
		var lineNumberText string
		if strings.Contains(line.Line, "-") {
			lineNumberText = strings.Split(line.Line, "-")[0]
		} else {
			lineNumberText = line.Line
		}

		// try parsing the line number and log any potential errors
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

		// similarly to the line number, use only the starting column if the column field contains a range (e.g. 10-30)
		var columnText string
		if strings.Contains(line.Column, "-") {
			columnText = strings.Split(line.Column, "-")[0]
		} else {
			columnText = line.Column
		}

		// try parsing the column number and log any potential errors
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

		// finally, build the gosec finding object and write it to CSV. If there should be an error, log it
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