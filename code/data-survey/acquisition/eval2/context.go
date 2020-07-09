package eval2

import (
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"strings"
)

func getCodeContext(parsedPkg *packages.Package, n ast.Node) (string, string) {
	nodePosition := parsedPkg.Fset.File(n.Pos()).Position(n.Pos())
	lineNumber := nodePosition.Line  // 1-based
	filename := nodePosition.Filename

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(fileData), "\n")

	if lineNumber > len(fileLines) {
		return "invalid-line-number", "invalid-line-number"
	}

	startLine := lexical.Max(1, lineNumber - 5)
	endLine := lexical.Min(len(fileLines), lineNumber + 6)

	text := strings.Trim(fileLines[lineNumber - 1], "\n\t")
	context := strings.Join(fileLines[startLine - 1 : endLine], "\n")

	return text, context
}
