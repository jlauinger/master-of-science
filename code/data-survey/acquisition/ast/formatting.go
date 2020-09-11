package ast

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"go/ast"
	"go/token"
	"strings"
)

/**
 * prints an appropriate amount of spaces depending on the given indent depth
 */
func printIndent (indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("    ")
	}
}

/**
 * prints a given AST node's rough string representation to stdout
 */
func printNode(n ast.Node) {
	// depending on the node type, output its type name and for some additional information like specific identifier
	// names and such
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

/**
 * prints the root of the findings tree
 */
func (t *TreeNode) printRoot(fset *token.FileSet) {
	// go through all children of the tree
	for _, child := range t.Children {
		// and print those. The root itself is not of interest
		printRecursive(child, fset, 0)
	}
}

/**
 * recursively prints a node in the findings tree
 */
func printRecursive(t *TreeNode, fset *token.FileSet, indent int) {
	// print the correct indent for this depth
	printIndent(indent)
	// print the node itself
	printNode(t.Node)
	// print the finding counts
	fmt.Printf(" (%d, %d)\n", t.UnsafePointerCount, t.UintptrCount)

	// then go through the children of this node
	for _, child := range t.Children {
		// and recursively print those as well
		printRecursive(child, fset, indent + 1)
	}
}

/**
 * prints the functions in the findings tree to stdout
 */
func formatFunctions(findingsTree *TreeNode, fset *token.FileSet, lines []string) {
	// get all functions in the findings tree
	functions := findingsTree.findFunctions()
	// go through the functions, represented by the function and their findings (leaves in the findings tree)
	for f, leaves := range *functions {
		// print the function name and its unsafe findings counts (unsafe pointer and uintptr)
		fmt.Printf("FUNC %s (%d, %d):\n", f.Node.(*ast.FuncDecl).Name, f.UnsafePointerCount, f.UintptrCount)
		// go through all leaves (findings) for this function
		for _, leaf := range leaves {
			// print the findings tree node to stdout
			printRecursive(leaf, fset, 1)
			// then print the code position (filename and line) in a new line with indent 2
			printIndent(2)
			fmt.Println(fset.Position(leaf.Node.Pos()))
			// and print the code line itself as a context
			printIndent(2)
			fmt.Println(lines[fset.Position(leaf.Node.Pos()).Line-1])
		}
		fmt.Println()
	}
}

/**
 * prints the statements in the findings tree to stdout
 */
func formatStatements(findingsTree *TreeNode, fset *token.FileSet, lines []string) {
	// get all statements in the findings tree
	statements := findingsTree.findStatements()
	// go through the statements, represented by the statement and their findings (leaves in the findings tree)
	for s, leaves := range *statements {
		// print the statement and its unsafe findings counts (unsafe pointer and uintptr)
		fmt.Printf("STMT ")
		printNode(s.Node)
		fmt.Printf(" (%d, %d):\n", s.UnsafePointerCount, s.UintptrCount)
		// go through all leaves (findings) for this statement
		for _, leaf := range leaves {
			// print the findings tree node to stdout
			printRecursive(leaf, fset, 1)
			// then print the code position (filename and line) in a new line with indent 2
			printIndent(2)
			fmt.Println(fset.Position(leaf.Node.Pos()))
			// and print the code line itself as a context
			printIndent(2)
			fmt.Println(lines[fset.Position(leaf.Node.Pos()).Line-1])
		}
		fmt.Println()
	}
}

/**
 * saves the findings in the findings tree to disk
 */
func saveFindings(findingsTree *TreeNode, fset *token.FileSet, lines []string, pkg *base.PackageData) {
	// first, save all findings. go through all the leaves (findings) in the findings tree
	for _, finding := range findingsTree.collectLeaves() {
		// and write them to disk
		err := base.WriteAstFinding(base.AstFindingData{
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
		// if an error occurred, save it to the log
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "finding-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(finding.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}

	// then, save the function analysis data. Find all functions in the findings tree
	for function := range *findingsTree.findFunctions() {
		// for the function, identify the first and last lines of the function and join them together to get the code
		startLine := fset.Position(function.Node.Pos()).Line
		endLine := fset.Position(function.Node.End()).Line
		text := strings.Join(lines[startLine-1:endLine], "\n")
		// then write the function data to disk
		err := base.WriteAstFunction(base.AstFunctionData{
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
		// if an error occurred, save it to the log
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "function-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(function.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}

	// finally, save the statement analysis data. Find all statements in the findings tree
	for statement := range *findingsTree.findStatements() {
		// for the statement, identify the first and last lines of the statement and join them together to get the code
		startLine := fset.Position(statement.Node.Pos()).Line
		endLine := fset.Position(statement.Node.End()).Line
		text := strings.Join(lines[startLine-1:endLine], "\n")
		// then write the statement data to disk
		err := base.WriteAstStatement(base.AstStatementData{
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
		// if an error occurred, save it to the log
		if err != nil {
			_ = base.WriteErrorCondition(base.ErrorConditionData{
				Stage:             "statement-write",
				ProjectName:       pkg.ProjectName,
				PackageImportPath: pkg.ImportPath,
				FileName:          fset.Position(statement.Node.Pos()).Filename,
				Message:           err.Error(),
			})
		}
	}
}

/**
 * returns a string representation fitting my data set nomenclature depending on the match type of a given finding
 * tree node
 */
func matchTypeFor(n *TreeNode) string {
	// check which type the node has and return its string representation
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
