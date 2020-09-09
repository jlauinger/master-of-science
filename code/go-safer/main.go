package main

import (
	"github.com/stg-tud/thesis-2020-Lauinger-code/go-safer/passes/sliceheader"
	"github.com/stg-tud/thesis-2020-Lauinger-code/go-safer/passes/structcast"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	// invoke go vet main function with go-safer analyzers
	multichecker.Main(sliceheader.Analyzer, structcast.Analyzer)
}
