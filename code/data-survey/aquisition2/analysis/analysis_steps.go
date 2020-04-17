package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func analyzeProject(project *ProjectData) error {
	modules, err := getProjectModules(project)
	if err != nil {
		return err
	}

	files := make([]string, 0, 500)
	fileToModuleMap := map[string]ModuleData{}

	for _, module := range modules {
		err := WriteModule(module)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "module",
				ProjectName:      project.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         "",
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		for _, file := range module.PackageGoFiles {
			fullFilename := fmt.Sprintf("%s/%s", module.PackageDir, file)
			files = append(files, fullFilename)
			fileToModuleMap[fullFilename] = module
		}
	}

	fileToLineCountMap, err := countLines(files)
	if err != nil {
		return err
	}


	fileToByteCountMap, err := countBytes(files)
	if err != nil {
		return err
	}

	parsedGrepLines, err := grepForUnsafe(files)
	if err != nil {
		return err
	}
	analyzeGrepLines(parsedGrepLines, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)

	vetFindings := runVet(project, modules)
	analyzeVetFindings(vetFindings, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)

	gosecFindings, _ := runGosec(project, modules)
	analyzeGosecFindings(gosecFindings, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)

	return nil
}

func getProjectModules(project *ProjectData) ([]ModuleData, error) {
	cmd := exec.Command("go", "list", "-deps", "-json")
	cmd.Dir = project.ProjectCheckoutPath

	jsonOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	modules := make([]ModuleData, 0, 500)

	for {
		var pkg GoListOutputPackage

		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var moduleRegistry string
		if pkg.Standard {
			moduleRegistry = "std"
		} else {
			moduleRegistry = getRegistryFromImportPath(pkg.ImportPath)
		}

		modules = append(modules, ModuleData{
			ProjectName:          project.ProjectName,
			ModuleImportPath:     pkg.ImportPath,
			ModuleRegistry:       moduleRegistry,
			ModuleVersion:        "",
			ModuleNumberGoFiles:  len(pkg.GoFiles),
			PackageDir:			  pkg.Dir,
			PackageGoFiles: 	  pkg.GoFiles,
		})
	}

	return modules, nil
}

func countLines(files []string) (map[string]int, error) {
	args := []string{"-l"}
	args = append(args, files...)

	cmd := exec.Command("wc", args...)

	wcOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	filesToLineCount := map[string]int{}

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesToLineCount[components[1]] = count
	}

	return filesToLineCount, nil
}

func countBytes(files []string) (map[string]int, error) {
	args := []string{"-c"}
	args = append(args, files...)

	cmd := exec.Command("wc", args...)

	wcOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	filesToByteCount := map[string]int{}

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesToByteCount[components[1]] = count
	}

	return filesToByteCount, nil
}

func grepForUnsafe(files []string) ([]RipgrepOutputLine, error) {
	args := []string{"unsafe.Pointer", "--context", "5", "--json"}
	args = append(args, files...)

	cmd := exec.Command("rg", args...)

	rgOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(rgOutput))
	parsedLines := make([]RipgrepOutputLine, 0, 1000)

	for {
		var message RipgrepOutputLine

		err := dec.Decode(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		parsedLines = append(parsedLines, message)
	}

	return parsedLines, nil
}

func analyzeGrepLines(parsedLines []RipgrepOutputLine, fileToModuleMap map[string]ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {
	for lineIdx, line := range parsedLines {
		if line.MessageType == "match" {
			contextLines := []string{line.Data.Lines.Text}

			// context before line
			for contextIdx := lineIdx - 1; contextIdx > Max(0, lineIdx - 5); contextIdx-- {
				contextLine := parsedLines[contextIdx]
				if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
					contextLines = append([]string{contextLine.Data.Lines.Text}, contextLines...)
				} else {
					break
				}
			}

			// context after line
			for contextIdx := lineIdx + 1; contextIdx < Min(len(parsedLines), lineIdx + 6); contextIdx++ {
				contextLine := parsedLines[contextIdx]
				if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
					contextLines = append(contextLines, contextLine.Data.Lines.Text)
				} else {
					break
				}
			}

			context := strings.Join(contextLines, "")

			fullFilename := line.Data.Path.Text
			module := fileToModuleMap[fullFilename]
			filename := fullFilename[len(module.PackageDir)+1:]

			err := WriteMatchResult(MatchResultData{
				ProjectName:          module.ProjectName,
				ModuleImportPath:     module.ModuleImportPath,
				ModuleRegistry:       module.ModuleRegistry,
				ModuleVersion:        module.ModuleVersion,
				ModuleNumberGoFiles:  module.ModuleNumberGoFiles,
				ModuleCheckoutFolder: module.PackageDir,
				FileName:             filename,
				FileSizeBytes:        fileToByteCountMap[fullFilename],
				FileSizeLines:        fileToLineCountMap[fullFilename],
				FileImportsUnsafePkg: false, // TODO
				FileGoVetOutput:      "", // TODO
				Text:                 line.Data.Lines.Text,
				Context:              context,
				LineNumber:           line.Data.LineNumber,
				ByteOffset:           line.Data.LineNumber,
				MatchType:            line.Data.SubMatches[0].Match.Text,
			})

			if err != nil {
				_ = WriteErrorCondition(ErrorConditionData{
					Stage:            "ripgrep-parse",
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
}

func runVet(project *ProjectData, modules []ModuleData) []VetFindingLine {
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
	vetFindings := make([]VetFindingLine, 0)

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

		vetFindings = append(vetFindings, VetFindingLine{
			Message:     messageLine,
			ContextLine: strings.Join(contextLines, "\n"),
		})
	}

	return vetFindings
}

func analyzeVetFindings(vetFindings []VetFindingLine, fileToModuleMap map[string]ModuleData,
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
			_ = WriteErrorCondition(ErrorConditionData{
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
			_ = WriteErrorCondition(ErrorConditionData{
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

		err = WriteVetFinding(VetFindingData{
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
			_ = WriteErrorCondition(ErrorConditionData{
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

func runGosec(project *ProjectData, modules []ModuleData) ([]GosecIssueOutput, error) {
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
	var gosecResult GosecOutput

	err := dec.Decode(&gosecResult)
	if err != nil {
		_ = WriteErrorCondition(ErrorConditionData{
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

func analyzeGosecFindings(gosecFindings []GosecIssueOutput, fileToModuleMap map[string]ModuleData,
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
			_ = WriteErrorCondition(ErrorConditionData{
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
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:            "gosec-parse-column",
				ProjectName:      module.ProjectName,
				ModuleImportPath: module.ModuleImportPath,
				FileName:         shortFilename,
				Message:          err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		err = WriteGosecFinding(GosecFindingData{
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
			_ = WriteErrorCondition(ErrorConditionData{
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