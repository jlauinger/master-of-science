package geiger

import (
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/base"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"strings"
)

/**
 * returns a one lines and +/- 5 lines source code context for the given AST ndoe
 */
func getCodeContext(parsedPkg *packages.Package, n ast.Node) (string, string) {
	// identify the line number and file name for the AST node from the parsing file set
	nodePosition := parsedPkg.Fset.File(n.Pos()).Position(n.Pos())
	lineNumber := nodePosition.Line  // 1-based
	filename := nodePosition.Filename

	// read in all the source file and split the lines
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(fileData), "\n")

	// assert that the line number is actually contained in the file
	if lineNumber > len(fileLines) {
		return "invalid-line-number", "invalid-line-number"
	}

	// define start and end line for the +/- 5 lines context. Take care of the file boundaries, i.e. if there are less
	// than 5 lines left we need to shorten the context
	startLine := base.Max(1, lineNumber - 5)
	endLine := base.Min(len(fileLines), lineNumber + 6)

	// extract the 1 and 5 lines context from the lines list
	text := strings.Trim(fileLines[lineNumber-1], "\n\t")
	context := strings.Join(fileLines[startLine-1:endLine-1], "\n")

	return text, context
}
