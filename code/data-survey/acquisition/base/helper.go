package base

import (
	"strings"
)

/**
 * return the bigger of two integers
 */
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

/**
 * returns the smaller of two integers
 */
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

/**
 * returns the registry part from a package import path string
 */
func getRegistryFromImportPath(importPath string) string {
	// split the path by slashes
	pathComponents := strings.Split(importPath, "/")

	// if there is no or only one slash, use the first component as registry
	if len(pathComponents) <= 1 {
		return pathComponents[0]
	}

	// otherwise, use the first component and/or add the second if it is an x, to accomodate for constructs such as
	// golang.org/x, which I use as the full registry name
	var registryComponents []string
	if pathComponents[1] == "x" {
		registryComponents = pathComponents[0:2]
	} else {
		registryComponents = pathComponents[0:1]
	}

	// since the components are a list now, join them together with the original slashes
	return strings.Join(registryComponents, "/")
}
