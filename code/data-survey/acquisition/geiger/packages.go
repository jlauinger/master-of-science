package geiger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"io"
	"os/exec"
	"strings"
)

func GetProjectPackages(project *base.ProjectData) ([]*base.PackageData, error) {
	fmt.Println("  identifying relevant packages...")

	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = project.CheckoutPath

	jsonOutput, _ := cmd.Output()

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	packages := make([]*base.PackageData, 0, 500)

	for {
		var pkg base.GoListOutputPackage

		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		modulePath, moduleVersion, moduleRegistry, moduleIsIndirect := getModuleData(pkg, project)

		packages = append(packages, &base.PackageData{
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

func getModuleData(pkg base.GoListOutputPackage, project *base.ProjectData) (modulePath, moduleVersion, moduleRegistry string, moduleIsIndirect bool) {
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
	} else if pkg.Module.Replace == nil {
		modulePath = pkg.Module.Path
		moduleRegistry = base.GetRegistryFromImportPath(pkg.Module.Path)
		moduleIsIndirect = pkg.Module.Indirect
		if pkg.Module.Version != "" {
			moduleVersion = pkg.Module.Version
		} else if pkg.Module.Path == project.RootModule || strings.HasPrefix(pkg.Module.Path, "./") {
			moduleVersion = "project"
			moduleRegistry = base.GetRegistryFromImportPath(project.RootModule)
		} else {
			moduleVersion = "unknown"
		}
	} else {
		modulePath = pkg.Module.Replace.Path
		moduleRegistry = base.GetRegistryFromImportPath(pkg.Module.Replace.Path)
		moduleIsIndirect = pkg.Module.Replace.Indirect
		if pkg.Module.Replace.Version != "" {
			moduleVersion = pkg.Module.Replace.Version
		} else if pkg.Module.Replace.Path == project.RootModule || strings.HasPrefix(pkg.Module.Replace.Path, "./") {
			moduleVersion = "project"
			moduleRegistry = base.GetRegistryFromImportPath(project.RootModule)
		} else {
			moduleVersion = "unknown"
		}
	}

	return
}

func fillPackageLOC(packages []*base.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
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

func writePackages(packages []*base.PackageData) {
	fmt.Println("  writing package results to disk...")

	for _, pkg := range packages {
		err := base.WritePackage(*pkg)
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
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
