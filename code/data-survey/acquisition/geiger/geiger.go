package geiger

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
	"strings"
)

/**
 * runs the go-geiger logic on the given packages
 */
func geigerPackages(project *base.ProjectData, pkgs []*base.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	fmt.Println("  parsing packages and counting unsafe using geiger...")

	// initialize a hash map from package import paths to the package data structure and fill it with all packages
	pkgsMap := make(map[string]*base.PackageData)
	for _, pkg := range pkgs {
		pkgsMap[pkg.ImportPath] = pkg
	}

	// initialize and fill a list of all package import paths
	paths := make([]string, 0)
	for _, pkg := range pkgs {
		paths = append(paths, pkg.ImportPath)
	}

	// parse all packages by their import paths with full parsing
	parsedPkgs, err := packages.Load(&packages.Config{
		Mode:       packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
					packages.NeedFiles | packages.NeedName | packages.NeedTypes,
		Tests:      false,
		Dir:		project.CheckoutPath,
	}, paths...)
	if err != nil {
		panic("error loading packages")
	}

	// go through all the parsed packages
	for _, parsedPkg := range parsedPkgs {
		// identify the acquisition tool package structure for this parsed package
		pkg := pkgsMap[parsedPkg.PkgPath]
		// and go-geiger it
		geigerSinglePackage(parsedPkg, pkg, fileToLineCountMap, fileToByteCountMap)
	}
}

/**
 * analyzes a single package using go-geiger logic
 */
func geigerSinglePackage(parsedPkg *packages.Package, pkg *base.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	// initialize an AST inspector to make walking the tree easier
	inspectResult := inspector.New(parsedPkg.Syntax)

	// initialize a hash set of expressions that were already seen to avoid double-counting with the Stack
	// inspector function
	seenSelectorExprs := map[*ast.SelectorExpr]bool{}
	// go through all selector expression nodes
	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		// check if the node has already been analyzed and skip it if so. Otherwise now mark it seen
		_, ok := seenSelectorExprs[node]
		if ok {
			return true
		}
		seenSelectorExprs[node] = true

		// initialize match and context types with their placeholders
		matchType := "unknown"
		contextType := "unknown"

		// one by another, check the different match and context types, set the variables accordingly and increment
		// the corresponding package unsafe counts
		if isUnsafePointer(node) {
			pkg.UnsafePointerSum++
			matchType = "unsafe.Pointer"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.UnsafePointerAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.UnsafePointerCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.UnsafePointerParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.UnsafePointerVariable++
			} else {
				contextType = "other"
				pkg.UnsafePointerOther++
			}
		}
		if isUnsafeSizeof(node) {
			pkg.UnsafeSizeofSum++
			matchType = "unsafe.Sizeof"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.UnsafeSizeofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.UnsafeSizeofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.UnsafeSizeofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.UnsafeSizeofVariable++
			} else {
				contextType = "other"
				pkg.UnsafeSizeofOther++
			}
		}
		if isUnsafeOffsetof(node) {
			pkg.UnsafeOffsetofSum++
			matchType = "unsafe.Offsetof"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.UnsafeOffsetofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.UnsafeOffsetofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.UnsafeOffsetofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.UnsafeOffsetofVariable++
			} else {
				contextType = "other"
				pkg.UnsafeOffsetofOther++
			}
		}
		if isUnsafeAlignof(node) {
			pkg.UnsafeAlignofSum++
			matchType = "unsafe.Alignof"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.UnsafeAlignofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.UnsafeAlignofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.UnsafeAlignofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.UnsafeAlignofVariable++
			} else {
				contextType = "other"
				pkg.UnsafeAlignofOther++
			}
		}
		if isReflectSliceHeader(node) {
			pkg.SliceHeaderSum++
			matchType = "reflect.SliceHeader"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.SliceHeaderAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.SliceHeaderCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.SliceHeaderParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.SliceHeaderVariable++
			} else {
				contextType = "other"
				pkg.SliceHeaderOther++
			}
		}
		if isReflectStringHeader(node) {
			pkg.StringHeaderSum++
			matchType = "reflect.StringHeader"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.StringHeaderAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.StringHeaderCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.StringHeaderParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.StringHeaderVariable++
			} else {
				contextType = "other"
				pkg.StringHeaderOther++
			}
		}

		// if there was a match, save it to CSV
		if matchType != "unknown" {
			writeData(n, parsedPkg, pkg, matchType, contextType, fileToLineCountMap, fileToByteCountMap)
		}

		// return true to continue walking the AST
		return true
	})

	// similar to the selector expressions, initialize a hash set of identifiers that were already seen to avoid
	// double-counting with the Stack inspector function
	seenIdents := map[*ast.Ident]bool{}
	// go through all identifier expression nodes
	inspectResult.WithStack([]ast.Node{(*ast.Ident)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.Ident)
		// check if the node has already been analyzed and skip it if so. Otherwise now mark it seen
		_, ok := seenIdents[node]
		if ok {
			return true
		}
		seenIdents[node] = true

		// initialize match and context types with their placeholders
		matchType := "unknown"
		contextType := "unknown"

		// one by another, check the different match and context types, set the variables accordingly and increment
		// the corresponding package unsafe counts
		if isUintptr(node) {
			pkg.UintptrSum++
			matchType = "uintptr"
			if isInAssignment(stack) {
				contextType = "assignment"
				pkg.UintptrAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				pkg.UintptrCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				pkg.UintptrParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				pkg.UintptrVariable++
			} else {
				contextType = "other"
				pkg.UintptrOther++
			}
		}

		// if there was a match, save it to CSV
		if matchType != "unknown" {
			writeData(n, parsedPkg, pkg, matchType, contextType, fileToLineCountMap, fileToByteCountMap)
		}

		// return true to continue walking the AST
		return true
	})

	pkg.UnsafeSum = pkg.UnsafePointerSum + pkg.UnsafeSizeofSum + pkg.UnsafeOffsetofSum +
		pkg.UnsafeAlignofSum + pkg.StringHeaderSum + pkg.SliceHeaderSum + pkg.UintptrSum
}

/**
 * writes a go-geiger finding to disk
 */
func writeData(n ast.Node, parsedPkg *packages.Package, pkg *base.PackageData, matchType, contextType string,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	// identify the node position from the parsing file set and get the 1 and +/- 5 lines contexts
	nodePosition := parsedPkg.Fset.File(n.Pos()).Position(n.Pos())
	text, context := getCodeContext(parsedPkg, n)

	// build the filename that will be saved into the CSV file
	var filename string
	if strings.Contains(nodePosition.Filename, ".cache/go-build") {
		// if it's a file from the build cache (e.g. Cgo file), use the full file path
		filename = nodePosition.Filename
	} else {
		// otherwise store it without the package directory to save space
		filename = nodePosition.Filename[len(pkg.Dir)+1:]
	}

	// write the finding data to CSV
	err := base.WriteGeigerFinding(base.GeigerFindingData{
		Text:              text,
		Context:           context,
		LineNumber:        nodePosition.Line,
		Column:            nodePosition.Column,
		AbsoluteOffset:    nodePosition.Offset,
		MatchType:         matchType,
		ContextType:       contextType,
		FileName:          filename,
		FileLoc:           fileToLineCountMap[nodePosition.Filename],
		FileByteSize:      fileToByteCountMap[nodePosition.Filename],
		PackageImportPath: pkg.ImportPath,
		PackageDir:        pkg.Dir,
		ModulePath:        pkg.ModulePath,
		ModuleVersion:     pkg.ModuleVersion,
		ProjectName:       pkg.ProjectName,
	})

	// and log any potential error
	if err != nil {
		_ = base.WriteErrorCondition(base.ErrorConditionData{
			Stage:             "geiger-save",
			ProjectName:       pkg.ProjectName,
			PackageImportPath: pkg.ImportPath,
			FileName:          parsedPkg.Fset.File(n.Pos()).Name(),
			Message:           err.Error(),
		})
		fmt.Println("SAVING ERROR!")
	}
}
