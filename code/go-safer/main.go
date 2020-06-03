package main

import (
	"github.com/stg-tud/thesis-2020-Lauinger-code/go-safer/passes/sliceheader"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(sliceheader.Analyzer)
}
