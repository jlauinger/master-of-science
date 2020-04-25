package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

func printIndent (indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("    ")
	}
}

func printNode(n ast.Node, fset *token.FileSet) {
	switch n := n.(type) {
	case *ast.Comment:
		fmt.Printf("Comment")
	case *ast.CommentGroup:
		fmt.Printf("CommentGroup")
	case *ast.Field:
		fmt.Printf("Field: %s", n.Names[0].Name)
	case *ast.FieldList:
		fmt.Printf("FieldList: %s", n.List[0].Names[0].Name)
	case *ast.BadExpr:
		fmt.Printf("BadExpr")
	case *ast.Ident:
		fmt.Printf("Ident: %s", n.Name)
	case *ast.BasicLit:
		fmt.Printf("BasicLit: %s", n.Value)
	case *ast.Ellipsis:
		fmt.Printf("Ellipsis")
	case *ast.FuncLit:
		fmt.Printf("FuncLit")
	case *ast.CompositeLit:
		fmt.Printf("CompositeLit")
	case *ast.ParenExpr:
		fmt.Printf("ParenExpr")
	case *ast.SelectorExpr:
		fmt.Printf("SelectorExpr")
	case *ast.IndexExpr:
		fmt.Printf("IndexExpr")
	case *ast.SliceExpr:
		fmt.Printf("SliceExpr")
	case *ast.TypeAssertExpr:
		fmt.Printf("TypeAssertExpr")
	case *ast.CallExpr:
		fmt.Printf("CallExpr %s(...)", n.Fun) //, argumentsToString(n.Args, fset))
	case *ast.StarExpr:
		fmt.Printf("StarExpr")
	case *ast.UnaryExpr:
		fmt.Printf("UnaryExpr")
	case *ast.BinaryExpr:
		fmt.Printf("BinaryExpr")
	case *ast.KeyValueExpr:
		fmt.Printf("KeyValueExpr")
	case *ast.ArrayType:
		fmt.Printf("ArrayType")
	case *ast.StructType:
		fmt.Printf("StructType")
	case *ast.FuncType:
		fmt.Printf("FuncType")
	case *ast.InterfaceType:
		fmt.Printf("InterfaceType")
	case *ast.MapType:
		fmt.Printf("MapType")
	case *ast.ChanType:
		fmt.Printf("ChanType")
	case *ast.BadStmt:
		fmt.Printf("BadStmt")
	case *ast.DeclStmt:
		fmt.Printf("DeclStmt")
	case *ast.EmptyStmt:
		fmt.Printf("EmptyStmt")
	case *ast.LabeledStmt:
		fmt.Printf("LabeledStmt")
	case *ast.ExprStmt:
		fmt.Printf("ExprStmt")
	case *ast.SendStmt:
		fmt.Printf("SendStmt")
	case *ast.IncDecStmt:
		fmt.Printf("IncDecStmt:")
	case *ast.AssignStmt:
		fmt.Printf("AssignStmt: %s", n.Tok.String())
	case *ast.GoStmt:
		fmt.Printf("GoStmt")
	case *ast.DeferStmt:
		fmt.Printf("DeferStmt")
	case *ast.ReturnStmt:
		fmt.Printf("ReturnStmt")
	case *ast.BranchStmt:
		fmt.Printf("BranchStmt: %s", n.Tok.String())
	case *ast.BlockStmt:
		fmt.Printf("BlockStmt")
	case *ast.IfStmt:
		fmt.Printf("IfStmt")
	case *ast.CaseClause:
		fmt.Printf("CaseClause")
	case *ast.SwitchStmt:
		fmt.Printf("SwitchStmt")
	case *ast.TypeSwitchStmt:
		fmt.Printf("TypeSwitchStmt")
	case *ast.CommClause:
		fmt.Printf("CommClause")
	case *ast.SelectStmt:
		fmt.Printf("SelectStmt")
	case *ast.ForStmt:
		fmt.Printf("ForStmt")
	case *ast.RangeStmt:
		fmt.Printf("RangeStmt")
	case *ast.ImportSpec:
		fmt.Printf("ImportSpec: %s", n.Name)
	case *ast.ValueSpec:
		fmt.Printf("ValueSpec: %s", n.Names[0].Name)
	case *ast.TypeSpec:
		fmt.Printf("TypeSepc: %s", n.Name)
	case *ast.BadDecl:
		fmt.Printf("BadDecl")
	case *ast.GenDecl:
		fmt.Printf("GenDecl: %s", n.Tok)
	case *ast.FuncDecl:
		fmt.Printf("FunDecl: %s", n.Name)
	case *ast.File:
		fmt.Printf("File: %s", n.Name)
	case *ast.Package:
		fmt.Printf("Package: %s", n.Name)
	default:
		fmt.Printf("unknown node type %T", n)
	}
}

func argumentsToString(args []ast.Expr, fset *token.FileSet) string {
	buf := bytes.NewBufferString("")

	for _, arg := range args {
		printer.Fprint(buf, fset, arg)
		buf.WriteString(", ")
	}

	return buf.String()
}

func formatFunctions(findingsTree *AstTreeNode, fset *token.FileSet, lines []string) {
	functions := findingsTree.findFunctions()
	for f, leaves := range *functions {
		fmt.Printf("FUNC %s (%d, %d):\n", f.Node.(*ast.FuncDecl).Name, f.UnsafePointerCount, f.UintptrCount)
		for _, leaf := range leaves {
			printIter(leaf, fset, 1)
			printIndent(2)
			fmt.Println(fset.Position(leaf.Node.Pos()))
			printIndent(2)
			fmt.Println(lines[fset.Position(leaf.Node.Pos()).Line-1])
		}
		fmt.Println()
	}
}

func formatStatements(findingsTree *AstTreeNode, fset *token.FileSet, lines []string) {
	statements := findingsTree.findStatements()
	for s, leaves := range *statements {
		fmt.Printf("STMT ")
		printNode(s.Node, fset)
		fmt.Printf(" (%d, %d):\n", s.UnsafePointerCount, s.UintptrCount)
		for _, leaf := range leaves {
			printIter(leaf, fset, 1)
			printIndent(2)
			fmt.Println(fset.Position(leaf.Node.Pos()))
			printIndent(2)
			fmt.Println(lines[fset.Position(leaf.Node.Pos()).Line-1])
		}
		fmt.Println()
	}
}