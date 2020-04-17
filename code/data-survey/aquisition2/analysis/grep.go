package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func grepForUnsafe(files []string) ([]RipgrepOutputLine, error) {
	args := []string{"unsafe.Pointer", "--context", "5", "--json"}
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

func analyzeGrepLines(parsedLines []RipgrepOutputLine, fileToModuleMap map[string]ModuleData,
	fileToLineCountMap map[string]int, fileToByteCountMap map[string]int) {
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
			module := fileToModuleMap[fullFilename]
			filename := fullFilename[len(module.PackageDir)+1:]

			err := WriteMatchResult(MatchResultData{
				ProjectName:          module.ProjectName,
				ModuleImportPath:     module.ModuleImportPath,
				ModuleRegistry:       module.ModuleRegistry,
				ModuleVersion:        module.ModuleVersion,
				ModuleNumberGoFiles:  module.ModuleNumberGoFiles,
				ModuleCheckoutFolder: module.PackageDir,
				FileName:             filename,
				FileSizeBytes:        fileToByteCountMap[fullFilename],
				FileSizeLines:        fileToLineCountMap[fullFilename],
				FileImportsUnsafePkg: false, // TODO
				FileGoVetOutput:      "", // TODO
				Text:                 line.Data.Lines.Text,
				Context:              context,
				LineNumber:           line.Data.LineNumber,
				ByteOffset:           line.Data.LineNumber,
				MatchType:            line.Data.SubMatches[0].Match.Text,
			})

			if err != nil {
				_ = WriteErrorCondition(ErrorConditionData{
					Stage:            "ripgrep-parse",
					ProjectName:      module.ProjectName,
					ModuleImportPath: module.ModuleImportPath,
					FileName:         filename,
					Message:          err.Error(),
				})
				fmt.Println("SAVING ERROR!")
				continue
			}
		}
	}
}