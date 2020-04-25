package ast

import (
	"go/ast"
)

func isUnsafePointer(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" && n.Sel.Name == "Pointer" {
				return true
			}
		}
	}
	return false
}

func isUnsafeSizeof(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" && n.Sel.Name == "Sizeof" {
				return true
			}
		}
	}
	return false
}

func isUnsafeAlignof(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" && n.Sel.Name == "Alignof" {
				return true
			}
		}
	}
	return false
}

func isUnsafeOffsetof(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" && n.Sel.Name == "Offsetof" {
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

func isSliceHeader(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "reflect" && n.Sel.Name == "SliceHeader" {
				return true
			}
		}
	}
	return false
}

func isStringHeader(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "reflect" && n.Sel.Name == "StringHeader" {
				return true
			}
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
	case *ast.DeclStmt:
		return true
	case *ast.EmptyStmt:
		return true
	case *ast.ExprStmt:
		return true
	case *ast.ForStmt:
		return false
	case *ast.GoStmt:
		return true
	case *ast.IfStmt:
		return true
	case *ast.IncDecStmt:
		return true
	case *ast.LabeledStmt:
		return true
	case *ast.RangeStmt:
		return false
	case *ast.ReturnStmt:
		return true
	case *ast.SelectStmt:
		return true
	case *ast.SendStmt:
		return true
	case *ast.SwitchStmt:
		return false
	case *ast.TypeSwitchStmt:
		return false
	}
	return false
}
