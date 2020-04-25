package ast

import (
	"bytes"
	"data-aquisition/analysis"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

func printIndent (indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("    ")
	}
}

func printNode(n ast.Node) {
	switch n := n.(type) {
	case *ast.Comment:
		fmt.Printf("Comment")
	case *ast.CommentGroup:
		fmt.Printf("CommentGroup")
	case *ast.Field:
		fmt.Printf("Field")
	case *ast.FieldList:
		fmt.Printf("FieldList")
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
		fmt.Printf("CallExpr %s(...)", n.Fun)
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
		_ = printer.Fprint(buf, fset, arg)
		buf.WriteString(", ")
	}

	return buf.String()
}

func formatFunctions(findingsTree *TreeNode, fset *token.FileSet, lines []string) {
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

func formatStatements(findingsTree *TreeNode, fset *token.FileSet, lines []string) {
	statements := findingsTree.findStatements()
	for s, leaves := range *statements {
		fmt.Printf("STMT ")
		printNode(s.Node)
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

func saveFindings(findingsTree *TreeNode, fset *token.FileSet, lines []string, pkg *analysis.PackageData) {
	for _, finding := range findingsTree.collectLeaves() {
		err := WriteAstFinding(FindingData{
			MatchType:            matchTypeFor(finding),
			LineNumber:           fset.Position(finding.Node.Pos()).Line,
			Column:               fset.Position(finding.Node.Pos()).Column,
			Text:                 lines[fset.Position(finding.Node.Pos()).Line - 1],
			FileName:             fset.Position(finding.Node.Pos()).Filename,
			PackageImportPath:    pkg.ImportPath,
			ModulePath:           pkg.ModulePath,
			ModuleVersion:        pkg.ModuleVersion,
			ProjectName:          pkg.ProjectName,
		})

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "finding-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(finding.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}

	for function := range *findingsTree.findFunctions() {
		startLine := fset.Position(function.Node.Pos()).Line
		endLine := fset.Position(function.Node.End()).Line
		text := strings.Join(lines[startLine-1:endLine], "\n")

		err := WriteFunction(FunctionData{
			LineNumber:           fset.Position(function.Node.Pos()).Line,
			Column:               fset.Position(function.Node.Pos()).Column,
			Text:                 text,
			NumberUnsafePointer:  function.UnsafePointerCount,
			NumberUnsafeSizeof:   function.UnsafeSizeofCount,
			NumberUnsafeAlignof:  function.UnsafeAlignofCount,
			NumberUnsafeOffsetof: function.UnsafeOffsetOfCount,
			NumberUintptr:        function.UintptrCount,
			NumberSliceHeader:    function.SliceHeaderCount,
			NumberStringHeader:   function.StringHeaderCount,
			FileName:             fset.Position(function.Node.Pos()).Filename,
			PackageImportPath:    pkg.ImportPath,
			ModulePath:           pkg.ModulePath,
			ModuleVersion:        pkg.ModuleVersion,
			ProjectName:          pkg.ProjectName,
		})
		
		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "function-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(function.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}

	for statement := range *findingsTree.findStatements() {
		startLine := fset.Position(statement.Node.Pos()).Line
		endLine := fset.Position(statement.Node.End()).Line
		text := strings.Join(lines[startLine-1:endLine], "\n")

		err := WriteStatement(StatementData{
			LineNumber:           fset.Position(statement.Node.Pos()).Line,
			Column:               fset.Position(statement.Node.Pos()).Column,
			Text:                 text,
			NumberUnsafePointer:  statement.UnsafePointerCount,
			NumberUnsafeSizeof:   statement.UnsafeSizeofCount,
			NumberUnsafeAlignof:  statement.UnsafeAlignofCount,
			NumberUnsafeOffsetof: statement.UnsafeOffsetOfCount,
			NumberUintptr:        statement.UintptrCount,
			NumberSliceHeader:    statement.SliceHeaderCount,
			NumberStringHeader:   statement.StringHeaderCount,
			FileName:             fset.Position(statement.Node.Pos()).Filename,
			PackageImportPath:    pkg.ImportPath,
			ModulePath:           pkg.ModulePath,
			ModuleVersion:        pkg.ModuleVersion,
			ProjectName:          pkg.ProjectName,
		})

		if err != nil {
			_ = WriteErrorCondition(ErrorConditionData{
				Stage:             "statement-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(statement.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}
}

func matchTypeFor(n *TreeNode) string {
	switch {
	case isUnsafePointer(n.Node):
		return "unsafe.Pointer"
	case isUnsafeSizeof(n.Node):
		return "unsafe.Sizeof"
	case isUnsafeAlignof(n.Node):
		return "unsafe.Alignof"
	case isUnsafeOffsetof(n.Node):
		return "unsafe.Offsetof"
	case isUintptr(n.Node):
		return "uintptr"
	case isSliceHeader(n.Node):
		return "reflect.SliceHeader"
	case isStringHeader(n.Node):
		return "reflect.StringHeader"
	default:
		return "unknown"
	}
}