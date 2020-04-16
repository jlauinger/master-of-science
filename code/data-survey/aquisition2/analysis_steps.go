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
		}

		for _, file := range module.PackageGoFiles {
			fullFilename := fmt.Sprintf("%s/%s", module.PackageDir, file)
			files = append(files, fullFilename)
			fileToModuleMap[fullFilename] = module
		}
	}

	filesWithLineCount, err := countLines(project, files)
	if err != nil {
		return err
	}
	for _, fileWithCount := range filesWithLineCount[:10] {
		fmt.Printf("%d %s (%s)\n", fileWithCount.Count, fileWithCount.Filename,
			fileToModuleMap[fileWithCount.Filename].ModuleImportPath)
	}
	fmt.Println("---------------------------")


	filesWithByteCount, err := countBytes(project, files)
	if err != nil {
		return err
	}
	for _, fileWithCount := range filesWithByteCount[:10] {
		fmt.Printf("%d %s (%s)\n", fileWithCount.Count, fileWithCount.Filename,
			fileToModuleMap[fileWithCount.Filename].ModuleImportPath)
	}
	fmt.Println("---------------------------")

	parsedGrepLines, err := grepForUnsafe(project, files)
	if err != nil {
		return err
	}
	for _, parsedLine := range parsedGrepLines[:20] {
		fmt.Printf("%s %d %s (%s)\n", parsedLine.MessageType, parsedLine.Data.LineNumber,
			parsedLine.Data.Path.Text, fileToModuleMap[parsedLine.Data.Path.Text].ModuleImportPath)
		if parsedLine.MessageType == "match" {
			fmt.Println(parsedLine.Data.Lines.Text)
		}
	}

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

func countLines(project *ProjectData, files []string) ([]FilenameWithCount, error) {
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

	filesWithLineCount := make([]FilenameWithCount, 0, len(outputLines))

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesWithLineCount = append(filesWithLineCount, FilenameWithCount{
			Filename: components[1],
			Count:    count,
		})
	}

	return filesWithLineCount, nil
}

func countBytes(project *ProjectData, files []string) ([]FilenameWithCount, error) {
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

	filesWithByteCount := make([]FilenameWithCount, 0, len(outputLines))

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesWithByteCount = append(filesWithByteCount, FilenameWithCount{
			Filename: components[1],
			Count:    count,
		})
	}

	return filesWithByteCount, nil
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
	parsedLines := []RipgrepOutputLine{}

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