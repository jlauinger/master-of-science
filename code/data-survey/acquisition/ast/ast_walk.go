package ast

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"strings"
)

/**
 * analyzes a single file using the AST analysis step
 */
func AnalyzeAstSingleFile(mode, filename string, pkg *base.PackageData) {
	// read in the source code and split by lines (the lines list is needed later to easily get the context lines for a
	// finding)
	code, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(code), "\n")

	// initialize a token file set and parse the file with it
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, path.Base(filename), code, parser.ParseComments)

	// initialize an empty TreeNode structure which will be the root of the tree of my findings
	findingsTree := &TreeNode{
		Node:     nil,
		Children: []*TreeNode{},
	}

	// go through the AST with a custom visitor object that identified unsafe usages
	ast.Walk(UnsafeVisitor{
		fileset:      fset,
		context:      []ast.Node{},
		findingsTree: findingsTree,
	}, node)

	// on the finished findings tree, count the findings
	findingsTree.countFindings()

	// the AST analysis can be used to print a tree of usages as well as counts by function and statement. Most
	// importantly however, the save mode saves findings to disk. Depending on the requested mode, call the right
	// function
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

/**
 * this visitor structure can hold a current node stack called context as well as put findings into a new tree structure
 * for easy analysis later on
 */
type UnsafeVisitor struct {
	fileset *token.FileSet
	context []ast.Node
	findingsTree *TreeNode
}

/**
 * the Visit function will be called by the AST Walk operation and identifies unsafe pointers and uintptrs
 */
func (uv UnsafeVisitor) Visit(n ast.Node) ast.Visitor {
	// check if this node is an unsafe.Pointer or uintptr call site
	if isUnsafePointer(n) || isUintptr(n) {
		// if so, append it to the findings tree by weaving in the node stack
		_ = uv.findingsTree.addPath(append(uv.context, n))
	}

	// then return a new visitor object which will recursively be used to dig deeper into the AST. It saves this node
	// as part of the node stack for potential upcoming unsafe nodes in the context field.
	return UnsafeVisitor{
		fileset: uv.fileset,
		context: append(uv.context, n),
		findingsTree: uv.findingsTree,
	}
}