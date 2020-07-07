package eval2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"io"
	"os/exec"
)

func Run(dataDir string, offset, length int) {
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

	for projectIdx, project := range projects[offset:offset+length] {
		if !project.UsesModules {
			fmt.Printf("skipping %s because it does not use modules\n", project.Name)
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

		fmt.Printf("  %s\n", pkg.ImportPath)

		/*var modulePath, moduleVersion, moduleRegistry string
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
			moduleRegistry = lexical.GetRegistryFromImportPath(pkg.Module.Path)
			moduleIsIndirect = pkg.Module.Indirect
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
		})*/
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