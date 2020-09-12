package base

import (
	"github.com/gocarina/gocsv"
	"os"
)

// these are global variables for any opened file, so that they can be used as singletons. Each file handle is
// accompanied by a boolean indicating if the CSV header has already been written, because the header should only
// be written once
var projectsFile *os.File
var projectsFileHeaderWritten = false
var packagesFile *os.File
var packagesFileHeaderWritten = false
var geigerFindingsFile *os.File
var geigerFindingsFileHeaderWritten = false
var grepFindingsFile *os.File
var grepFindingsFileHeaderWritten = false
var vetFindingsFile *os.File
var vetFindingsFileHeaderWritten = false
var gosecFindingsFile *os.File
var gosecFindingsFileHeaderWritten = false
var gosaferFindingsFile *os.File
var gosaferFindingsFileHeaderWritten = false
var astFindingsFile *os.File
var astFindingsFileHeaderWritten = false
var astFunctionsFile *os.File
var astFunctionsFileHeaderWritten = false
var astStatementsFile *os.File
var astStatementsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

/**
 * opens the projects file for writing and stores the handle in the global variable
 */
func OpenProjectsFile(projectsFilename string) error {
	var err error
	projectsFile, err = os.OpenFile(projectsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the packages file for writing and stores the handle in the global variable
 */
func OpenPackagesFile(packagesFilename string) error {
	var err error
	packagesFile, err = os.OpenFile(packagesFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the geiger analysis findings file for writing and stores the handle in the global variable
 */
func OpenGeigerFindingsFile(geigerFilename string) error {
	var err error
	geigerFindingsFile, err = os.OpenFile(geigerFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the grep analysis findings file for writing and stores the handle in the global variable
 */
func OpenGrepFindingsFile(grepFilename string) error {
	var err error
	grepFindingsFile, err = os.OpenFile(grepFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the go vet analysis findings file for writing and stores the handle in the global variable
 */
func OpenVetFindingsFile(vetFilename string) error {
	var err error
	vetFindingsFile, err = os.OpenFile(vetFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the gosec analysis findings file for writing and stores the handle in the global variable
 */
func OpenGosecFindingsFile(gosecFilename string) error {
	var err error
	gosecFindingsFile, err = os.OpenFile(gosecFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the go-safer analysis findings file for writing and stores the handle in the global variable
 */
func OpenGosaferFindingsFile(linterFilename string) error {
	var err error
	gosaferFindingsFile, err = os.OpenFile(linterFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * opens the AST analysis findings file for writing and stores the handle in the global variable
 */
func OpenAstFindingsFile(astFilename string) error {
	var err error
	astFindingsFile, err = os.OpenFile(astFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

/**
 * opens the AST analysis functions results file for writing and stores the handle in the global variable
 */
func OpenAstFunctionsFile(functionsFilename string) error {
	var err error
	astFunctionsFile, err = os.OpenFile(functionsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

/**
 * opens the AST analysis statements results file for writing and stores the handle in the global variable
 */
func OpenAstStatementsFile(statementsFilename string) error {
	var err error
	astStatementsFile, err = os.OpenFile(statementsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

/**
 * opens the error log file for writing and stores the handle in the global variable
 */
func OpenErrorConditionsFile(errorsFilename string) error {
	var err error
	errorConditionsFile, err = os.OpenFile(errorsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

/**
 * closes all files that are currently open
 */
func CloseFiles() {
	if projectsFile != nil {
		projectsFile.Close()
	}
	if packagesFile != nil {
		packagesFile.Close()
	}
	if geigerFindingsFile != nil {
		geigerFindingsFile.Close()
	}
	if grepFindingsFile != nil {
		grepFindingsFile.Close()
	}
	if vetFindingsFile != nil {
		vetFindingsFile.Close()
	}
	if gosecFindingsFile != nil {
		gosecFindingsFile.Close()
	}
	if gosaferFindingsFile != nil {
		gosaferFindingsFile.Close()
	}
	if astFindingsFile != nil {
		astFindingsFile.Close()
	}
	if astFunctionsFile != nil {
		astFunctionsFile.Close()
	}
	if astStatementsFile != nil {
		astStatementsFile.Close()
	}
	if errorConditionsFile != nil {
		errorConditionsFile.Close()
	}
}

/**
 * reads the project data from the specified file path to the CSV file
 */
func ReadProjects(filename string)([]*ProjectData, error) {
	// open the file and later close it again
	f, err := os.Open(filename)
	if err != nil {
		return []*ProjectData{}, err
	}
	defer f.Close()

	// initialize the list of projects
	var projects []*ProjectData

	// use the gocsv library to parse the projects data
	if err := gocsv.UnmarshalFile(f, &projects); err != nil {
		return []*ProjectData{}, err
	}

	// return the projects
	return projects, nil
}

/**
 * writes a project line to disk
 */
func WriteProject(project ProjectData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if projectsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]ProjectData{project}, projectsFile)
	} else {
		projectsFileHeaderWritten = true
		return gocsv.Marshal([]ProjectData{project}, projectsFile)
	}
}

/**
 * writes a package line to disk
 */
func WritePackage(module PackageData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if packagesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]PackageData{module}, packagesFile)
	} else {
		packagesFileHeaderWritten = true
		return gocsv.Marshal([]PackageData{module}, packagesFile)
	}
}

/**
 * writes a geiger analysis finding line to disk
 */
func WriteGeigerFinding(geigerFinding GeigerFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if geigerFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GeigerFindingData{geigerFinding}, geigerFindingsFile)
	} else {
		geigerFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GeigerFindingData{geigerFinding}, geigerFindingsFile)
	}
}

/**
 * writes a grep analysis finding line to disk
 */
func WriteGrepFinding(grepFinding GrepFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if grepFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GrepFindingData{grepFinding}, grepFindingsFile)
	} else {
		grepFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GrepFindingData{grepFinding}, grepFindingsFile)
	}
}

/**
 * writes a go vet analysis finding line to disk
 */
func WriteVetFinding(vetFinding VetFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if vetFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]VetFindingData{vetFinding}, vetFindingsFile)
	} else {
		vetFindingsFileHeaderWritten = true
		return gocsv.Marshal([]VetFindingData{vetFinding}, vetFindingsFile)
	}
}

/**
 * writes a gosec analysis finding line to disk
 */
func WriteGosecFinding(gosecFinding GosecFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if gosecFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GosecFindingData{gosecFinding}, gosecFindingsFile)
	} else {
		gosecFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GosecFindingData{gosecFinding}, gosecFindingsFile)
	}
}

/**
 * writes a go-safer analysis finding line to disk
 */
func WriteGosaferFinding(linterFinding GosaferFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if gosaferFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GosaferFindingData{linterFinding}, gosaferFindingsFile)
	} else {
		gosaferFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GosaferFindingData{linterFinding}, gosaferFindingsFile)
	}
}

/**
 * writes an AST analysis finding line to disk
 */
func WriteAstFinding(astFinding AstFindingData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if astFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]AstFindingData{astFinding}, astFindingsFile)
	} else {
		astFindingsFileHeaderWritten = true
		return gocsv.Marshal([]AstFindingData{astFinding}, astFindingsFile)
	}
}

/**
 * writes an AST analysis function results line to disk
 */
func WriteAstFunction(function AstFunctionData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if astFunctionsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]AstFunctionData{function}, astFunctionsFile)
	} else {
		astFunctionsFileHeaderWritten = true
		return gocsv.Marshal([]AstFunctionData{function}, astFunctionsFile)
	}
}

/**
 * writes an AST analysis statement results line to disk
 */
func WriteAstStatement(statement AstStatementData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if astStatementsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]AstStatementData{statement}, astStatementsFile)
	} else {
		astStatementsFileHeaderWritten = true
		return gocsv.Marshal([]AstStatementData{statement}, astStatementsFile)
	}
}

/**
 * writes an error log line to disk
 */
func WriteErrorCondition(errorCondition ErrorConditionData) error {
	// check if the header has already been written. If so, store the data without header. Otherwise, store including
	// the CSV header and flag that the header has now been written.
	if errorConditionsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]ErrorConditionData{errorCondition}, errorConditionsFile)
	} else {
		errorConditionsFileHeaderWritten = true
		return gocsv.Marshal([]ErrorConditionData{errorCondition}, errorConditionsFile)
	}
}