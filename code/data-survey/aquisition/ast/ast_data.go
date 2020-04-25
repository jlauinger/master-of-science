package ast

import (
	"go/ast"
	"go/token"
	"fmt"
)

type AstTreeNode struct {
	Node               ast.Node
	Children           []*AstTreeNode
	UnsafePointerCount int
	UintptrCount       int
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
	fmt.Printf(" (%d, %d)\n", t.UnsafePointerCount, t.UintptrCount)

	for _, child := range t.Children {
		printIter(child, fset, indent + 1)
	}
}

func (t *AstTreeNode) countUnsafePointer() int {
	count := 0
	for _, child := range t.Children {
		count += child.countUnsafePointer()
	}

	if isUnsafePointer(t.Node) {
		count++
	}

	t.UnsafePointerCount = count

	return count
}

func (t *AstTreeNode) countUintptr() int {
	count := 0
	for _, child := range t.Children {
		count += child.countUintptr()
	}

	if isUintptr(t.Node) {
		count++
	}

	t.UintptrCount = count

	return count
}