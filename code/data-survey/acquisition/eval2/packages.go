package eval2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"io"
	"os/exec"
)

func Run(dataDir string, offset, length int, skipProjects []string) {
	packagesFilename := fmt.Sprintf("%s/packages_%d_%d.csv", dataDir, offset, offset + length - 1)
	errorsFilename := fmt.Sprintf("%s/lexical/errors_grep_%d_%d.csv", dataDir, offset, offset + length - 1)

	if err := lexical.OpenPackagesFile(packagesFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	if err := lexical.OpenErrorConditionsFile(errorsFilename); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	defer lexical.CloseFiles()

	projectsFilename := fmt.Sprintf("%s/projects.csv", dataDir)

	fmt.Println("reading projects data...")
	projects, err := lexical.ReadProjects(projectsFilename)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	skipProjectMap := make(map[string]struct{}, len(skipProjects))
	for _, skipProject := range skipProjects {
		skipProjectMap[skipProject] = struct{}{}
	}

	for projectIdx, project := range projects[offset:offset+length] {
		if _, ok := skipProjectMap[project.Name]; ok {
			fmt.Printf("%d/%d (#%d): Skipping %s as requested\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		if !project.UsesModules {
			fmt.Printf("%d/%d (#%d): Skipping %s because it does not use modules\n", projectIdx+1, length, projectIdx+1+offset, project.Name)
			continue
		}

		fmt.Printf("%d/%d (#%d): Analyzing %s\n", projectIdx+1, length, projectIdx+1+offset, project.Name)

		err := analyzeProject(project)
		if err != nil {
			_ = lexical.WriteErrorCondition(lexical.ErrorConditionData{
				Stage:             "project",
				ProjectName:       project.Name,
				PackageImportPath: "",
				FileName:          "",
				Message:           err.Error(),
			})
			fmt.Println("SAVING ERROR!")
			continue
		}
	}
}

func analyzeProject(project *lexical.ProjectData) error {
	packages, err := GetProjectPackages(project)
	if err != nil {
		return err
	}

	fullFilenames := make([]string, 0, 500)
	fileToPackageMap := map[string]*lexical.PackageData{}

	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			fullFilenames = append(fullFilenames, fullFilename)
			fileToPackageMap[fullFilename] = pkg
		}
	}

	analyzeDepTree(packages)

	for _, pkg := range packages {
		fmt.Printf("%s (%s): %d\n", pkg.ImportPath, pkg.ModulePath, pkg.HopCount)
	}

	writePackages(packages)

	return nil
}

func GetProjectPackages(project *lexical.ProjectData) ([]*lexical.PackageData, error) {
	fmt.Println("  identifying relevant packages...")

	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = project.CheckoutPath

	jsonOutput, _ := cmd.Output()

	dec := json.NewDecoder(bytes.NewReader(jsonOutput))
	packages := make([]*lexical.PackageData, 0, 500)

	for {
		var pkg lexical.GoListOutputPackage

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
		} else if pkg.Module.Replace == nil {
			modulePath = pkg.Module.Path
			moduleVersion = pkg.Module.Version
			moduleRegistry = lexical.GetRegistryFromImportPath(pkg.Module.Path)
			moduleIsIndirect = pkg.Module.Indirect
		} else {
			modulePath = pkg.Module.Replace.Path
			moduleVersion = pkg.Module.Replace.Version
			moduleRegistry = lexical.GetRegistryFromImportPath(pkg.Module.Replace.Path)
			moduleIsIndirect = pkg.Module.Replace.Indirect
		}

		packages = append(packages, &lexical.PackageData{
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
		})
	}

	return packages, nil
}

func writePackages(packages []*lexical.PackageData) {
	fmt.Println("  writing package results to disk...")

	for _, pkg := range packages {
		err := lexical.WritePackage(*pkg)
		if err != nil {
			_ = lexical.WriteErrorCondition(lexical.ErrorConditionData{
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

func analyzeDepTree(packages []*lexical.PackageData) {
	packagesGetImported := make(map[string]bool, len(packages))
	packagesMap := make(map[string]*lexical.PackageData, len(packages))

	for _, pkg := range packages {
		packagesGetImported[pkg.ImportPath] = false
		packagesMap[pkg.ImportPath] = pkg
	}

	for _, pkg := range packages {
		for _, childPath := range pkg.Imports {
			if childPath == "C" {
				continue
			}
			child := packagesMap[childPath]
			packagesGetImported[child.ImportPath] = true
		}
	}

	rootPackages := make([]*lexical.PackageData, 0)

	for pkgPath, getsImported := range packagesGetImported {
		if getsImported {
			continue
		}
		pkg := packagesMap[pkgPath]
		if pkg.ImportPath == "runtime/cgo" {
			continue
		}

		rootPackages = append(rootPackages, pkg)
	}

	packagesHopCountMap := make(map[string]int, len(packages))

	for _, pkg := range rootPackages {
		analyzeHopCount(pkg, &packagesMap, &packagesHopCountMap, 0)
	}

	for _, pkg := range packages {
		hopCount, ok := packagesHopCountMap[pkg.ImportPath]
		if !ok {
			pkg.HopCount = 0
		} else {
			pkg.HopCount = hopCount
		}
	}
}

func analyzeHopCount(pkg *lexical.PackageData, packagesMap *map[string]*lexical.PackageData, packagesHopCountMap *map[string]int, hopCount int) {
	previousCount, ok := (*packagesHopCountMap)[pkg.ImportPath]
	if !ok || previousCount > hopCount {
		(*packagesHopCountMap)[pkg.ImportPath] = hopCount
	}

	for _, childPath := range pkg.Imports {
		if childPath == "C" {
			continue
		}
		child, ok := (*packagesMap)[childPath]
		if !ok {
			fmt.Printf("ERROR fetching child path %s\n", childPath)
		}
		analyzeHopCount(child, packagesMap, packagesHopCountMap, hopCount + 1)
	}
}
