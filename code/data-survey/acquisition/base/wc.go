package base

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

/**
 * runs wc on the specified files and returns their line counts as a hash map using the file name as key
 */
func CountLines(files []string) (map[string]int, error) {
	fmt.Println("  running wc to count LOC")

	// build the wc command as wc -l file1 file2
	args := []string{"-l"}
	args = append(args, files...)
	cmd := exec.Command("wc", args...)

	// run wc and capture the output
	wcOutput, _ := cmd.Output()

	// split the output by lines because wc will return each line count in a separate line
	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	// initialize the map of counts
	filesToLineCount := map[string]int{}

	// go through the output lines
	for _, outputLine := range outputLines {
		// file name and count are separated by spaces, therefore separate by those
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		// try to read the count which is the first component
		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		// save the count in the hash map
		filesToLineCount[components[1]] = count
	}

	return filesToLineCount, nil
}

/**
 * runs wc on the specified files and returns their byte counts as a hash map using the file name as key
 */
func CountBytes(files []string) (map[string]int, error) {
	// build the wc command as wc -c file1 file2
	args := []string{"-c"}
	args = append(args, files...)
	cmd := exec.Command("wc", args...)

	// run wc and capture the output
	wcOutput, _ := cmd.Output()

	// split the output by lines because wc will return each line count in a separate line
	outputLines := strings.Split(string(wcOutput), "\n")

	// remove count summary and trailing newline
	outputLines = outputLines[:len(outputLines)-2]

	// initialize the map of counts
	filesToByteCount := map[string]int{}

	// go through the output lines
	for _, outputLine := range outputLines {
		// file name and count are separated by spaces, therefore separate by those
		components := strings.Split(strings.Trim(outputLine, " "), " ")

		// try to read the count which is the first component
		count, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}

		// save the count in the hash map
		filesToByteCount[components[1]] = count
	}

	return filesToByteCount, nil
}