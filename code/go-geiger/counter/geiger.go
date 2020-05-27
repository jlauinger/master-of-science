package counter

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"os"
)

type Config struct {
	MaxIndent            int
	ShortenSeenPackages  bool
	ShowStandardPackages bool
	PrintLinkToPkgGoDev  bool
	PrintUnsafeLines     bool
}

func Run(config Config, paths... string) {
	mode := packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedFiles | packages.NeedName

	if config.PrintUnsafeLines {
		mode |= packages.NeedTypes
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:       mode,
		Tests:      false,
	}, paths...)

	if err != nil {
		panic(err)
	}

	initCache()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Local", "Total", "Package"})
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, pkg := range pkgs {
		printPkgTree(pkg, []IndentType{}, config, table, &map[*packages.Package]bool{})
	}

	if config.PrintUnsafeLines {
		fmt.Println()
	}

	table.Render()

	printLegend()
}

func printLegend() {
	fmt.Println()

	fmt.Printf("%s have no unsafe.Pointer usages\n", color.GreenString("Packages in green"))
	fmt.Printf("%s contain unsafe.Pointer usages\n", color.RedString("Packages in red"))
	fmt.Printf("%s import packages with unsafe.Pointer usages\n", color.WhiteString("Packages in white"))
}
