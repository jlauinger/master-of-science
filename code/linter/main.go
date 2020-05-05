package main

import (
	"linter/passes/literalheader"
	"golang.org/x/tools/go/analysis/singlechecker"
)


func main() {
	singlechecker.Main(literalheader.Analyzer)
}
