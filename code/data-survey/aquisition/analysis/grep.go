package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var matchTypes = []string{"unsafe\\.Pointer", "unsafe\\.Sizeof", "unsafe\\.Alignof", "unsafe\\.Offsetof",
	"uintptr", "reflect\\.SliceHeader", "reflect\\.StringHeader"}


func grepForUnsafe(packages []*PackageData) ([]RipgrepOutputLine, error) {
	files := make([]string, 0, 1000)

	fmt.Println("  running ripgrep...")

	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			fullFilename := fmt.Sprintf("%s/%s", pkg.Dir, file)
			files = append(files, fullFilename)
		}
	}

	args := []string{strings.Join(matchTypes, "|"), "--context", "5", "--json"}
	args = append(args, files...)

	cmd := exec.Command("rg", args...)

	rgOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(rgOutput))
	parsedLines := make([]RipgrepOutputLine, 0, 1000)

	for {
		var message RipgrepOutputLine

		err := dec.Decode(&message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		parsedLines = append(parsedLines, message)
	}

	return parsedLines, nil
}

func analyzeGrepLines(parsedLines []RipgrepOutputLine, fileToPackageMap map[string]*PackageData,
	fileToLineCountMap, fileToByteCountMap map[string]int) map[string]string {

	fmt.Println("  analyzing ripgrep output")

	var filesToCopy = make(map[string]string, 500)

	for lineIdx, line := range parsedLines {
		if line.MessageType == "match" {
			contextLines := []string{line.Data.Lines.Text}

			// context before line
			for contextIdx := lineIdx - 1; contextIdx > Max(0, lineIdx - 5); contextIdx-- {
				contextLine := parsedLines[contextIdx]
				if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
					contextLines = append([]string{contextLine.Data.Lines.Text}, contextLines...)
				} else {
					break
				}
			}

			// context after line
			for contextIdx := lineIdx + 1; contextIdx < Min(len(parsedLines), lineIdx + 6); contextIdx++ {
				contextLine := parsedLines[contextIdx]
				if contextLine.MessageType == "context" || contextLine.MessageType == "match" {
					contextLines = append(contextLines, contextLine.Data.Lines.Text)
				} else {
					break
				}
			}

			context := strings.Join(contextLines, "")

			fullFilename := line.Data.Path.Text

			pkg, ok := fileToPackageMap[fullFilename]
			if !ok {
				pkg = &PackageData{
					ImportPath: "unknown-vet-error",
				}
			}

			filename := fullFilename[len(pkg.Dir)+1:]

			for _, subMatch := range line.Data.SubMatches {
				copyDestination := fmt.Sprintf("%s/%s", pkg.ImportPath, filename)

				err := WriteGrepFinding(GrepFindingData{
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
					FileCopyPath:      copyDestination,
				})

				filesToCopy[fullFilename] = copyDestination

				if err != nil {
					_ = WriteErrorCondition(ErrorConditionData{
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

	return filesToCopy
}