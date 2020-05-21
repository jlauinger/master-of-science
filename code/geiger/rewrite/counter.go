package rewrite

import (
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var packageUnsafeCountCache map[*packages.Package]int

func initCache() {
	packageUnsafeCountCache = map[*packages.Package]int{}
}

func isUnsafePointer(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Pointer" {
			return true
		}
	}
	return false
}

func getUnsafeCount(pkg *packages.Package, config Config) int {
	if config.ShowStandardPackages == false && isStandardPackage(pkg) {
		return 0
	}

	count, ok := packageUnsafeCountCache[pkg]
	if ok {
		return count
	}

	inspectResult := inspector.New(pkg.Syntax)
	unsafePointerCount := 0

	inspectResult.Preorder([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node) {
		node := n.(*ast.SelectorExpr)
		if isUnsafePointer(node) {
			unsafePointerCount++
		}
	})

	packageUnsafeCountCache[pkg] = unsafePointerCount

	return unsafePointerCount
}

func getTotalCount(pkg *packages.Package, config Config, seen *map[*packages.Package]bool) int {
	_, ok := (*seen)[pkg]
	if ok {
		return 0
	}
	(*seen)[pkg] = true

	totalCount := getUnsafeCount(pkg, config)

	for _, child := range pkg.Imports {
		totalCount += getTotalCount(child, config, seen)
	}

	return totalCount
}
