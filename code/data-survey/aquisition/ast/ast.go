package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func AnalyzeAst(mode string) {
	filename := "/home/johannes/studium/s14/masterarbeit/download/bosun/vendor/github.com/bradfitz/slice/slice.go"
	code, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(code), "\n")

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "slice.go", code, parser.ParseComments)

	findingsTree := &AstTreeNode{
		Node:     nil,
		Children: []*AstTreeNode{},
	}

	ast.Walk(UnsafeVisitor{
		fileset:      fset,
		context:      []ast.Node{},
		findingsTree: findingsTree,
	}, node)

	findingsTree.countUnsafePointer()
	findingsTree.countUintptr()

	switch mode {
	case "tree":
		findingsTree.printRoot(fset)
	case "func":
		formatFunctions(findingsTree, fset, lines)
	case "stmt":
		formatStatements(findingsTree, fset, lines)
	default:
		fmt.Printf("unknown mode %s not in tree,func,stmt\n", mode)
	}
}

type UnsafeVisitor struct {
	fileset *token.FileSet
	context []ast.Node
	findingsTree *AstTreeNode
}

func (uv UnsafeVisitor) Visit(n ast.Node) ast.Visitor {
	if isUnsafePointer(n) || isUintptr(n) {
		uv.findingsTree.addPath(append(uv.context, n))
	}

	return UnsafeVisitor{
		fileset: uv.fileset,
		context: append(uv.context, n),
		findingsTree: uv.findingsTree,
	}
}

func isUnsafePointer(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" {
				return true
			}
		}
	}
	return false
}

func isUintptr(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.Ident:
		if n.Name == "uintptr" {
			return true
		}
	}
	return false
}

func isFunction(n ast.Node) bool {
	switch n.(type) {
	case *ast.FuncDecl:
		return true
	case *ast.FuncLit:
		return false // deliberately ignore here
	}
	return false
}

func isStatement(n ast.Node) bool {
	switch n.(type) {
	case *ast.AssignStmt:
		return true
	case *ast.BranchStmt:
		return true
	case *ast.DeclStmt:
		return true
	case *ast.EmptyStmt:
		return true
	case *ast.ExprStmt:
		return true
	case *ast.ForStmt:
		return true
	case *ast.GoStmt:
		return true
	case *ast.IfStmt:
		return true
	case *ast.IncDecStmt:
		return true
	case *ast.LabeledStmt:
		return true
	case *ast.RangeStmt:
		return true
	case *ast.ReturnStmt:
		return true
	case *ast.SelectStmt:
		return true
	case *ast.SendStmt:
		return true
	case *ast.SwitchStmt:
		return true
	case *ast.TypeSwitchStmt:
		return true
	}
	return false
}