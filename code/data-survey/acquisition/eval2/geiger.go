package eval2

import (
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

func geigerPackages(project *lexical.ProjectData, pkgs []*lexical.PackageData) {
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
		localUnsafeCounts := getUnsafeCounts(parsedPkg)
		localUnsafeSum := sumUp(localUnsafeCounts)

		pkgsMap[parsedPkg.PkgPath].UnsafeSum = localUnsafeSum

		if localUnsafeSum > 0 {
			fmt.Printf("%s: total %d\n", parsedPkg.PkgPath, localUnsafeSum)
		}
	}
}

func getUnsafeCounts(pkg *packages.Package) LocalPackageCounts {
	cachedCounts, ok := packageUnsafeCountsCache[pkg]
	if ok {
		return cachedCounts
	}

	inspectResult := inspector.New(pkg.Syntax)
	localPackageCounts := LocalPackageCounts{}

	seenSelectorExprs := map[*ast.SelectorExpr]bool{}
	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		_, ok := seenSelectorExprs[node]
		if ok {
			return true
		}
		seenSelectorExprs[node] = true

		if isUnsafePointer(node) {
			localPackageCounts.UnsafePointerLocal++
			if isInAssignment(stack) {
				localPackageCounts.UnsafePointerAssignment++
			} else if isArgument(stack) {
				localPackageCounts.UnsafePointerCall++
			} else if isParameter(stack) {
				localPackageCounts.UnsafePointerParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.UnsafePointerVariable++
			} else {
				localPackageCounts.UnsafePointerOther++
			}
		}
		if isUnsafeSizeof(node) {
			localPackageCounts.UnsafeSizeofLocal++
			if isInAssignment(stack) {
				localPackageCounts.UnsafeSizeofAssignment++
			} else if isArgument(stack) {
				localPackageCounts.UnsafeSizeofCall++
			} else if isParameter(stack) {
				localPackageCounts.UnsafeSizeofParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.UnsafeSizeofVariable++
			} else {
				localPackageCounts.UnsafeSizeofOther++
			}
		}
		if isUnsafeOffsetof(node) {
			localPackageCounts.UnsafeOffsetofLocal++
			if isInAssignment(stack) {
				localPackageCounts.UnsafeOffsetofAssignment++
			} else if isArgument(stack) {
				localPackageCounts.UnsafeOffsetofCall++
			} else if isParameter(stack) {
				localPackageCounts.UnsafeOffsetofParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.UnsafeOffsetofVariable++
			} else {
				localPackageCounts.UnsafeOffsetofOther++
			}
		}
		if isUnsafeAlignof(node) {
			localPackageCounts.UnsafeAlignofLocal++
			if isInAssignment(stack) {
				localPackageCounts.UnsafeAlignofAssignment++
			} else if isArgument(stack) {
				localPackageCounts.UnsafeAlignofCall++
			} else if isParameter(stack) {
				localPackageCounts.UnsafeAlignofParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.UnsafeAlignofVariable++
			} else {
				localPackageCounts.UnsafeAlignofOther++
			}
		}
		if isReflectSliceHeader(node) {
			localPackageCounts.SliceHeaderLocal++
			if isInAssignment(stack) {
				localPackageCounts.SliceHeaderAssignment++
			} else if isArgument(stack) {
				localPackageCounts.SliceHeaderCall++
			} else if isParameter(stack) {
				localPackageCounts.SliceHeaderParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.SliceHeaderVariable++
			} else {
				localPackageCounts.SliceHeaderOther++
			}
		}
		if isReflectStringHeader(node) {
			localPackageCounts.StringHeaderLocal++
			if isInAssignment(stack) {
				localPackageCounts.StringHeaderAssignment++
			} else if isArgument(stack) {
				localPackageCounts.StringHeaderCall++
			} else if isParameter(stack) {
				localPackageCounts.StringHeaderParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.StringHeaderVariable++
			} else {
				localPackageCounts.StringHeaderOther++
			}
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

		if isUintptr(node) {
			localPackageCounts.UintptrLocal++
			if isInAssignment(stack) {
				localPackageCounts.UintptrAssignment++
			} else if isArgument(stack) {
				localPackageCounts.UintptrCall++
			} else if isParameter(stack) {
				localPackageCounts.UintptrParameter++
			} else if isInVariableDefinition(stack) {
				localPackageCounts.UintptrVariable++
			} else {
				localPackageCounts.UintptrOther++
			}
		}

		return true
	})

	packageUnsafeCountsCache[pkg] = localPackageCounts

	return localPackageCounts
}

func sumUp(counts LocalPackageCounts) int {
	return counts.UnsafePointerLocal + counts.UnsafeSizeofLocal + counts.UnsafeOffsetofLocal +
		counts.UnsafeAlignofLocal + counts.StringHeaderLocal + counts.SliceHeaderLocal + counts.UintptrLocal
}

func getTotalUnsafeCounts(pkg *packages.Package, seen *map[*packages.Package]bool) TotalPackageCounts {
	_, ok := (*seen)[pkg]
	if ok {
		return TotalPackageCounts{}
	}
	(*seen)[pkg] = true

	unsafeCounts := getUnsafeCounts(pkg)

	totalCounts := TotalPackageCounts{
		UnsafePointerTotal:  unsafeCounts.UnsafeAlignofLocal,
		UnsafeSizeofTotal:   unsafeCounts.UnsafeSizeofLocal,
		UnsafeOffsetofTotal: unsafeCounts.UnsafeOffsetofLocal,
		UnsafeAlignofTotal:  unsafeCounts.UnsafeAlignofLocal,
		SliceHeaderTotal:    unsafeCounts.SliceHeaderLocal,
		StringHeaderTotal:   unsafeCounts.StringHeaderLocal,
		UintptrTotal:        unsafeCounts.UintptrLocal,
	}

	for _, child := range pkg.Imports {
		totalCountsChild := getTotalUnsafeCounts(child, seen)

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
