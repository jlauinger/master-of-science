package lexical

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runVet(project *ProjectData, packages []*PackageData) []VetFindingLine {
	packagePaths := make([]string, len(packages))

	fmt.Println("  running go vet")

	for i, pkg := range packages {
		packagePaths[i] = pkg.ImportPath
	}

	args := []string{"vet", "-c=0"}
	args = append(args, packagePaths...)

	cmd := exec.Command("go", args...)
	cmd.Dir = project.CheckoutPath

	vetOutput, _ := cmd.CombinedOutput()

	vetLines := strings.Split(string(vetOutput), "\n")
	vetFindings := make([]VetFindingLine, 0)

	for i := 0; i < len(vetLines); i++ {
		messageLine := vetLines[i]

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
			contextLine := vetLines[i+1]

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

		vetFindings = append(vetFindings, VetFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return vetFindings
}

func analyzeVetFindings(vetFindings []VetFindingLine, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int, project *ProjectData) map[string]string {

	fmt.Println("  analyzing go vet output")

	var filesToCopy = make(map[string]string, 500)

	for _, line := range vetFindings {
		components := strings.Split(line.Message, ":")

		if len(components) < 4 {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "vet-ensure-components-length",
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

		if components[0] == "vet" {
			components = components[1:]
		}

		fullFilename = strings.Trim(components[0], " ")
		pkg, ok := fileToPackageMap[fullFilename]

		if !ok {
			pkg = &PackageData{
				ImportPath: "unknown-vet-error",
			}
		}

		filename := fullFilename[len(pkg.Dir)+1:]

		if strings.Contains(filename, "test") {
			continue
		}

		lineNumber, err := strconv.Atoi(components[1])
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "vet-parse-linenumber",
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
				Stage:             "vet-parse-column",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
		message = strings.Trim(components[3], " ")

		copyDestination := fmt.Sprintf("%s/%s", pkg.ImportPath, filename)

		err = WriteVetFinding(VetFindingData{
			Message:           message,
			Context:           line.ContextLine,
			LineNumber:        lineNumber,
			Column:            column,
			RawOutput:         line.Message,
			FileName:          filename,
			FileLoc:           fileToLineCountMap[fullFilename],
			FileByteSize:      fileToByteCountMap[fullFilename],
			PackageImportPath: pkg.ImportPath,
			ModulePath:        pkg.ModulePath,
			ModuleVersion:     pkg.ModuleVersion,
			ProjectName:       pkg.ProjectName,
			FileCopyPath:      copyDestination,
		})

		filesToCopy[fullFilename] = copyDestination

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "vet-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          filename,
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}

	return filesToCopy
}