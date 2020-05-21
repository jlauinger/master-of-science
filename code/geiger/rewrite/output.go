package rewrite

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"strconv"
)

type IndentType int
const (
	Space IndentType = iota
	I
	T
	L
)

func printPkgTree(pkg *packages.Package, indents []IndentType, config Config, table *tablewriter.Table, seen *map[*packages.Package]bool) {
	(*seen)[pkg] = true

	countInThisPackage := getUnsafeCount(pkg, config)
	totalCount := getTotalCount(pkg, config, &map[*packages.Package]bool{})
	nameString := fmt.Sprintf("%s%s", getIndentString(indents), pkg.PkgPath)

	var colors []tablewriter.Colors
	if countInThisPackage > 0 {
		colors = []tablewriter.Colors{{tablewriter.FgRedColor}, {tablewriter.FgRedColor}, {tablewriter.FgRedColor}}
	} else if totalCount == 0 {
		colors = []tablewriter.Colors{{tablewriter.FgGreenColor}, {tablewriter.FgGreenColor}, {tablewriter.FgGreenColor}}
	} else {
		colors = []tablewriter.Colors{{tablewriter.Normal}, {tablewriter.Normal}, {tablewriter.Normal}}
	}

	table.Rich([]string{strconv.Itoa(countInThisPackage), strconv.Itoa(totalCount), nameString}, colors)

	childCount, _ := getImportsCount(pkg.Imports, config)

	nextIndents := getNextIndents(indents)

	if len(indents) == config.MaxIndent && childCount > 0 {
		table.Append([]string{"", "", fmt.Sprintf("%sMaximum depth reached. Use --level= to increase it",
			getIndentString(append(nextIndents, L)))})
		return
	}

	childIndex := 0
	for _, child := range pkg.Imports {
		if config.ShowStandardPackages == false && isStandardPackage(child) {
			continue
		}

		childIndex++
		isLast := childIndex == childCount

		_, ok := (*seen)[child]
		if config.ShortenSeenPackages && ok {
			table.Append([]string{"", "", fmt.Sprintf("%s%s...", getIndentString(nextIndents), child.PkgPath)})
			continue
		}

		if isLast {
			printPkgTree(child, append(nextIndents, L), config, table, seen)
		} else {
			printPkgTree(child, append(nextIndents, T), config, table, seen)
		}
	}
}

func getNextIndents(indents []IndentType) []IndentType {
	var nextIndents []IndentType
	if len(indents) > 0 {
		nextIndents = indents[0 : len(indents)-1]
		if indents[len(indents)-1] == L || indents[len(indents)-1] == Space {
			nextIndents = append(nextIndents, Space)
		} else {
			nextIndents = append(nextIndents, I)
		}
	} else {
		nextIndents = []IndentType{}
	}
	return nextIndents
}

func getIndentString(indents []IndentType) string {
	str := ""
	for _, indent := range indents {
		switch indent {
		case Space:
			str = fmt.Sprintf("%s%s", str, "  ")
		case I:
			str = fmt.Sprintf("%s%s", str, "│ ")
		case T:
			str = fmt.Sprintf("%s%s", str, "├─")
		case L:
			str = fmt.Sprintf("%s%s", str, "└─")
		}
	}
	return str
}

func getImportsCount(pkgs map[string]*packages.Package, config Config) (childCount, stdLibCount int) {
	for _, pkg := range pkgs {
		if isStandardPackage(pkg) {
			stdLibCount++
			if config.ShowStandardPackages {
				childCount++
			}
		} else {
			childCount++
		}
	}
	return
}