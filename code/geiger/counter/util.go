package counter

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"strings"
)

func getPrintedPackageName(pkg *packages.Package, config Config) string {
	if config.PrintLinkToPkgGoDev && pathLooksLikeAUrl(pkg.PkgPath) {
		return fmt.Sprintf("https://pkg.go.dev/%s", pkg.PkgPath)
	} else {
		return pkg.PkgPath
	}
}

func pathLooksLikeAUrl(path string) bool {
	components := strings.Split(path, "/")

	if len(components) <= 1 {
		return false
	}

	domain := components[0]

	return strings.Contains(domain, ".")
}

func getImportsCount(pkgs map[string]*packages.Package, config Config) (childCount, stdLibCount int) {
	for _, pkg := range pkgs {
		if isStandardPackage(pkg) {
			stdLibCount++
			if config.ShowStandardPackages {
				childCount++
			}
		} else {
			childCount++
		}
	}
	return
}
