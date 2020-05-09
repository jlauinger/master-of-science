package prettyprint

import (
	"geiger/facts"
	"go/types"
)

type PackageTreeNode struct {
	Pkg *types.Package
	Info *facts.PackageInfo
	Children []*PackageTreeNode
}

func buildForest(results map[*types.Package]*facts.PackageInfo) []*PackageTreeNode {
	roots := findRoots(results)

	trees := make([]*PackageTreeNode, len(roots))

	for i, root := range roots {
		trees[i] = buildTree(root, results)
	}

	return trees
}

func findRoots(results map[*types.Package]*facts.PackageInfo) []*types.Package {
	possibleRoots := make(map[*types.Package]bool, len(results))
	for pkg := range results {
		possibleRoots[pkg] = true
	}

	for pkg := range results {
		for _, childPkg := range pkg.Imports() {
			possibleRoots[childPkg] = false
		}
	}

	roots := make([]*types.Package, 0)
	for pkg, possibleRoot := range possibleRoots {
		if possibleRoot {
			roots = append(roots, pkg)
		}
	}
	
	return roots
}

func buildTree(root *types.Package, results map[*types.Package]*facts.PackageInfo) *PackageTreeNode {
	children := make([]*PackageTreeNode, len(root.Imports()))

	for i, child := range root.Imports() {
		children[i] = buildTree(child, results)
	}

	return &PackageTreeNode{
		Pkg:      root,
		Info:     results[root],
		Children: children,
	}
}