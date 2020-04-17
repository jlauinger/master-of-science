package analysis

import (
	"fmt"
	"os"
	"strings"
)

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func goModExists(project *ProjectData) bool {
	_, err := os.Stat(fmt.Sprintf("%s/go.mod", project.ProjectCheckoutPath))
	return err == nil
}

func getRegistryFromImportPath(importPath string) string {
	pathComponents := strings.Split(importPath, "/")

	if len(pathComponents) <= 1 {
		return "std"
	}

	var registryComponents []string
	if pathComponents[1] == "x" {
		registryComponents = pathComponents[0:2]
	} else {
		registryComponents = pathComponents[0:1]
	}

	return strings.Join(registryComponents, "/")
}