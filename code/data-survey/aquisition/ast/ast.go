package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
)

func AnalyzeAst() {
	filename := "/home/johannes/studium/s14/masterarbeit/download/bosun/vendor/github.com/bradfitz/slice/slice.go"
	code, _ := ioutil.ReadFile(filename)

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "slice.go", string(code), parser.ParseComments)

	/*for k, v := range node.Scope.Objects {
		fmt.Println(k)
		fmt.Println(v.Kind.String())
		fmt.Println(v.Decl)
		fmt.Println("----")
	}

	fmt.Println("-----------------")

	for _, k := range node.Decls {
		fmt.Println(k)
	}

	fmt.Println("-----------------")
	fmt.Println("Functions:")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}*/

	/*fmt.Println("-----------------")
	ast.Inspect(node, func (n ast.Node) bool {
		// Find unsafe identifier
		ret, ok := n.(*ast.Ident)
		if ok && ret.Name == "unsafe" {
			//fmt.Printf("unsafe identifier found on line %d:\n", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			fmt.Println()
			//fmt.Println()
			return true
		}
		return true
	})*/

	/*fmt.Println("-----------------")
	ast.Inspect(node, func (n ast.Node) bool {
		// Find unsafe identifier
		ret, ok := n.(*ast.SelectorExpr)
		if ok {
			lhs, ok := ret.X.(*ast.Ident)
			if ok && lhs.Name == "unsafe" {
				fmt.Printf("unsafe identifier found on line %d:\n", fset.Position(ret.Pos()).Line)
				printer.Fprint(os.Stdout, fset, ret)
				fmt.Println()
				//fmt.Println()
			}
			return true
		}
		return true
	})*/

	ast.Walk(UnsafeVisitor{fileset: fset}, node)
	//ast.Print(fset, node)
}

type UnsafeVisitor struct {
	fileset *token.FileSet
	context []ast.Node
}

func (uv UnsafeVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch X := n.X.(type) {
		case *ast.Ident:
			if X.Name == "unsafe" {
				unrollContext(uv, n)
			}
		}
	case *ast.Ident:
		if n.Name == "uintptr" {
			unrollContext(uv, n)
		}
	}

	return UnsafeVisitor{
		fileset: uv.fileset,
		context: append(uv.context, n),
	}
}

func unrollContext(uv UnsafeVisitor, n ast.Node) {
	for i, cn := range uv.context {
		printIndent(i)
		printNode(cn, uv.fileset)
	}
	printIndent(len(uv.context) + 1)
	fmt.Println(n)
	fmt.Println("---")
}

func printNode(n ast.Node, fset *token.FileSet) {
	switch n := n.(type) {
	case *ast.Comment:
		fmt.Printf("Comment\n")
	case *ast.CommentGroup:
		fmt.Printf("CommentGroup\n")
	case *ast.Field:
		fmt.Printf("Field: %s\n", n.Names[0].Name)
	case *ast.FieldList:
		fmt.Printf("FieldList: %s\n", n.List[0].Names[0].Name)
	case *ast.BadExpr:
		fmt.Printf("BadExpr\n")
	case *ast.Ident:
		fmt.Printf("Ident: %s\n", n.Name)
	case *ast.BasicLit:
		fmt.Printf("BasicLit: %s\n", n.Value)
	case *ast.Ellipsis:
		fmt.Printf("Ellipsis\n")
	case *ast.FuncLit:
		fmt.Printf("FuncLit\n")
	case *ast.CompositeLit:
		fmt.Printf("CompositeLit\n")
	case *ast.ParenExpr:
		fmt.Printf("ParenExpr\n")
	case *ast.SelectorExpr:
		fmt.Printf("SelectorExpr\n")
	case *ast.IndexExpr:
		fmt.Printf("IndexExpr\n")
	case *ast.SliceExpr:
		fmt.Printf("SliceExpr\n")
	case *ast.TypeAssertExpr:
		fmt.Printf("TypeAssertExpr\n")
	case *ast.CallExpr:
		fmt.Printf("CallExpr %s(%s)\n", n.Fun, argumentsToString(n.Args, fset))
	case *ast.StarExpr:
		fmt.Printf("StarExpr\n")
	case *ast.UnaryExpr:
		fmt.Printf("UnaryExpr\n")
	case *ast.BinaryExpr:
		fmt.Printf("BinaryExpr\n")
	case *ast.KeyValueExpr:
		fmt.Printf("KeyValueExpr\n")
	case *ast.ArrayType:
		fmt.Printf("ArrayType\n")
	case *ast.StructType:
		fmt.Printf("StructType\n")
	case *ast.FuncType:
		fmt.Printf("FuncType\n")
	case *ast.InterfaceType:
		fmt.Printf("InterfaceType\n")
	case *ast.MapType:
		fmt.Printf("MapType\n")
	case *ast.ChanType:
		fmt.Printf("ChanType\n")
	case *ast.BadStmt:
		fmt.Printf("BadStmt\n")
	case *ast.DeclStmt:
		fmt.Printf("DeclStmt\n")
	case *ast.EmptyStmt:
		fmt.Printf("EmptyStmt\n")
	case *ast.LabeledStmt:
		fmt.Printf("LabeledStmt\n")
	case *ast.ExprStmt:
		fmt.Printf("ExprStmt\n")
	case *ast.SendStmt:
		fmt.Printf("SendStmt\n")
	case *ast.IncDecStmt:
		fmt.Printf("IncDecStmt:\n")
	case *ast.AssignStmt:
		fmt.Printf("AssignStmt: %s\n", n.Tok.String())
	case *ast.GoStmt:
		fmt.Printf("GoStmt\n")
	case *ast.DeferStmt:
		fmt.Printf("DeferStmt\n")
	case *ast.ReturnStmt:
		fmt.Printf("ReturnStmt\n")
	case *ast.BranchStmt:
		fmt.Printf("BranchStmt: %s\n", n.Tok.String())
	case *ast.BlockStmt:
		fmt.Printf("BlockStmt\n")
	case *ast.IfStmt:
		fmt.Printf("IfStmt\n")
	case *ast.CaseClause:
		fmt.Printf("CaseClause\n")
	case *ast.SwitchStmt:
		fmt.Printf("SwitchStmt\n")
	case *ast.TypeSwitchStmt:
		fmt.Printf("TypeSwitchStmt\n")
	case *ast.CommClause:
		fmt.Printf("CommClause\n")
	case *ast.SelectStmt:
		fmt.Printf("SelectStmt\n")
	case *ast.ForStmt:
		fmt.Printf("ForStmt\n")
	case *ast.RangeStmt:
		fmt.Printf("RangeStmt\n")
	case *ast.ImportSpec:
		fmt.Printf("ImportSpec: %s\n", n.Name)
	case *ast.ValueSpec:
		fmt.Printf("ValueSpec: %s\n", n.Names[0].Name)
	case *ast.TypeSpec:
		fmt.Printf("TypeSepc: %s\n", n.Name)
	case *ast.BadDecl:
		fmt.Printf("BadDecl\n")
	case *ast.GenDecl:
		fmt.Printf("GenDecl: %s\n", n.Tok)
	case *ast.FuncDecl:
		fmt.Printf("FunDecl: %s\n", n.Name)
	case *ast.File:
		fmt.Printf("File: %s\n", n.Name)
	case *ast.Package:
		fmt.Printf("Package: %s\n", n.Name)
	default:
		fmt.Printf("unknown node type %T\n", n)
	}
}

func printIndent (indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
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