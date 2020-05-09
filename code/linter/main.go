package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"linter/passes/literalheader"
)

func main() {
	singlechecker.Main(literalheader.Analyzer)
}
