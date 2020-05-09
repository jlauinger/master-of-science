package unsafecountpass_test

import (
	"geiger/unsafecountpass"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	testPackages := []string{
		"some_unsafe",
	}
	analysistest.Run(t, testdata, unsafecountpass.Analyzer, testPackages...)
}


