package geiger

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func countLines(files []string) (map[string]int, error) {
	fmt.Println("  running wc to count LOC...")

	args := []string{"-l"}
	args = append(args, files...)

	cmd := exec.Command("wc", args...)

	wcOutput, _ := cmd.Output()

	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	filesToLineCount := map[string]int{}

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesToLineCount[components[1]] = count
	}

	return filesToLineCount, nil
}

func countBytes(files []string) (map[string]int, error) {
	args := []string{"-c"}
	args = append(args, files...)

	cmd := exec.Command("wc", args...)

	wcOutput, _ := cmd.Output()

	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	filesToByteCount := map[string]int{}

	for _, outputLine := range outputLines {
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		filesToByteCount[components[1]] = count
	}

	return filesToByteCount, nil
}