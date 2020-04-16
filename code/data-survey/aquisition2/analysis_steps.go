package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	fileToLineCountMap, err := countLines(project, files)
	if err != nil {
		return err
	}


	fileToByteCountMap, err := countBytes(project, files)
	if err != nil {
		return err
	}

	parsedGrepLines, err := grepForUnsafe(project, files)
	if err != nil {
		return err
	}

	analyzeGrepLines(parsedGrepLines, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)

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

func countLines(project *ProjectData, files []string) (map[string]int, error) {
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

func countBytes(project *ProjectData, files []string) (map[string]int, error) {
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

func grepForUnsafe(project *ProjectData, files []string) ([]RipgrepOutputLine, error) {
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