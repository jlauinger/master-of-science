package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func AnalyzeAst() {
	filename := "/home/johannes/studium/s14/masterarbeit/download/bosun/vendor/github.com/bradfitz/slice/slice.go"
	code, _ := ioutil.ReadFile(filename)

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "slice.go", string(code), parser.ParseComments)

	findingsTree := &AstTreeNode{
		NodeType: NODE_TYPE_ROOT,
		Node:     nil,
		Children: []*AstTreeNode{},
	}

	ast.Walk(UnsafeVisitor{
		fileset:      fset,
		context:      []ast.Node{},
		findingsTree: findingsTree,
	}, node)

	findingsTree.print(fset)
}

type UnsafeVisitor struct {
	fileset *token.FileSet
	context []ast.Node
	findingsTree *AstTreeNode
}

func (uv UnsafeVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" {
				//unrollContext(uv, n)
				uv.findingsTree.addPath(uv.context)
			}
		}
	case *ast.Ident:
		if n.Name == "uintptr" {
			//unrollContext(uv, n)
			uv.findingsTree.addPath(uv.context)
		}
	}

	return UnsafeVisitor{
		fileset: uv.fileset,
		context: append(uv.context, n),
		findingsTree: uv.findingsTree,
	}
}
