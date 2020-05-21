package rewrite

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

type Config struct {
	MaxIndent            int
	ShortenSeenPackages  bool
	ShowStandardPackages bool
}


func Run(config Config, paths... string) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedFiles | packages.NeedName,
		Tests:      false,
	}, paths...)

	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		printPkgTree(pkg, 0, config)
	}
}

func printPkgTree(pkg *packages.Package, indent int, config Config) {
	if config.ShowStandardPackages == false && isStandardPackage(pkg.ID) {
		return
	}

	if indent > config.MaxIndent {
		return
	}
	printIndent(indent)

	countInThisPackage := getUnsafeCount(pkg)
	totalCount := getTotalCount(pkg, &map[*packages.Package]bool{})
	fmt.Printf("%s %d/%d\n", pkg.PkgPath, countInThisPackage, totalCount)

	for _, child := range pkg.Imports {
		printPkgTree(child, indent + 1, config)
	}
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
}

func getUnsafeCount(pkg *packages.Package) int {
	inspectResult := inspector.New(pkg.Syntax)

	unsafePointerCount := 0

	inspectResult.Preorder([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node) {
		node := n.(*ast.SelectorExpr)
		if isUnsafePointer(node) {
			unsafePointerCount++
		}
	})

	return unsafePointerCount
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

func getTotalCount(pkg *packages.Package, seen *map[*packages.Package]bool) int {
	_, ok := (*seen)[pkg]
	if ok {
		return 0
	}
	(*seen)[pkg] = true

	totalCount := getUnsafeCount(pkg)

	for _, child := range pkg.Imports {
		totalCount += getTotalCount(child, seen)
	}

	return totalCount
}