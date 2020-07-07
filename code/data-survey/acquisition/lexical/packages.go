package lexical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func analyzeProject(project *ProjectData, writePackagesToFile bool,
	operator func(*ProjectData, []*PackageData, map[string]*PackageData, map[string]int, map[string]int) map[string]string) (map[string]string, error) {

	packages, err := GetProjectPackages(project)
	if err != nil {
		return map[string]string{}, err
	}

	fullFilenames := make([]string, 0, 500)
	fileToPackageMap := map[string]*PackageData{}

	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			fullFilenames = append(fullFilenames, fullFilename)
			fileToPackageMap[fullFilename] = pkg
		}
	}

	fileToLineCountMap, err := countLines(fullFilenames)
	if err != nil {
		return map[string]string{}, err
	}


	fileToByteCountMap, err := countBytes(fullFilenames)
	if err != nil {
		return map[string]string{}, err
	}

	fillPackageLOC(packages, fileToLineCountMap, fileToByteCountMap)

	if writePackagesToFile {
		writePackages(packages)
	}

	filesToCopy := operator(project, packages, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)

	return filesToCopy, nil
}

func GetProjectPackages(project *ProjectData) ([]*PackageData, error) {
	fmt.Println("  identifying relevant packages...")

	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = project.CheckoutPath

	jsonOutput, _ := cmd.Output()

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	packages := make([]*PackageData, 0, 500)

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

		if pkg.Standard {
			modulePath = "std"
			moduleVersion = "std"
			moduleRegistry = "std"
			moduleIsIndirect = false
		} else if pkg.Module == nil {
			modulePath = "unknown"
			moduleVersion = "unknown"
			moduleRegistry = "unknown"
			moduleIsIndirect = false
		} else {
			modulePath = pkg.Module.Path
			moduleVersion = pkg.Module.Version
			moduleRegistry = GetRegistryFromImportPath(pkg.Module.Path)
			moduleIsIndirect = pkg.Module.Indirect
		}

		packages = append(packages, &PackageData{
			Name:             pkg.Name,
			ImportPath:       pkg.ImportPath,
			Dir:              pkg.Dir,
			IsStandard:       pkg.Standard,
			IsDepOnly:        pkg.DepOnly,
			NumberOfGoFiles:  len(pkg.GoFiles),
			Loc:              0, // filled later
			ByteSize:         0, // filled later
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

func writePackages(packages []*PackageData) {
	fmt.Println("  writing package results to disk...")

	for _, pkg := range packages {
		err := WritePackage(*pkg)
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "package",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          "",
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}

func fillPackageLOC(packages []*PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	for _, pkg := range packages {
		var loc, byteSize int

		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			loc += fileToLineCountMap[fullFilename]
			byteSize += fileToByteCountMap[fullFilename]
		}

		pkg.Loc = loc
		pkg.ByteSize = byteSize
	}
}