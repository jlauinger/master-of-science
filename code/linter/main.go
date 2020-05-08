package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"linter/passes/literalheader"
	"linter/passes/unsafecount"
)

func main() {
	unitchecker.Main(literalheader.Analyzer, unsafecount.Analyzer,)
}
