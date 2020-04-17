package analysis

import (
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