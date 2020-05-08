package unsafecount_test

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"linter/passes/unsafecount"
	"testing"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	testPackages := []string{
		"some_unsafe",
	}
	analysistest.Run(t, testdata, unsafecount.Analyzer, testPackages...)
}


