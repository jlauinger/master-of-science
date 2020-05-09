package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"linter/passes/literalheader"
)

func main() {
	multichecker.Main(literalheader.Analyzer,)
}
