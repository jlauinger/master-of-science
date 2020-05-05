package analysis

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
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
	fmt.Printf("") // to "need" fmt Package

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.AssignStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CompositeLit:
			if (compositeLiteralIsReflectHeader(n, pass)) {
				pass.Reportf(n.Pos(), "reflect header composite literal found: %q", render(pass.Fset, n))
			}
		case *ast.AssignStmt:
			if (assigningToReflectHeader(n, pass)) {
				pass.Reportf(n.Pos(), "assigning to reflect header object: %q", render(pass.Fset, n))
			}
		default:
			return
		}
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

func compositeLiteralIsReflectHeader(cl *ast.CompositeLit, pass *analysis.Pass) bool {
	literalType, ok := pass.TypesInfo.Types[cl]
	if !ok {
		return false
	}

	sliceHeaderType := types.NewStruct([]*types.Var{
		types.NewVar(token.NoPos, nil, "Data", types.Typ[types.Uintptr]),
		types.NewVar(token.NoPos, nil, "Len", types.Typ[types.Int]),
		types.NewVar(token.NoPos, nil, "Cap", types.Typ[types.Int]),
	}, nil)
	stringHeaderType := types.NewStruct([]*types.Var{
		types.NewVar(token.NoPos, nil, "Data", types.Typ[types.Uintptr]),
		types.NewVar(token.NoPos, nil, "Len", types.Typ[types.Int]),
	}, nil)

	if (types.AssignableTo(sliceHeaderType, literalType.Type.Underlying()) ||
		types.AssignableTo(stringHeaderType, literalType.Type.Underlying())) {
		return true
	}

	return false
}

func assigningToReflectHeader(assignStmt *ast.AssignStmt, pass *analysis.Pass) bool {
	for _, expr := range assignStmt.Lhs {
		lhs, ok := expr.(*ast.SelectorExpr)
		if !ok {
			return false
		}

		lhsType, ok := pass.TypesInfo.Types[lhs.X]
		if !ok {
			return false
		}

		t := lhsType.Type.String()

		if t == "*reflect.SliceHeader" || t == "reflect.SliceHeader" || t == "*reflect.StringHeader" || t == "reflect.StringHeader" {
			return true
		}
	}
	return false
}