package ast

import (
	"go/ast"
)

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

func isStatement(n ast.Node) bool {
	switch n.(type) {
	case *ast.AssignStmt:
		return true
	case *ast.BranchStmt:
		return true
	case *ast.DeclStmt:
		return true
	case *ast.EmptyStmt:
		return true
	case *ast.ExprStmt:
		return true
	case *ast.ForStmt:
		return true
	case *ast.GoStmt:
		return true
	case *ast.IfStmt:
		return true
	case *ast.IncDecStmt:
		return true
	case *ast.LabeledStmt:
		return true
	case *ast.RangeStmt:
		return true
	case *ast.ReturnStmt:
		return true
	case *ast.SelectStmt:
		return true
	case *ast.SendStmt:
		return true
	case *ast.SwitchStmt:
		return true
	case *ast.TypeSwitchStmt:
		return true
	}
	return false
}
