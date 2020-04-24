package ast

import (
	"go/ast"
	"go/token"
)

const(
	NODE_TYPE_ROOT = "root"
	NODE_TYPE_FUNCTION = "func"
	NODE_TYPE_STATEMENT = "stmt"
)

type AstTreeNode struct {
	NodeType string
	Node ast.Node
	Children []*AstTreeNode
}

func (t *AstTreeNode) addPath(path []ast.Node) error {
	if len(path) == 0 {
		return nil
	}

	var head ast.Node
	var tail []ast.Node
	if len(path) == 1 {
		head, tail = path[0], []ast.Node{}
	} else {
		head, tail = path[0], path[1:]
	}

	var childToAdd *AstTreeNode
	for _, child := range t.Children {
		if child.Node == head {
			childToAdd = child
		}
	}
	if childToAdd == nil {
		childToAdd = &AstTreeNode{
			NodeType: "foo",
			Node:     head,
			Children: []*AstTreeNode{},
		}
		t.Children = append(t.Children, childToAdd)
	}

	err := childToAdd.addPath(tail)
	if err != nil {
		return err
	}

	return nil
}

func (t *AstTreeNode) print(fset *token.FileSet) {
	for _, child := range t.Children {
		printIter(child, fset, 0)
	}
}

func printIter(t *AstTreeNode, fset *token.FileSet, indent int) {
	printIndent(indent)
	printNode(t.Node, fset)
	for _, child := range t.Children {
		printIter(child, fset, indent + 1)
	}
}
