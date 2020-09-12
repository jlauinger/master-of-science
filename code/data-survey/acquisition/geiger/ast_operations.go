package geiger

import "go/ast"

/**
 * returns true if the given AST node is an unsafe pointer node
 */
func isUnsafePointer(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Pointer" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is an unsafe sizeof node
 */
func isUnsafeSizeof(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Sizeof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is an unsafe offsetof node
 */
func isUnsafeOffsetof(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Offsetof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is an unsafe alignof node
 */
func isUnsafeAlignof(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Alignof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is a reflect stringheader node
 */
func isReflectStringHeader(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "StringHeader" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is a reflect sliceheader node
 */
func isReflectSliceHeader(node *ast.SelectorExpr) bool {
	// check if the node is a selector expression with matching names
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "SliceHeader" {
			return true
		}
	}
	return false
}

/**
 * returns true if the given AST node is a uintptr node
 */
func isUintptr(node *ast.Ident) bool {
	// check if the identified name matches
	return node.Name == "uintptr"
}

/**
 * returns true if the given AST node stack represents an unsafe finding as an argument in a function call
 */
func isArgument(stack []ast.Node) bool {
	// skip the last stack elements because the unsafe.Pointer SelectorExpr is itself a call expression.
	// the selector expression is in function position of a call, and we are not interested in that.
	// go through the stack backwards because we are interested in the innermost context occurrence
	for i := len(stack) - 2; i > 0; i-- {
		n := stack[i - 1]
		// if there is a CallExpr node, this is a function call argument
		_, ok := n.(*ast.CallExpr)
		if ok {
			return true
		}
	}
	// otherwise, it is something else
	return false
}

/**
 * returns true if the given AST node stack represents an unsafe finding as an assignment
 */
func isInAssignment(stack []ast.Node) bool {
	// go through the stack backwards because we are interested in the innermost context occurrence
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// if there is an assignment statement, composite literal or return statement we count this as an assignment
		_, ok := n.(*ast.AssignStmt)
		if ok {
			return true
		}
		_, ok = n.(*ast.CompositeLit)
		if ok {
			return true
		}
		_, ok = n.(*ast.ReturnStmt)
		if ok {
			return true
		}
	}
	// otherwise, it is something else
	return false
}

/**
 * returns true if the given AST node stack represents an unsafe finding as a parameter in a function definition
 */
func isParameter(stack []ast.Node) bool {
	// go through the stack backwards because we are interested in the innermost context occurrence
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// if there is a function type node, this is a function definition
		_, ok := n.(*ast.FuncType)
		if ok {
			return true
		}
	}
	// otherwise, it is something else
	return false
}

/**
 * returns true if the given AST node stack represents an unsafe finding as a variable definition
 */
func isInVariableDefinition(stack []ast.Node) bool {
	// go through the stack backwards because we are interested in the innermost context occurrence
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// if there is a generic declaration node, this is a variable definition
		_, ok := n.(*ast.GenDecl)
		if ok {
			return true
		}
	}
	// otherwise, it is something else
	return false
}
