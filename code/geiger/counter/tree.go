package counter

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"sort"
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
	totalCount := getTotalUnsafeCount(pkg, config, &map[*packages.Package]bool{})
	nameString := fmt.Sprintf("%s%s", getIndentString(indents), pkg.PkgPath)

	colors := getColors(countInThisPackage, totalCount)

	table.Rich([]string{strconv.Itoa(countInThisPackage), strconv.Itoa(totalCount), nameString}, colors)

	childCount, _ := getImportsCount(pkg.Imports, config)
	nextIndents := getNextIndents(indents)

	if len(indents) == config.MaxIndent && childCount > 0 {
		table.Append([]string{"", "", fmt.Sprintf("%sMaximum depth reached. Use --level= to increase it",
			getIndentString(append(nextIndents, L)))})
		return
	}

	childKeys := make([]string, 0, len(pkg.Imports))
	for childKey := range pkg.Imports {
		childKeys = append(childKeys, childKey)
	}
	sort.Strings(childKeys)

	childIndex := 0
	for _, childKey := range childKeys {
		child := pkg.Imports[childKey]

		if config.ShowStandardPackages == false && isStandardPackage(child) {
			continue
		}

		childIndex++
		childIndents := getChildIndents(childIndex, childCount, nextIndents)

		_, ok := (*seen)[child]
		if config.ShortenSeenPackages && ok {
			countInChild := getUnsafeCount(child, config)
			totalCountInChild := getTotalUnsafeCount(child, config, &map[*packages.Package]bool{})
			table.Rich([]string{strconv.Itoa(countInChild), strconv.Itoa(totalCountInChild),
				fmt.Sprintf("%s%s...", getIndentString(childIndents), child.PkgPath)},
				getColors(0, totalCountInChild))
			continue
		}

		printPkgTree(child, childIndents, config, table, seen)
	}
}

func getChildIndents(childIndex int, childCount int, nextIndents []IndentType) []IndentType {
	isLast := childIndex == childCount

	var nextChildIndents []IndentType
	if isLast {
		nextChildIndents = append(nextIndents, L)
	} else {
		nextChildIndents = append(nextIndents, T)
	}
	return nextChildIndents
}

func getColors(countInThisPackage int, totalCount int) []tablewriter.Colors {
	var colors []tablewriter.Colors
	if countInThisPackage > 0 {
		colors = []tablewriter.Colors{{tablewriter.FgRedColor}, {tablewriter.FgRedColor}, {tablewriter.FgRedColor}}
	} else if totalCount == 0 {
		colors = []tablewriter.Colors{{tablewriter.FgGreenColor}, {tablewriter.FgGreenColor}, {tablewriter.FgGreenColor}}
	} else {
		colors = []tablewriter.Colors{{tablewriter.Normal}, {tablewriter.Normal}, {tablewriter.Normal}}
	}
	return colors
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