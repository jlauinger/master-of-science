package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"linter/passes/literalheader"
	"linter/passes/unsafecount"
)

func main() {
	multichecker.Main(literalheader.Analyzer, unsafecount.Analyzer,)
}
