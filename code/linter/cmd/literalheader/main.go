package main

import (
	"linter/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)


func main() {
	singlechecker.Main(analysis.LiteralHeaderAnalyzer)
}