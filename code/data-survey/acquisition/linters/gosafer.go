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
	linterFindings := runGosafer(project, packages)
	analyzeGosaferFindings(linterFindings, fileToPackageMap, project)
}

func runGosafer(project *base.ProjectData, packages []*base.PackageData) []base.GosaferFindingLine {
	chunkLength := 100
	lines := make([]base.GosaferFindingLine, 0, 1000)
	for i := 0; i < len(packages); i += chunkLength {
		pkgs := packages[i:base.Min(len(packages), i+chunkLength)]
		lines = append(lines, runGosaferEx(project, pkgs, i, len(packages))...)
	}
	return lines
}

func runGosaferEx(project *base.ProjectData, packages []*base.PackageData, start, length int) []base.GosaferFindingLine {
	packagePaths := make([]string, 0)

	fmt.Printf("  running go-safer for %d of %d...\n", start, length)

	for _, pkg := range packages {
		if pkg.ImportPath == "runtime" {
			continue
		}
		packagePaths = append(packagePaths, pkg.ImportPath)
	}

	args := []string{"-c=0"}
	args = append(args, packagePaths...)

	cmd := exec.Command("go-safer", args...)
	cmd.Dir = project.CheckoutPath

	gosaferOutput, _ := cmd.CombinedOutput()

	gosaferLines := strings.Split(string(gosaferOutput), "\n")
	gosaferFindings := make([]base.GosaferFindingLine, 0)

	for i := 0; i < len(gosaferLines); i++ {
		messageLine := gosaferLines[i]

		if len(messageLine) <= 0 || messageLine[0] == '#' ||
			(len(messageLine) > 12 && messageLine[0:12] == " downloading") ||
			(len(messageLine) > 5 && messageLine[0:5] == "open ") {
			continue
		}
		if len(messageLine) > 18 && messageLine[0:18] == "can't load package" {
			i += 2  // skip GOPATH and GOROOT paths of missing package
			continue
		}

		var contextLines []string

		for {
			contextLine := gosaferLines[i+1]

			components := strings.Split(contextLine, "\t")

			if len(components) <= 1 {
				break
			}

			_, err := strconv.Atoi(components[0])
			if err != nil {
				break
			}

			contextLines = append(contextLines, strings.Join(components[1:], "\t"))
			i++
		}

		gosaferFindings = append(gosaferFindings, base.GosaferFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return gosaferFindings
}

func analyzeGosaferFindings(gosaferFindings []base.GosaferFindingLine, fileToPackageMap map[string]*base.PackageData,
	project *base.ProjectData) {

	fmt.Println("  analyzing go-safer output")

	for _, line := range gosaferFindings {
		components := strings.Split(line.Message, ":")

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

		var fullFilename string
		var lineNumber int
		var column int
		var message string

		if components[0] == "go-safer" || components[0] == "vet" {
			components = components[1:]
		}

		fullFilename = strings.Trim(components[0], " ")
		pkg, ok := fileToPackageMap[fullFilename]

		if !ok {
			pkg = &base.PackageData{
				ImportPath: "unknown-gosafer-error",
			}
		}

		filename := fullFilename[len(pkg.Dir)+1:]

		if strings.Contains(filename, "_test.go") {
			continue
		}

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
		column, err = strconv.Atoi(components[2])
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
		message = strings.Trim(components[3], " ")

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