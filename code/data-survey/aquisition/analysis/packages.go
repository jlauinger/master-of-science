package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func analyzeProject(project *ProjectData,
	operator func(*ProjectData, []PackageData, map[string]PackageData, map[string]int, map[string]int)) error {

	packages, err := getProjectPackages(project)
	if err != nil {
		return err
	}

	files := make([]string, 0, 500)
	fileToPackageMap := map[string]PackageData{}

	for _, pkg := range packages {
		err := WritePackage(pkg)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "package",
				ProjectName:       project.Name,
				PackageImportPath: pkg.ImportPath,
				FileName:          "",
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}

		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			files = append(files, fullFilename)
			fileToPackageMap[fullFilename] = pkg
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

	operator(project, packages, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)

	return nil
}

func getProjectPackages(project *ProjectData) ([]PackageData, error) {
	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = project.CheckoutPath

	jsonOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	packages := make([]PackageData, 0, 500)

	for {
		var pkg GoListOutputPackage

		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var modulePath, moduleVersion, moduleRegistry string
		var moduleIsIndirect bool

		if pkg.Module == nil {
			modulePath = "std"
			moduleVersion = "std"
			moduleRegistry = "std"
			moduleIsIndirect = false
		} else {
			modulePath = pkg.Module.Path
			moduleVersion = pkg.Module.Version
			moduleRegistry = getRegistryFromImportPath(pkg.Module.Path)
			moduleIsIndirect = pkg.Module.Indirect
		}

		packages = append(packages, PackageData{
			Name:             pkg.Name,
			ImportPath:       pkg.ImportPath,
			Dir:              pkg.Dir,
			IsStandard:       pkg.Standard,
			IsDepOnly:        pkg.DepOnly,
			NumberOfGoFiles:  len(pkg.GoFiles),
			Loc:              0,
			ByteSize:         0,
			ModulePath:       modulePath,
			ModuleVersion:    moduleVersion,
			ModuleRegistry:   moduleRegistry,
			ModuleIsIndirect: moduleIsIndirect,
			ProjectName:      project.Name,
			GoFiles:          pkg.GoFiles,
		})
	}

	return packages, nil
}