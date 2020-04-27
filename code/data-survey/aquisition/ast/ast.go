package ast

import (
	"data-aquisition/lexical"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"strings"
)

func AnalyzeAstSingleFile(mode, filename string, pkg *lexical.PackageData) {
	code, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(code), "\n")

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, path.Base(filename), code, parser.ParseComments)

	findingsTree := &TreeNode{
		Node:     nil,
		Children: []*TreeNode{},
	}

	ast.Walk(UnsafeVisitor{
		fileset:      fset,
		context:      []ast.Node{},
		findingsTree: findingsTree,
	}, node)

	findingsTree.countFindings()

	switch mode {
	case "tree":
		findingsTree.printRoot(fset)
	case "func":
		formatFunctions(findingsTree, fset, lines)
	case "stmt":
		formatStatements(findingsTree, fset, lines)
	case "save":
		saveFindings(findingsTree, fset, lines, pkg)
	default:
		fmt.Printf("unknown mode %s not in tree,func,stmt,save\n", mode)
	}
}

type UnsafeVisitor struct {
	fileset *token.FileSet
	context []ast.Node
	findingsTree *TreeNode
}

func (uv UnsafeVisitor) Visit(n ast.Node) ast.Visitor {
	if isUnsafePointer(n) || isUintptr(n) {
		_ = uv.findingsTree.addPath(append(uv.context, n))
	}

	return UnsafeVisitor{
		fileset: uv.fileset,
		context: append(uv.context, n),
		findingsTree: uv.findingsTree,
	}
}