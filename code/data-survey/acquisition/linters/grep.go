package linters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"io"
	"os/exec"
	"strings"
)

/*
 * these are the strings to grep for
 */
var matchTypes = []string{"unsafe\\.Pointer", "unsafe\\.Sizeof", "unsafe\\.Alignof", "unsafe\\.Offsetof",
	"uintptr", "reflect\\.SliceHeader", "reflect\\.StringHeader"}

/**
 * callback handling the grep analysis coordination. This is called for each project
 */
func callbackGrep(project *base.ProjectData, packages []*base.PackageData, fileToPackageMap map[string]*base.PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	// run grep, log any potential error, then analyze the findings and write them to disk
	parsedGrepLines, err := grepForUnsafe(packages)
	if err != nil {
		_ = base.WriteErrorCondition(base.ErrorConditionData{
			Stage:             "grep-parse",
			ProjectName:       project.Name,
			PackageImportPath: "",
			FileName:          "",
			Message:           err.Error(),
		})
		fmt.Println("SAVING ERROR!")
	}
	analyzeGrepLines(parsedGrepLines, fileToPackageMap, fileToLineCountMap, fileToByteCountMap)
}

/**
 * runs ripgrep to find any unsafe usages in the given packages
 */
func grepForUnsafe(packages []*base.PackageData) ([]base.RipgrepOutputLine, error) {
	fmt.Println("  running ripgrep...")

	// initialize a list of files to grep through, then go through all packages and their files and append their
	// full effective paths to the list
	files := make([]string, 0, 1000)
	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			files = append(files, fullFilename)
		}
	}

	// build the ripgrep command using JSON output and including a context of 5 lines
	args := []string{strings.Join(matchTypes, "|"), "--context", "5", "--json"}
	args = append(args, files...)
	cmd := exec.Command("rg", args...)

	// run ripgrep and capture its output
	rgOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// initialize a JSON decoder and initialize a buffer for parsed lines
	dec := json.NewDecoder(bytes.NewReader(rgOutput))
	parsedLines := make([]base.RipgrepOutputLine, 0, 1000)

	// try decoding until the stream ends
	for {
		var message base.RipgrepOutputLine

		// decode the next ripgrep message if there any left, otherwise terminate
		err := dec.Decode(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// store the parsed line for later analysis
		parsedLines = append(parsedLines, message)
	}

	return parsedLines, nil
}

/**
 * analyzes the identified grep findings and writes them to disk
 */
func analyzeGrepLines(parsedLines []base.RipgrepOutputLine, fileToPackageMap map[string]*base.PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) {

	fmt.Println("  analyzing ripgrep output")

	// go through all grep finding lines
	for lineIdx, line := range parsedLines {
		// if the ripgrep message is not match, ignore it. There are also summary messages for example
		if line.MessageType != "match" {
			continue
		}

		// get the context (+/- 5 code lines) for this finding
		context := extractContext(line, lineIdx, parsedLines)

		// extract the filename where ripgrep found this finding
		fullFilename := line.Data.Path.Text

		// try to find the package containing that file, use a dummy if there should be none matching
		pkg, ok := fileToPackageMap[fullFilename]
		if !ok {
			pkg = &base.PackageData{
				ImportPath: "unknown-grep-error",
			}
		}

		// shorten the filename to not include the package directory when saving it to CSV
		filename := fullFilename[len(pkg.Dir)+1:]

		// go through all submatches, since a line could have several unsafe findings
		for _, subMatch := range line.Data.SubMatches {
			// build the finding data and write it to disk
			err := base.WriteGrepFinding(base.GrepFindingData{
				Text:              line.Data.Lines.Text,
				Context:           context,
				LineNumber:        line.Data.LineNumber,
				Column:            subMatch.Start,
				AbsoluteOffset:    line.Data.AbsoluteOffset,
				MatchType:         subMatch.Match.Text,
				FileName:          filename,
				FileLoc:           fileToByteCountMap[fullFilename],
				FileByteSize:      fileToLineCountMap[fullFilename],
				PackageImportPath: pkg.ImportPath,
				ModulePath:        pkg.ModulePath,
				ModuleVersion:     pkg.ModuleVersion,
				ProjectName:       pkg.ProjectName,
			})
			// and log any error
			if err != nil {
				_ = base.WriteErrorCondition(base.ErrorConditionData{
					Stage:             "ripgrep-parse",
					ProjectName:       pkg.ProjectName,
					PackageImportPath: pkg.ImportPath,
					FileName:          filename,
					Message:           err.Error(),
				})
				fmt.Println("SAVING ERROR!")
				continue
			}
		}
	}
}

/**
 * extracts the +/- 5 lines source code context for a given file
 */
func extractContext(line base.RipgrepOutputLine, lineIdx int, parsedLines []base.RipgrepOutputLine) string {
	// start with the matched line itself
	contextLines := []string{line.Data.Lines.Text}

	// add the context before that line
	for contextIdx := lineIdx - 1; contextIdx > base.Max(0, lineIdx-5); contextIdx-- {
		// extract the message line at that position
		contextLine := parsedLines[contextIdx]
		// only context or match messages can be appended for context, other messages denote boundaries and shall
		// abort the context generation
		if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
			contextLines = append([]string{contextLine.Data.Lines.Text}, contextLines...)
		} else {
			break
		}
	}

	// then add the context after that line
	for contextIdx := lineIdx + 1; contextIdx < base.Min(len(parsedLines), lineIdx+6); contextIdx++ {
		// extract the message line at that position
		contextLine := parsedLines[contextIdx]
		// only context or match messages can be appended for context, other messages denote boundaries and shall
		// abort the context generation
		if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
			contextLines = append(contextLines, contextLine.Data.Lines.Text)
		} else {
			break
		}
	}

	// finally merge the context lines to one string. The lines already terminate with \n so need to add it here
	return strings.Join(contextLines, "")
}