package analysis

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var LiteralHeaderAnalyzer = &analysis.Analyzer{
	Name: "literalHeader",
	Doc:  "reports reflect.SliceHeader and reflect.StringHeader composite literals",
	Run:  run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		cl := n.(*ast.CompositeLit)

		if (!compositeLiteralIsReflectHeader(cl)) {
			return
		}

		pass.Reportf(cl.Pos(), "reflect header composite literal found: %q", render(pass.Fset, cl))
	})

	return nil, nil
}

func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func compositeLiteralIsReflectHeader(cl *ast.CompositeLit) bool {
	switch typ := cl.Type.(type) {
	case *ast.SelectorExpr:
		switch X := typ.X.(type) {
		case *ast.Ident:
			return X.Name == "reflect" && (typ.Sel.Name == "SliceHeader" || typ.Sel.Name == "StringHeader")
		default:
			return false
		}
	default:
		return false
	}
}

