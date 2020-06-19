package counter

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"io"
)

type Config struct {
	MaxIndent            int
	ShortenSeenPackages  bool
	ShowStandardPackages bool
	PrintLinkToPkgGoDev  bool
	PrintUnsafeLines     bool
	DetailedStats        bool
	Output               io.Writer
	Filter               string
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

		table := tablewriter.NewWriter(config.Output)
		if config.DetailedStats {
			table.SetHeader([]string{"Local", "Total", "Variable", "Parameter", "Assignment", "Call", "Other", "Package"})
			table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER,
				tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER,
				tablewriter.ALIGN_LEFT})
		} else {
			table.SetHeader([]string{"Local", "Total", "Package"})
			table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
		}
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetAutoWrapText(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

		stats := printPkgTree(pkg, []IndentType{}, config, table, &map[*packages.Package]bool{})

		if config.PrintUnsafeLines {
			_, _ = fmt.Fprintln(config.Output)
		}

		table.Render()

		printStats(pkg, stats, config)
	}

	printLegend(config)
}

func printLegend(config Config) {
	_, _ = fmt.Fprintf(config.Output, "%s have no unsafe.Pointer usages\n", color.GreenString("Packages in green"))
	_, _ = fmt.Fprintf(config.Output, "%s contain unsafe.Pointer usages\n", color.RedString("Packages in red"))
	_, _ = fmt.Fprintf(config.Output, "%s import packages with unsafe.Pointer usages\n", color.WhiteString("Packages in white"))
}

func printStats(pkg *packages.Package, stats Stats, config Config) {
	_, _ = fmt.Fprintln(config.Output)

	_, _ = fmt.Fprintf(config.Output, "Package %s including imports effectively makes up %d packages\n", pkg.PkgPath, stats.ImportCount+1)

	if stats.UnsafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.RedString("  %d of those contain unsafe.Pointer usages\n", stats.UnsafeCount))
	}
	if stats.TransitivelyUnsafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.WhiteString("  %d of those further import packages that contain unsafe.Pointer usages\n",
			stats.TransitivelyUnsafeCount))
	}
	if stats.SafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.GreenString("  %d of those do not contain any unsafe.Pointer usages\n", stats.SafeCount))
	}

	_, _ = fmt.Fprintln(config.Output)
}
