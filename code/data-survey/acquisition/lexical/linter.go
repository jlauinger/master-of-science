package lexical

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runLinter(project *ProjectData, packages []*PackageData) []LinterFindingLine {
	chunkLength := 100
	lines := make([]LinterFindingLine, 0, 1000)
	for i := 0; i < len(packages); i += chunkLength {
		pkgs := packages[i:Min(len(packages), i+chunkLength)]
		lines = append(lines, runLinterEx(project, pkgs, i, len(packages))...)
	}
	return lines
}

func runLinterEx(project *ProjectData, packages []*PackageData, start, length int) []LinterFindingLine {
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

	linterOutput, _ := cmd.CombinedOutput()

	linterLines := strings.Split(string(linterOutput), "\n")
	linterFindings := make([]LinterFindingLine, 0)

	for i := 0; i < len(linterLines); i++ {
		messageLine := linterLines[i]

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
			contextLine := linterLines[i+1]

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

		linterFindings = append(linterFindings, LinterFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return linterFindings
}

func analyzeLinterFindings(linterFindings []LinterFindingLine, fileToPackageMap map[string]*PackageData,
	project *ProjectData) map[string]string {

	fmt.Println("  analyzing linter output")

	for _, line := range linterFindings {
		components := strings.Split(line.Message, ":")

		if len(components) < 4 {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "linter-ensure-components-length",
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

		if components[0] == "linter" || components[0] == "vet" {
			components = components[1:]
		}

		fullFilename = strings.Trim(components[0], " ")
		pkg, ok := fileToPackageMap[fullFilename]

		if !ok {
			pkg = &PackageData{
				ImportPath: "unknown-linter-error",
			}
		}

		filename := fullFilename[len(pkg.Dir)+1:]

		if strings.Contains(filename, "_test.go") {
			continue
		}

		lineNumber, err := strconv.Atoi(components[1])
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "linter-parse-linenumber",
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
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "linter-parse-column",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		message = strings.Trim(components[3], " ")

		err = WriteLinterFinding(LinterFindingData{
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
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "linter-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}

	return map[string]string{}
}