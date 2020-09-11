package ast

import (
	"go/ast"
)

/**
 * TreeNode represents a node in the findings tree. It contains the count of all unsafe call sites in its children, as
 * well as the children itself. The node is a Go AST node, but the children are recursive TreeNode instances. Therefore,
 * this acts as an augmented AST tree structure with my findings counts
 */
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

/**
 * adds a path (a context stack of AST nodes) to the findings tree. This is done if the top-most node in the stack to
 * be added is an unsafe findings, so the findings tree always has findings in its leaves
 */
func (t *TreeNode) addPath(path []ast.Node) error {
	// if the path is empty, do nothing. This terminates the recursion
	if len(path) == 0 {
		return nil
	}

	// separate head and tail from the stack, taking care of the possibility that there is no tail (any more)
	var head ast.Node
	var tail []ast.Node
	if len(path) == 1 {
		head, tail = path[0], []ast.Node{}
	} else {
		head, tail = path[0], path[1:]
	}

	// check if the head of the stack to be added is already a child of this findings tree node
	var childToAdd *TreeNode
	for _, child := range t.Children {
		if child.Node == head {
			// if so, recursively add the path onto this branch of the findings tree instead of adding a new child
			// this will weave in the new path into the existing findings tree
			childToAdd = child
		}
	}
	// if the head of the stack to be added is not already a child
	if childToAdd == nil {
		// then create a new child that represents that head node
		childToAdd = &TreeNode{
			Node:     head,
			Children: []*TreeNode{},
		}
		// and append it to the children of this findings tree node
		t.Children = append(t.Children, childToAdd)
	}

	// then recursively weave in the rest of the stack, returning back any potential error
	return childToAdd.addPath(tail)
}

/**
 * populates all the findings counts fields in the tree
 */
func (t *TreeNode) countFindings() {
	// each count is counted separately, so count all of them
	t.countUnsafePointer()
	t.countUnsafeSizeof()
	t.countUnsafeAlignof()
	t.countUnsafeOffsetof()
	t.countUintptr()
	t.countSliceHeader()
	t.countStringHeader()
}

/**
 * populates the unsafe pointer counts fields
 */
func (t *TreeNode) countUnsafePointer() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countUnsafePointer()
	}
	// if this node is an unsafe pointer itself, add 1
	if isUnsafePointer(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.UnsafePointerCount = count
	return count
}
/**
 * populates the unsafe sizeof counts fields
 */
func (t *TreeNode) countUnsafeSizeof() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countUnsafeSizeof()
	}
	// if this node is an unsafe sizeof itself, add 1
	if isUnsafeSizeof(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.UnsafeSizeofCount = count
	return count
}
/**
 * populates the unsafe alignof counts fields
 */
func (t *TreeNode) countUnsafeAlignof() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countUnsafeAlignof()
	}
	// if this node is an unsafe alignof itself, add 1
	if isUnsafeAlignof(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.UnsafeAlignofCount = count
	return count
}
/**
 * populates the unsafe offsetof counts fields
 */
func (t *TreeNode) countUnsafeOffsetof() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countUnsafeOffsetof()
	}
	// if this node is an unsafe offsetof itself, add 1
	if isUnsafeOffsetof(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.UnsafeOffsetOfCount = count
	return count
}
/**
 * populates the uintptr counts fields
 */
func (t *TreeNode) countUintptr() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countUintptr()
	}
	// if this node is a uintptr itself, add 1
	if isUintptr(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.UintptrCount = count
	return count
}
/**
 * populates the reflect sliceheader counts fields
 */
func (t *TreeNode) countSliceHeader() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countSliceHeader()
	}
	// if this node is a reflect sliceheader itself, add 1
	if isSliceHeader(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.SliceHeaderCount = count
	return count
}
/**
 * populates the reflect stringheader counts fields
 */
func (t *TreeNode) countStringHeader() int {
	count := 0
	// add all counts form the children recursively
	for _, child := range t.Children {
		count += child.countStringHeader()
	}
	// if this node is a reflect stringheader itself, add 1
	if isStringHeader(t.Node) {
		count++
	}
	// store the count in this node and return it for recursive calculation
	t.StringHeaderCount = count
	return count
}

/**
 * returns all leaves of the tree
 */
func (t *TreeNode) collectLeaves() []*TreeNode {
	// if there are no children, there are no leaves at all. This special case is needed because the root of my
	// findings tree is always there but does not represent a leaf
	if len(t.Children) == 0 {
		return []*TreeNode{}
	}

	// otherwise, collect the leaves starting with the
	return collectLeavesRecursive(t)
}

/**
 * returns all leaves of the given tree node. Used internally for recursion
 */
func collectLeavesRecursive(t *TreeNode) []*TreeNode {
	// if there are no children, this is a leaf, so return it
	if len(t.Children) == 0 {
		return []*TreeNode{t}
	}

	// otherwise, go through all children, recursively collect their leaves and splice them together
	var leaves []*TreeNode
	for _, child := range t.Children {
		leaves = append(leaves, collectLeavesRecursive(child)...)
	}
	return leaves
}

/**
 * returns all function nodes in the findings tree
 */
func (t *TreeNode) findFunctions() *map[*TreeNode][]*TreeNode {
	// recursively start with the tree root and no functions
	return findFunctionsIter(t, &map[*TreeNode][]*TreeNode{})
}

/**
 * internal recursion function to get functions in the given tree node
 */
func findFunctionsIter(t *TreeNode, functions *map[*TreeNode][]*TreeNode) *map[*TreeNode][]*TreeNode {
	// if this node is a function, add it to the functions map using the function as key and its leaves (findings) as
	// values.
	if isFunction(t.Node) {
		(*functions)[t] = t.collectLeaves()
	} else {
		// otherwise, call the function finder recursively on all children
		for _, child := range t.Children {
			findFunctionsIter(child, functions)
		}
	}
	// return the populated functions map
	return functions
}

/** returns all statement nodes in the findings tree
 */
func (t *TreeNode) findStatements() *map[*TreeNode][]*TreeNode {
	// recursively start with the tree root and no statements
	return findStatementsIter(t, &map[*TreeNode][]*TreeNode{})
}

/**
 * internal recursion function to get statements in the given tree node
 */
func findStatementsIter(t *TreeNode, statements *map[*TreeNode][]*TreeNode) *map[*TreeNode][]*TreeNode {
	// for statements, we need the inner most statement that does not contain any more sub statements. There can be
	// made a choice which should count as a statement that does not get divided anymore (e.g. assignments, or if, or
	// for, ...) and that choice is made through the isStatement function
	shouldCollectLeaves := false
	// check if this node is a statement at all
	if isStatement(t.Node) {
		// then check if the children contain more statements. If there is at least one branch that does not contain
		// more sub statements, this node here should be collected as the statement found
		for _, child := range t.Children {
			if !child.containsStatement() {
				shouldCollectLeaves = true
			}
		}
	}
	// if this node should be collected, save it to the statements hash map using the statement node as key and the
	// leaves (findings) as values
	if shouldCollectLeaves {
		(*statements)[t] = t.collectLeaves()
	} else {
		// otherwise, call the statement finder recursively on all children
		for _, child := range t.Children {
			findStatementsIter(child, statements)
		}
	}
	// finally return the populated statements map
	return statements
}

/**
 * checks if a tree node recursively contains any statement node
 */
func (t *TreeNode) containsStatement() bool {
	// if this node is a statement itself, return true immediately
	if isStatement(t.Node) {
		return true
	}
	// otherwise, check all children recursively
	for _, child := range t.Children {
		if child.containsStatement() {
			return true
		}
	}
	// if the node itself and none of the children contain statements, return false
	return false
}