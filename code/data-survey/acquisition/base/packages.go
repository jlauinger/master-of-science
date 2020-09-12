package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

/**
 * identifies all packages that belong to a given project using go list
 */
func getProjectPackages(project *ProjectData) ([]*PackageData, error) {
	fmt.Println("  identifying relevant packages...")

	// build the go list -deps -json ./... command to identify packages
	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = project.CheckoutPath

	// run the command and capture its output
	jsonOutput, _ := cmd.Output()

	// initialize a JSON decoder and a list for the packages
	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	packages := make([]*PackageData, 0, 500)

	// repeat until we reach the end of the data
	for {
		var pkg GoListOutputPackage

		// try to decode a go list output message (it's JSON)
		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// extract the module data from the output message
		modulePath, moduleVersion, moduleRegistry, moduleIsIndirect := getModuleData(pkg, project)

		// build the package structure by copying over the relevant data
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
			Imports:          pkg.Imports,
			Deps:             pkg.Deps,
			HopCount:         0,
		})
	}

	return packages, nil
}

/**
 * returns module information for a given package as returned by go list
 */
func getModuleData(pkg GoListOutputPackage, project *ProjectData) (modulePath, moduleVersion, moduleRegistry string, moduleIsIndirect bool) {
	if pkg.Standard {
		// if the package is part of the standard library, use the std placeholder for all module data as the standard
		// library is not modularized
		modulePath = "std"
		moduleVersion = "std"
		moduleRegistry = "std"
		moduleIsIndirect = false
	} else if pkg.Module == nil {
		// similarly, if there is no module information in the go list output, use unknown as a placeholder
		modulePath = "unknown"
		moduleVersion = "unknown"
		moduleRegistry = "unknown"
		moduleIsIndirect = false
	} else if pkg.Module.Replace == nil {
		// if the module information is present and there is no replace information (which would be ranked higher),
		// fill in the information accordingly. Get the registry using a helper function
		modulePath = pkg.Module.Path
		moduleRegistry = getRegistryFromImportPath(pkg.Module.Path)
		moduleIsIndirect = pkg.Module.Indirect
		// finding the module version is the most complicated step
		if pkg.Module.Version != "" {
			// if the version is directly present, use it
			moduleVersion = pkg.Module.Version
		} else if pkg.Module.Path == project.RootModule || strings.HasPrefix(pkg.Module.Path, "./") {
			// if the module is the project module itself or at least a local module (those are not usually versioned
			// because they don't get imported through go.mod), use project as placeholder and get the registry from
			// the project root module
			moduleVersion = "project"
			moduleRegistry = getRegistryFromImportPath(project.RootModule)
		} else {
			// otherwise, use the unknown placeholder
			moduleVersion = "unknown"
		}
	} else {
		// similarly to the regular module, use the same approach for a possible replace module information
		modulePath = pkg.Module.Replace.Path
		moduleRegistry = getRegistryFromImportPath(pkg.Module.Replace.Path)
		moduleIsIndirect = pkg.Module.Replace.Indirect
		if pkg.Module.Replace.Version != "" {
			moduleVersion = pkg.Module.Replace.Version
		} else if pkg.Module.Replace.Path == project.RootModule || strings.HasPrefix(pkg.Module.Replace.Path, "./") {
			moduleVersion = "project"
			moduleRegistry = getRegistryFromImportPath(project.RootModule)
		} else {
			moduleVersion = "unknown"
		}
	}

	return
}

/**
 * fill in the LOC information in the packages list
 */
func fillPackageLOC(packages []*PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	// go through all packages
	for _, pkg := range packages {
		// initialize sum values to identify total LOC for the package
		var loc, byteSize int

		// then to through all Go files in the package
		for _, file := range pkg.GoFiles {
			// build the full filename
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			// and use it to find the LOC and byte size. Then add those to the package totals
			loc += fileToLineCountMap[fullFilename]
			byteSize += fileToByteCountMap[fullFilename]
		}

		// finally store the total LOC information in the package data structure
		pkg.Loc = loc
		pkg.ByteSize = byteSize
	}
}

/**
 * writes the packages CSV file
 */
func writePackages(packages []*PackageData) {
	fmt.Println("  writing package results to disk...")

	// go through all the packages
	for _, pkg := range packages {
		// write the package to disk, and log any errors
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
