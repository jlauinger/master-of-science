package prettyprint

import (
	"fmt"
	"geiger/facts"
	"go/types"
)

func Print(results map[*types.Package]*facts.PackageInfo, maxIndent int, shortenSeen bool) {
	forest := buildForest(results)

	for _, tree := range forest {
		printTree(tree, 0, maxIndent, &map[*types.Package]bool{}, shortenSeen)
		fmt.Println("")
	}
}

func printTree(root *PackageTreeNode, indent, maxIndent int, seen *map[*types.Package]bool, shortenSeen bool) {
	if indent > maxIndent {
		return
	}
	printIndent(indent)

	totalCount := getTotalCount(root, &map[*types.Package]bool{})

	fmt.Printf("%s: %d (total %d)", root.Pkg.Path(), root.Info.ThisCount, totalCount)

	_, ok := (*seen)[root.Pkg]
	if ok && shortenSeen {
		if len(root.Children) > 0 {
			fmt.Printf("...\n")
		} else {
			fmt.Printf("\n")
		}
	} else {
		fmt.Printf("\n")
		(*seen)[root.Pkg] = true
		for _, child := range root.Children {
			printTree(child, indent + 1, maxIndent, seen, shortenSeen)
		}
	}
}

func printIndent (indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("    ")
	}
}

func getTotalCount(root *PackageTreeNode, seen *map[*types.Package]bool) int {
	_, ok := (*seen)[root.Pkg]
	if ok {
		return 0
	}
	(*seen)[root.Pkg] = true

	totalCount := root.Info.ThisCount

	for _, child := range root.Children {
		totalCount += getTotalCount(child, seen)
	}

	return totalCount
}