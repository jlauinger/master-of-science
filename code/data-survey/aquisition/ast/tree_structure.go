package ast

import (
	"fmt"
	"go/ast"
	"go/token"
)

type TreeNode struct {
	Node                ast.Node
	Children            []*TreeNode

	UnsafePointerCount  int
	UnsafeSizeofCount   int
	UnsafeAlignofCount  int
	UnsafeOffsetOfCount int
	UintptrCount        int
	SliceHeaderCount    int
	StringHeaderCount   int
}

func (t *TreeNode) addPath(path []ast.Node) error {
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

	var childToAdd *TreeNode
	for _, child := range t.Children {
		if child.Node == head {
			childToAdd = child
		}
	}
	if childToAdd == nil {
		childToAdd = &TreeNode{
			Node:     head,
			Children: []*TreeNode{},
		}
		t.Children = append(t.Children, childToAdd)
	}

	err := childToAdd.addPath(tail)
	if err != nil {
		return err
	}

	return nil
}

func (t *TreeNode) printRoot(fset *token.FileSet) {
	for _, child := range t.Children {
		printIter(child, fset, 0)
	}
}

func printIter(t *TreeNode, fset *token.FileSet, indent int) {
	printIndent(indent)
	printNode(t.Node)
	fmt.Printf(" (%d, %d)\n", t.UnsafePointerCount, t.UintptrCount)

	for _, child := range t.Children {
		printIter(child, fset, indent + 1)
	}
}

func (t *TreeNode) countFindings() {
	t.countUnsafePointer()
	t.countUnsafeSizeof()
	t.countUnsafeAlignof()
	t.countUnsafeOffsetof()
	t.countUintptr()
	t.countSliceHeader()
	t.countStringHeader()
}

func (t *TreeNode) countUnsafePointer() int {
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
func (t *TreeNode) countUnsafeSizeof() int {
	count := 0
	for _, child := range t.Children {
		count += child.countUnsafeSizeof()
	}
	if isUnsafeSizeof(t.Node) {
		count++
	}
	t.UnsafeSizeofCount = count
	return count
}
func (t *TreeNode) countUnsafeAlignof() int {
	count := 0
	for _, child := range t.Children {
		count += child.countUnsafeAlignof()
	}
	if isUnsafeAlignof(t.Node) {
		count++
	}
	t.UnsafeAlignofCount = count
	return count
}
func (t *TreeNode) countUnsafeOffsetof() int {
	count := 0
	for _, child := range t.Children {
		count += child.countUnsafeOffsetof()
	}
	if isUnsafeOffsetof(t.Node) {
		count++
	}
	t.UnsafeOffsetOfCount = count
	return count
}
func (t *TreeNode) countUintptr() int {
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
func (t *TreeNode) countSliceHeader() int {
	count := 0
	for _, child := range t.Children {
		count += child.countSliceHeader()
	}
	if isSliceHeader(t.Node) {
		count++
	}
	t.SliceHeaderCount = count
	return count
}
func (t *TreeNode) countStringHeader() int {
	count := 0
	for _, child := range t.Children {
		count += child.countStringHeader()
	}
	if isStringHeader(t.Node) {
		count++
	}
	t.StringHeaderCount = count
	return count
}

func (t *TreeNode) collectLeaves() []*TreeNode {
	if len(t.Children) == 0 {
		return []*TreeNode{}
	}

	return collectLeavesIter(t)
}

func collectLeavesIter(t *TreeNode) []*TreeNode {
	if len(t.Children) == 0 {
		return []*TreeNode{t}
	}

	var leaves []*TreeNode
	for _, child := range t.Children {
		leaves = append(leaves, collectLeavesIter(child)...)
	}
	return leaves
}

func (t *TreeNode) findFunctions() *map[*TreeNode][]*TreeNode {
	return findFunctionsIter(t, &map[*TreeNode][]*TreeNode{})
}

func findFunctionsIter(t *TreeNode, functions *map[*TreeNode][]*TreeNode) *map[*TreeNode][]*TreeNode {
	if isFunction(t.Node) {
		(*functions)[t] = t.collectLeaves()
	} else {
		for _, child := range t.Children {
			findFunctionsIter(child, functions)
		}
	}
	return functions
}

func (t *TreeNode) findStatements() *map[*TreeNode][]*TreeNode {
	return findStatementsIter(t, &map[*TreeNode][]*TreeNode{})
}

func findStatementsIter(t *TreeNode, statements *map[*TreeNode][]*TreeNode) *map[*TreeNode][]*TreeNode {
	shouldCollectLeaves := false
	if isStatement(t.Node) {
		for _, child := range t.Children {
			if !child.containsStatement() {
				shouldCollectLeaves = true
			}
		}
	}
	if shouldCollectLeaves {
		(*statements)[t] = t.collectLeaves()
	} else {
		for _, child := range t.Children {
			findStatementsIter(child, statements)
		}
	}
	return statements
}

func (t *TreeNode) containsStatement() bool {
	if isStatement(t.Node) {
		return true
	}
	for _, child := range t.Children {
		if child.containsStatement() {
			return true
		}
	}
	return false
}