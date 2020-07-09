package eval2

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

func geigerPackages(project *lexical.ProjectData, pkgs []*lexical.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	fmt.Println("  parsing packages and counting unsafe using geiger...")

	pkgsMap := make(map[string]*lexical.PackageData)
	for _, pkg := range pkgs {
		pkgsMap[pkg.ImportPath] = pkg
	}

	paths := make([]string, 0)
	for _, pkg := range pkgs {
		paths = append(paths, pkg.ImportPath)
	}

	parsedPkgs, err := packages.Load(&packages.Config{
		Mode:       packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
					packages.NeedFiles | packages.NeedName | packages.NeedTypes,
		Tests:      false,
		Dir:		project.CheckoutPath,
	}, paths...)
	if err != nil {
		panic("error loading packages")
	}

	// count each package on its own
	for _, parsedPkg := range parsedPkgs {
		pkg := pkgsMap[parsedPkg.PkgPath]
		geigerSinglePackage(parsedPkg, pkg, fileToLineCountMap, fileToByteCountMap)
	}

	// sum up the dependencies using the package hash map and a cache
	seen := make(map[string]bool)
	for _, pkg := range pkgs {
		sumUpDependencies(pkg, pkgsMap, seen)
	}
}

func geigerSinglePackage(parsedPkg *packages.Package, pkg *lexical.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) {
	inspectResult := inspector.New(parsedPkg.Syntax)

	seenSelectorExprs := map[*ast.SelectorExpr]bool{}
	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		_, ok := seenSelectorExprs[node]
		if ok {
			return true
		}
		seenSelectorExprs[node] = true

		matchType := "unknown"
		contextType := "unknown"

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

		if matchType != "unknown" {
			writeData(n, parsedPkg, pkg, matchType, contextType, fileToLineCountMap, fileToByteCountMap)
		}

		return true
	})

	seenIdents := map[*ast.Ident]bool{}
	inspectResult.WithStack([]ast.Node{(*ast.Ident)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.Ident)
		_, ok := seenIdents[node]
		if ok {
			return true
		}
		seenIdents[node] = true

		matchType := "unknown"
		contextType := "unknown"

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

		if matchType != "unknown" {
			writeData(n, parsedPkg, pkg, matchType, contextType, fileToLineCountMap, fileToByteCountMap)
		}

		return true
	})

	pkg.UnsafeSum = pkg.UnsafePointerSum + pkg.UnsafeSizeofSum + pkg.UnsafeOffsetofSum +
		pkg.UnsafeAlignofSum + pkg.StringHeaderSum + pkg.SliceHeaderSum + pkg.UintptrSum
}

func writeData(n ast.Node, parsedPkg *packages.Package, pkg *lexical.PackageData, matchType, contextType string,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	nodePosition := parsedPkg.Fset.File(n.Pos()).Position(n.Pos())

	text, context := getCodeContext(parsedPkg, n)

	err := lexical.WriteGeigerFinding(lexical.GeigerFindingData{
		Text:              text,
		Context:           context,
		LineNumber:        nodePosition.Line,
		Column:            nodePosition.Column,
		AbsoluteOffset:    nodePosition.Offset,
		MatchType:         matchType,
		ContextType:       contextType,
		FileName:          nodePosition.Filename,
		FileLoc:           fileToLineCountMap[nodePosition.Filename],
		FileByteSize:      fileToByteCountMap[nodePosition.Filename],
		PackageImportPath: pkg.ImportPath,
		ModulePath:        pkg.ModulePath,
		ModuleVersion:     pkg.ModuleVersion,
		ProjectName:       pkg.ProjectName,
	})

	if err != nil {
		_ = lexical.WriteErrorCondition(lexical.ErrorConditionData{
			Stage:             "geiger-save",
			ProjectName:       pkg.ProjectName,
			PackageImportPath: pkg.ImportPath,
			FileName:          parsedPkg.Fset.File(n.Pos()).Name(),
			Message:           err.Error(),
		})
		fmt.Println("SAVING ERROR!")
	}
}

func sumUpDependencies(pkg *lexical.PackageData, pkgsMap map[string]*lexical.PackageData, seen map[string]bool) {
	_, ok := seen[pkg.ImportPath]
	if ok {
		return
	}

	for _, childPath := range pkg.Imports {
		child, ok := pkgsMap[childPath]
		if !ok {
			panic("child not found")
		}
		sumUpDependencies(child, pkgsMap, seen)

		pkg.UnsafeSumWithDependencies += child.UnsafeSumWithDependencies
	}

	seen[pkg.ImportPath] = true
}

/*func getTotalUnsafeCounts(parsedPkg *packages.Package, seen *map[*packages.Package]bool, pkg *lexical.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) TotalPackageCounts {
	_, ok := (*seen)[parsedPkg]
	if ok {
		return TotalPackageCounts{}
	}
	(*seen)[parsedPkg] = true

	unsafeCounts := geigerSinglePackage(parsedPkg, pkg, fileToLineCountMap, fileToByteCountMap)

	totalCounts := TotalPackageCounts{
		UnsafePointerTotal:  unsafeCounts.UnsafeAlignofLocal,
		UnsafeSizeofTotal:   unsafeCounts.UnsafeSizeofLocal,
		UnsafeOffsetofTotal: unsafeCounts.UnsafeOffsetofLocal,
		UnsafeAlignofTotal:  unsafeCounts.UnsafeAlignofLocal,
		SliceHeaderTotal:    unsafeCounts.SliceHeaderLocal,
		StringHeaderTotal:   unsafeCounts.StringHeaderLocal,
		UintptrTotal:        unsafeCounts.UintptrLocal,
	}

	for _, child := range parsedPkg.Imports {
		totalCountsChild := getTotalUnsafeCounts(child, seen, pkg, fileToLineCountMap, fileToByteCountMap)

		totalCounts.UnsafePointerTotal += totalCountsChild.UnsafePointerTotal
		totalCounts.UnsafeSizeofTotal += totalCountsChild.UnsafeSizeofTotal
		totalCounts.UnsafeOffsetofTotal += totalCountsChild.UnsafeOffsetofTotal
		totalCounts.UnsafeAlignofTotal += totalCountsChild.UnsafeAlignofTotal
		totalCounts.SliceHeaderTotal += totalCountsChild.SliceHeaderTotal
		totalCounts.StringHeaderTotal += totalCountsChild.StringHeaderTotal
		totalCounts.UintptrTotal += totalCountsChild.UintptrTotal
	}

	return totalCounts
}*/
