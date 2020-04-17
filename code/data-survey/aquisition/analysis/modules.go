package analysis

import (
	"bytes"
	"data-aquisition/base"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func analyzeProject(project *base.ProjectData,
	operator func(*base.ProjectData, []base.ModuleData, map[string]base.ModuleData, map[string]int, map[string]int)) error {

	modules, err := getProjectModules(project)
	if err != nil {
		return err
	}

	files := make([]string, 0, 500)
	fileToModuleMap := map[string]base.ModuleData{}

	for _, module := range modules {
		err := WriteModule(module)
		if err != nil {
			_ = WriteErrorCondition(base.ErrorConditionData{
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

	operator(project, modules, fileToModuleMap, fileToLineCountMap, fileToByteCountMap)

	return nil
}

func getProjectModules(project *base.ProjectData) ([]base.ModuleData, error) {
	cmd := exec.Command("go", "list", "-deps", "-json")
	cmd.Dir = project.ProjectCheckoutPath

	jsonOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	modules := make([]base.ModuleData, 0, 500)

	for {
		var pkg base.GoListOutputPackage

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

		modules = append(modules, base.ModuleData{
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