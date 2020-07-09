package eval2

import (
	"github.com/stg-tud/thesis-2020-lauinger-code/data-survey/acquisition/lexical"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"strings"
)

func getCodeLine(parsedPkg *packages.Package, n ast.Node) string {
	file := parsedPkg.Fset.File(n.Pos())
	lineNumber := file.Position(n.Pos()).Line  // 1-based

	startLine := lineNumber
	endLine := lexical.Min(file.LineCount(), lineNumber + 1)

	start := file.Position(file.LineStart(startLine)).Offset
	end := file.Position(file.LineStart(endLine)).Offset
	length := end - start

	filename := file.Name()

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Seek(int64(start), 0)
	if err != nil {
		panic(err)
	}
	line := make([]byte, length)
	_, err = io.ReadAtLeast(f, line, length)
	if err != nil {
		panic(err)
	}

	return strings.Trim(string(line), "\n\t ")
}

func getCodeContext(parsedPkg *packages.Package, n ast.Node) string {
	file := parsedPkg.Fset.File(n.Pos())
	lineNumber := file.Position(n.Pos()).Line  // 1-based

	startLine := lexical.Max(1, lineNumber - 5)
	endLine := lexical.Min(file.LineCount(), lineNumber + 6)

	start := file.Position(file.LineStart(startLine)).Offset
	end := file.Position(file.LineStart(endLine)).Offset
	length := end - start

	filename := file.Name()

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Seek(int64(start), 0)
	if err != nil {
		panic(err)
	}
	line := make([]byte, length)
	_, err = io.ReadAtLeast(f, line, length)
	if err != nil {
		panic(err)
	}

	return strings.Trim(string(line), "\n\t ")
}