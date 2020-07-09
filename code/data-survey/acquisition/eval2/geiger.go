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

	initCache()

	for _, parsedPkg := range parsedPkgs {
		pkg := pkgsMap[parsedPkg.PkgPath]

		localUnsafeCounts := getUnsafeCounts(parsedPkg, pkg, fileToLineCountMap, fileToByteCountMap)
		localUnsafeSum := sumUp(localUnsafeCounts)

		pkg.UnsafeSum = localUnsafeSum

		pkg.UnsafePointerSum = localUnsafeCounts.UnsafePointerLocal
		pkg.UnsafePointerAssignment = localUnsafeCounts.UnsafePointerAssignment
		pkg.UnsafePointerCall = localUnsafeCounts.UnsafePointerCall
		pkg.UnsafePointerParameter = localUnsafeCounts.UnsafePointerParameter
		pkg.UnsafePointerVariable = localUnsafeCounts.UnsafePointerVariable
		pkg.UnsafePointerOther = localUnsafeCounts.UnsafePointerOther

		pkg.UnsafeSizeofSum = localUnsafeCounts.UnsafeSizeofLocal
		pkg.UnsafeSizeofAssignment = localUnsafeCounts.UnsafeSizeofAssignment
		pkg.UnsafeSizeofCall = localUnsafeCounts.UnsafeSizeofCall
		pkg.UnsafeSizeofParameter = localUnsafeCounts.UnsafeSizeofParameter
		pkg.UnsafeSizeofVariable = localUnsafeCounts.UnsafeSizeofVariable
		pkg.UnsafeSizeofOther = localUnsafeCounts.UnsafeSizeofOther

		pkg.UnsafeOffsetofSum = localUnsafeCounts.UnsafeOffsetofLocal
		pkg.UnsafeOffsetofAssignment = localUnsafeCounts.UnsafeOffsetofAssignment
		pkg.UnsafeOffsetofCall = localUnsafeCounts.UnsafeOffsetofCall
		pkg.UnsafeOffsetofParameter = localUnsafeCounts.UnsafeOffsetofParameter
		pkg.UnsafeOffsetofVariable = localUnsafeCounts.UnsafeOffsetofVariable
		pkg.UnsafeOffsetofOther = localUnsafeCounts.UnsafeOffsetofOther

		pkg.UnsafeAlignofSum = localUnsafeCounts.UnsafeAlignofLocal
		pkg.UnsafeAlignofAssignment = localUnsafeCounts.UnsafeAlignofAssignment
		pkg.UnsafeAlignofCall = localUnsafeCounts.UnsafeAlignofCall
		pkg.UnsafeAlignofParameter = localUnsafeCounts.UnsafeAlignofParameter
		pkg.UnsafeAlignofVariable = localUnsafeCounts.UnsafeAlignofVariable
		pkg.UnsafeAlignofOther = localUnsafeCounts.UnsafeAlignofOther

		pkg.SliceHeaderSum = localUnsafeCounts.SliceHeaderLocal
		pkg.SliceHeaderAssignment = localUnsafeCounts.SliceHeaderAssignment
		pkg.SliceHeaderCall = localUnsafeCounts.SliceHeaderCall
		pkg.SliceHeaderParameter = localUnsafeCounts.SliceHeaderParameter
		pkg.SliceHeaderVariable = localUnsafeCounts.SliceHeaderVariable
		pkg.SliceHeaderOther = localUnsafeCounts.SliceHeaderOther

		pkg.StringHeaderSum = localUnsafeCounts.StringHeaderLocal
		pkg.StringHeaderAssignment = localUnsafeCounts.StringHeaderAssignment
		pkg.StringHeaderCall = localUnsafeCounts.StringHeaderCall
		pkg.StringHeaderParameter = localUnsafeCounts.StringHeaderParameter
		pkg.StringHeaderVariable = localUnsafeCounts.StringHeaderVariable
		pkg.StringHeaderOther = localUnsafeCounts.StringHeaderOther

		pkg.UintptrSum = localUnsafeCounts.UintptrLocal
		pkg.UintptrAssignment = localUnsafeCounts.UintptrAssignment
		pkg.UintptrCall = localUnsafeCounts.UintptrCall
		pkg.UintptrParameter = localUnsafeCounts.UintptrParameter
		pkg.UintptrVariable = localUnsafeCounts.UintptrVariable
		pkg.UintptrOther = localUnsafeCounts.UintptrOther
	}
}

func getUnsafeCounts(parsedPkg *packages.Package, pkg *lexical.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) LocalPackageCounts {
	cachedCounts, ok := packageUnsafeCountsCache[parsedPkg.Name]
	if ok {
		return cachedCounts
	}

	inspectResult := inspector.New(parsedPkg.Syntax)
	localPackageCounts := LocalPackageCounts{}

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
			localPackageCounts.UnsafePointerLocal++
			matchType = "unsafe.Pointer"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.UnsafePointerAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.UnsafePointerCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.UnsafePointerParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.UnsafePointerVariable++
			} else {
				contextType = "other"
				localPackageCounts.UnsafePointerOther++
			}
		}
		if isUnsafeSizeof(node) {
			localPackageCounts.UnsafeSizeofLocal++
			matchType = "unsafe.Sizeof"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.UnsafeSizeofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.UnsafeSizeofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.UnsafeSizeofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.UnsafeSizeofVariable++
			} else {
				contextType = "other"
				localPackageCounts.UnsafeSizeofOther++
			}
		}
		if isUnsafeOffsetof(node) {
			localPackageCounts.UnsafeOffsetofLocal++
			matchType = "unsafe.Offsetof"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.UnsafeOffsetofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.UnsafeOffsetofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.UnsafeOffsetofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.UnsafeOffsetofVariable++
			} else {
				contextType = "other"
				localPackageCounts.UnsafeOffsetofOther++
			}
		}
		if isUnsafeAlignof(node) {
			localPackageCounts.UnsafeAlignofLocal++
			matchType = "unsafe.Alignof"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.UnsafeAlignofAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.UnsafeAlignofCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.UnsafeAlignofParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.UnsafeAlignofVariable++
			} else {
				contextType = "other"
				localPackageCounts.UnsafeAlignofOther++
			}
		}
		if isReflectSliceHeader(node) {
			localPackageCounts.SliceHeaderLocal++
			matchType = "reflect.SliceHeader"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.SliceHeaderAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.SliceHeaderCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.SliceHeaderParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.SliceHeaderVariable++
			} else {
				contextType = "other"
				localPackageCounts.SliceHeaderOther++
			}
		}
		if isReflectStringHeader(node) {
			localPackageCounts.StringHeaderLocal++
			matchType = "reflect.StringHeader"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.StringHeaderAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.StringHeaderCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.StringHeaderParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.StringHeaderVariable++
			} else {
				contextType = "other"
				localPackageCounts.StringHeaderOther++
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
			localPackageCounts.UintptrLocal++
			matchType = "uintptr"
			if isInAssignment(stack) {
				contextType = "assignment"
				localPackageCounts.UintptrAssignment++
			} else if isArgument(stack) {
				contextType = "call"
				localPackageCounts.UintptrCall++
			} else if isParameter(stack) {
				contextType = "parameter"
				localPackageCounts.UintptrParameter++
			} else if isInVariableDefinition(stack) {
				contextType = "variable"
				localPackageCounts.UintptrVariable++
			} else {
				contextType = "other"
				localPackageCounts.UintptrOther++
			}

			if matchType != "unknown" {
				writeData(n, parsedPkg, pkg, matchType, contextType, fileToLineCountMap, fileToByteCountMap)
			}
		}

		return true
	})

	packageUnsafeCountsCache[parsedPkg.Name] = localPackageCounts

	return localPackageCounts
}

func sumUp(counts LocalPackageCounts) int {
	return counts.UnsafePointerLocal + counts.UnsafeSizeofLocal + counts.UnsafeOffsetofLocal +
		counts.UnsafeAlignofLocal + counts.StringHeaderLocal + counts.SliceHeaderLocal + counts.UintptrLocal
}

func writeData(n ast.Node, parsedPkg *packages.Package, pkg *lexical.PackageData, matchType, contextType string,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	nodePosition := parsedPkg.Fset.Position(n.Pos())

	err := lexical.WriteGeigerFinding(lexical.GeigerFindingData{
		Text:              getCodeLine(parsedPkg, n),
		Context:           getCodeContext(parsedPkg, n),
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

func getTotalUnsafeCounts(parsedPkg *packages.Package, seen *map[*packages.Package]bool, pkg *lexical.PackageData, fileToLineCountMap, fileToByteCountMap map[string]int) TotalPackageCounts {
	_, ok := (*seen)[parsedPkg]
	if ok {
		return TotalPackageCounts{}
	}
	(*seen)[parsedPkg] = true

	unsafeCounts := getUnsafeCounts(parsedPkg, pkg, fileToLineCountMap, fileToByteCountMap)

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
}
