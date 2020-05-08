package printunsafecount

import (
	"flag"
	"fmt"
	"golang.org/x/tools/go/analysis"
	"linter/passes/unsafecount"
)

var Analyzer = &analysis.Analyzer{
	Name:             "printunsafecount",
	Doc:              "pretty prints usages of unsafe Pointer",
	Flags:            flag.FlagSet{},
	Run:              run,
	RunDespiteErrors: true,
	Requires:         []*analysis.Analyzer{unsafecount.Analyzer},
	ResultType:       nil,
	FactTypes:        nil,
}

func run(pass *analysis.Pass) (interface{}, error) {
	fmt.Println(pass.Pkg.Name())

	return nil, nil
}