package linters

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"os/exec"
	"strconv"
	"strings"
)

/**
 * callback handling the go-safer analysis coordination. This is called for each project
 */
func callbackGosafer(project *base.ProjectData, packages []*base.PackageData, fileToPackageMap map[string]*base.PackageData,
	_, _ map[string]int) {

	// run go-safer, then analyze the findings and write them to disk
	gosaferFindings := runGosafer(project, packages)
	analyzeGosaferFindings(gosaferFindings, fileToPackageMap, project)
}

/**
 * runs go-safer and captures the result
 */
func runGosafer(project *base.ProjectData, packages []*base.PackageData) []base.GosaferFindingLine {
	// analyze packages in chunks of 100 to avoid too much memory consumption when go-safer builds the CFG
	chunkLength := 100
	lines := make([]base.GosaferFindingLine, 0, 1000)
	// go through all packages in chunks of 100
	for i := 0; i < len(packages); i += chunkLength {
		// identify the packages in this chunk
		pkgs := packages[i:base.Min(len(packages), i+chunkLength)]
		// then run go-safer for this chunk
		lines = append(lines, runGosaferEx(project, pkgs, i, len(packages))...)
	}
	return lines
}

/**
 * runs go-safer on an actual chunk of packages
 */
func runGosaferEx(project *base.ProjectData, packages []*base.PackageData, start, length int) []base.GosaferFindingLine {
	fmt.Printf("  running go-safer for %d of %d...\n", start, length)

	// build a list of all import paths of the packages under analysis
	packagePaths := make([]string, 0)
	for _, pkg := range packages {
		// skip analyzing the runtime package however because it creates a stack overflow under the CFG analysis of
		// go-safer
		if pkg.ImportPath == "runtime" {
			continue
		}
		packagePaths = append(packagePaths, pkg.ImportPath)
	}

	// build the go-safer -c=0 command for the packages
	args := []string{"-c=0"}
	args = append(args, packagePaths...)
	cmd := exec.Command("go-safer", args...)
	cmd.Dir = project.CheckoutPath

	// run go-safer and capture the result
	gosaferOutput, _ := cmd.CombinedOutput()

	// split the output into lines and initialize a list of the final findings
	gosaferLines := strings.Split(string(gosaferOutput), "\n")
	gosaferFindings := make([]base.GosaferFindingLine, 0)

	// then go through the result lines
	for i := 0; i < len(gosaferLines); i++ {
		messageLine := gosaferLines[i]

		// if the line contains a comment '#' or starts with downloading or open then ignore it because those lines are
		// not actual results but status or error messages
		if len(messageLine) <= 0 || messageLine[0] == '#' ||
			(len(messageLine) > 12 && messageLine[0:12] == " downloading") ||
			(len(messageLine) > 5 && messageLine[0:5] == "open ") {
			continue
		}
		// if there is an error message indicating a package load error then actually skip the next two lines too
		// because after this error message go-safer will output path information on the missing package
		if len(messageLine) > 18 && messageLine[0:18] == "can't load package" {
			i += 2  // skip GOPATH and GOROOT paths of missing package
			continue
		}

		// build the code context from the following lines
		var contextLines []string
		// use a loop with break condition to capture all context lines that may follow
		for {
			// look at the next gosafer output line
			contextLine := gosaferLines[i+1]

			// see if there are any tab characters, which are included in code output by go-safer
			components := strings.Split(contextLine, "\t")

			// if there aren't, the context has finished so exit the loop
			if len(components) <= 1 {
				break
			}

			// check if the first component, before the tab(s), is a number, because it should be a line number. If not,
			// the context has finished so exit the loop
			_, err := strconv.Atoi(components[0])
			if err != nil {
				break
			}

			// if no exit conditions were met, append this line without the line number and advance lines
			contextLines = append(contextLines, strings.Join(components[1:], "\t"))
			i++
		}

		// finally append the finding line for later analysis
		gosaferFindings = append(gosaferFindings, base.GosaferFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return gosaferFindings
}

/**
 * analyzes identified go-safer finding lines and saves them to disk
 */
func analyzeGosaferFindings(gosaferFindings []base.GosaferFindingLine, fileToPackageMap map[string]*base.PackageData,
	project *base.ProjectData) {

	fmt.Println("  analyzing go-safer output")

	// go through all of the findings
	for _, line := range gosaferFindings {
		// split the message by colons, which separate file name, line number, column, and warning message
		components := strings.Split(line.Message, ":")

		// if there are less than 4 components that means there is an error, so log it and ignore this finding line
		if len(components) < 4 {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "gosafer-ensure-components-length",
				ProjectName:       project.Name,
				PackageImportPath: "",
				FileName:          "",
				Message:           line.Message,
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		// if the first component is the tool name, then ignore that component as it will be a prefix
		if components[0] == "go-safer" || components[0] == "vet" {
			components = components[1:]
		}

		// get the path of the file triggering the warning
		fullFilename := strings.Trim(components[0], " ")

		// try to find the corresponding package. It should be there, if not then use a dummy package
		pkg, ok := fileToPackageMap[fullFilename]
		if !ok {
			pkg = &base.PackageData{
				ImportPath: "unknown-gosafer-error",
			}
		}

		// build a short version of the filename without the package directory for saving to CSV
		filename := fullFilename[len(pkg.Dir)+1:]

		// if the file is a test then ignore it
		if strings.Contains(filename, "_test.go") {
			continue
		}

		// extract line number, column, and message from the line components. If there is a parsing error, log it and
		// ignore this finding line
		lineNumber, err := strconv.Atoi(components[1])
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "gosafer-parse-linenumber",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		column, err := strconv.Atoi(components[2])
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "gosafer-parse-column",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		message := strings.Trim(components[3], " ")

		// finally, save the go-safer finding to disk and log any potential error
		err = base.WriteGosaferFinding(base.GosaferFindingData{
			Message:           message,
			Context:           line.ContextLine,
			LineNumber:        lineNumber,
			Column:            column,
			RawOutput:         line.Message,
			FileName:          filename,
			PackageImportPath: pkg.ImportPath,
			ModulePath:        pkg.ModulePath,
			ModuleVersion:     pkg.ModuleVersion,
			ProjectName:       pkg.ProjectName,
		})
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "gosafer-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}