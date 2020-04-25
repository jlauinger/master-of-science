package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func AnalyzeAst() {
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

	//findingsTree.printRoot(fset)
	functions := findingsTree.findFunctions()
	
	for f, leaves := range *functions {
		fmt.Printf("FUNC %s (%d, %d):\n", f.Node.(*ast.FuncDecl).Name, f.UnsafePointerCount, f.UintptrCount)
		for _, leaf := range leaves {
			printIter(leaf, fset, 1)
			printIndent(2)
			fmt.Println(fset.Position(leaf.Node.Pos()))
			printIndent(2)
			fmt.Println(lines[fset.Position(leaf.Node.Pos()).Line-1])
		}
		fmt.Println()
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