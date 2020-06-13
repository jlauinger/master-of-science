package structcast_test

import (
	"github.com/stg-tud/thesis-2020-Lauinger-code/go-safer/passes/structcast"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	testPackages := []string{
		"bad/architecture_sized_variable",
	}
	analysistest.Run(t, testdata, structcast.Analyzer, testPackages...)
}

