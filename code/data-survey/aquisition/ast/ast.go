package ast

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"fmt"
	"os"
)

func AnalyzeAst() {
	filename := "/home/johannes/studium/s14/masterarbeit/download/bosun/vendor/github.com/bradfitz/slice/slice.go"
	code, _ := ioutil.ReadFile(filename)

	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "slice.go", string(code), parser.ParseComments)

	for k, v := range node.Scope.Objects {
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
	}

	fmt.Println("-----------------")
	ast.Inspect(node, func (n ast.Node) bool {
		// Find unsafe identifier
		ret, ok := n.(*ast.Ident)
		if ok && ret.Name == "unsafe" {
			fmt.Printf("unsafe identifier found on line %d:\n", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			fmt.Println()
			fmt.Println()
			return true
		}
		return true
	})
}