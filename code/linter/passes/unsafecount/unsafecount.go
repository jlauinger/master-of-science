package unsafecount

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "unsafecount",
	Doc:  "reports usages of unsafe Pointer",
	Run:  run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
	FactTypes: []analysis.Fact{new(UnsafeCount)},
}

type UnsafeCount struct {
	This int
	Total int
}

func (uc *UnsafeCount) AFact() {}

func (uc *UnsafeCount) String() string {
	return fmt.Sprintf("%d unsafe.Pointer usages", uc.This)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	unsafePointerCount := 0

	inspectResult.Preorder([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node) {
		node := n.(*ast.SelectorExpr)
		if isUnsafePointer(node) {
			unsafePointerCount++
		}
	})

	fact := &UnsafeCount{
		This:  unsafePointerCount,
		Total: unsafePointerCount,
	}

	for _, pkg := range pass.Pkg.Imports() {
		var pkgUnsafeCount UnsafeCount
		pass.ImportPackageFact(pkg, &pkgUnsafeCount)
		fact.Total += pkgUnsafeCount.Total
	}

	pass.ExportPackageFact(fact)

	return nil, nil
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