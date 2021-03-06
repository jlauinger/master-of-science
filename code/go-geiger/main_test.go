package main

import (
	"bytes"
	"github.com/acarl005/stripansi"
	"github.com/stg-tud/thesis-2020-lauinger-code/go-geiger/counter"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestMin (t *testing.T) {
	got := counter.Min(5, 10)
	if got != 5 {
		t.Errorf("Min(5, 10) = %d; want 5", got)
	}
}

func TestPointerCount (t *testing.T) {
	output := bytes.NewBufferString("")

	counter.Run(counter.Config{
		MaxDepth:             9999,
		ShortenSeenPackages:  true,
		PrintLinkToPkgGoDev:  false,

		DetailedStats:        true,
		HideStats:            false,
		PrintUnsafeLines:     false,

		ShowStandardPackages: false,
		MatchFilter:          "pointer",
		ContextFilter:        "all",
		Output:               output,
	}, "./testdata")

	statsLine := stripansi.Strip(strings.Split(output.String(), "\n")[2])
	zp := regexp.MustCompile(` +`)
	stats := zp.Split(statsLine, -1)

	total, _ := strconv.Atoi(stats[1])
	local, _ := strconv.Atoi(stats[2])
	variable, _ := strconv.Atoi(stats[3])
	parameter, _ := strconv.Atoi(stats[4])
	assignment, _ := strconv.Atoi(stats[5])
	call, _ := strconv.Atoi(stats[6])
	other, _ := strconv.Atoi(stats[7])

	if local != 7 {
		t.Errorf("local = %d; want 7", local)
	}
	if total != 7 {
		t.Errorf("total = %d; want 7", total)
	}
	if variable != 2 {
		t.Errorf("variable = %d; want 2", variable)
	}
	if parameter != 1 {
		t.Errorf("parameter = %d; want 1", parameter)
	}
	if assignment != 3 {
		t.Errorf("assignment = %d; want 3", assignment)
	}
	if call != 1 {
		t.Errorf("call = %d; want 1", call)
	}
	if other != 0 {
		t.Errorf("other = %d; want 0", other)
	}
}

func TestAllCount (t *testing.T) {
	output := bytes.NewBufferString("")

	counter.Run(counter.Config{
		MaxDepth:             9999,
		ShortenSeenPackages:  true,
		PrintLinkToPkgGoDev:  false,

		DetailedStats:        true,
		HideStats:            false,
		PrintUnsafeLines:     false,

		ShowStandardPackages: false,
		MatchFilter:          "all",
		ContextFilter:        "all",
		Output:               output,
	}, "./testdata")

	statsLine := stripansi.Strip(strings.Split(output.String(), "\n")[2])
	zp := regexp.MustCompile(` +`)
	stats := zp.Split(statsLine, -1)

	total, _ := strconv.Atoi(stats[1])
	local, _ := strconv.Atoi(stats[2])
	variable, _ := strconv.Atoi(stats[3])
	parameter, _ := strconv.Atoi(stats[4])
	assignment, _ := strconv.Atoi(stats[5])
	call, _ := strconv.Atoi(stats[6])
	other, _ := strconv.Atoi(stats[7])

	if local != 13 {
		t.Errorf("local = %d; want 13", local)
	}
	if total != 13 {
		t.Errorf("total = %d; want 13", total)
	}
	if variable != 3 {
		t.Errorf("variable = %d; want 3", variable)
	}
	if parameter != 1 {
		t.Errorf("parameter = %d; want 1", parameter)
	}
	if assignment != 8 {
		t.Errorf("assignment = %d; want 8", assignment)
	}
	if call != 1 {
		t.Errorf("call = %d; want 1", call)
	}
	if other != 0 {
		t.Errorf("other = %d; want 0", other)
	}
}