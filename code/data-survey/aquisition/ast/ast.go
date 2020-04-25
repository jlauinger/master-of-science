package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func AnalyzeAst(mode, filename string) {
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