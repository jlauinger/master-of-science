package prettyprint

import (
	"fmt"
	"geiger/facts"
	"go/types"
)

type Config struct {
	MaxIndent            int
	ShortenSeenPackages  bool
	ShowStandardPackages bool
}

func Print(results map[*types.Package]*facts.PackageInfo, config Config) {
	forest := buildForest(results)

	for _, tree := range forest {
		printTree(tree, 0, &map[*types.Package]bool{}, config)
		fmt.Println("")
	}

	if config.ShowStandardPackages == false {
		fmt.Println("Not showing packages from the standard library. Use GEIGER_SHOW_STD=true to include them.")
	}
}

func printTree(root *PackageTreeNode, indent int, seen *map[*types.Package]bool, config Config) {
	if config.ShowStandardPackages == false && isStandardPackage(root.Pkg.Path()) {
		return
	}

	if indent > config.MaxIndent {
		return
	}
	printIndent(indent)

	totalCount := getTotalCount(root, &map[*types.Package]bool{})

	fmt.Printf("%s: %d (total %d)", root.Pkg.Path(), root.Info.ThisCount, totalCount)

	_, ok := (*seen)[root.Pkg]
	if ok && config.ShortenSeenPackages {
		if len(root.Children) > 0 {
			fmt.Printf("...\n")
		} else {
			fmt.Printf("\n")
		}
	} else {
		fmt.Printf("\n")
		(*seen)[root.Pkg] = true
		for _, child := range root.Children {
			printTree(child, indent + 1, seen, config)
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