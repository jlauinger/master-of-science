package counter

import (
	"fmt"
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
}
