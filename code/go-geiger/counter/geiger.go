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

	for _, pkg := range pkgs {
		initCache()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Local", "Total", "Package"})
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetAutoWrapText(false)
		table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

		stats := printPkgTree(pkg, []IndentType{}, config, table, &map[*packages.Package]bool{})

		if config.PrintUnsafeLines {
			fmt.Println()
		}

		table.Render()

		printStats(pkg, stats)
	}

	printLegend()
}

func printLegend() {
	fmt.Printf("%s have no unsafe.Pointer usages\n", color.GreenString("Packages in green"))
	fmt.Printf("%s contain unsafe.Pointer usages\n", color.RedString("Packages in red"))
	fmt.Printf("%s import packages with unsafe.Pointer usages\n", color.WhiteString("Packages in white"))
}

func printStats(pkg *packages.Package, stats Stats) {
	fmt.Println()

	fmt.Printf("Package %s including imports effectively makes up %d packages\n", pkg.PkgPath, stats.ImportCount + 1)

	if stats.UnsafeCount > 0 {
		color.Red("  %d of those contain unsafe.Pointer usages\n", stats.UnsafeCount)
	}
	if stats.TransitivelyUnsafeCount > 0 {
		color.White("  %d of those further import packages that contain unsafe.Pointer usages\n",
			stats.TransitivelyUnsafeCount)
	}
	if stats.SafeCount > 0 {
		color.Green("  %d of those do not contain any unsafe.Pointer usages\n", stats.SafeCount)
	}

	fmt.Println()
}
