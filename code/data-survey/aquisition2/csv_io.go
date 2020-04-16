package main

import (
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

var modulesFile *os.File
var modulesFileHeaderWritten = false
var matchesFile *os.File
var matchesFileHeaderWritten = false
var vetResultsFile *os.File
var vetResultsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

func openFiles(modulesFilename, matchesFilename, vetResultsFilename, errorsFilename string) error {
	var err error

	modulesFile, err = os.OpenFile(modulesFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	matchesFile, err = os.OpenFile(matchesFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	vetResultsFile, err = os.OpenFile(vetResultsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	errorConditionsFile, err = os.OpenFile(errorsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return nil
}

func closeFiles() {
	modulesFile.Close()
	matchesFile.Close()
	vetResultsFile.Close()
	errorConditionsFile.Close()
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	// https://yourbasic.org/golang/format-parse-string-time-date-example/
	date.Time, err = time.Parse("2006-01-02 15:04:05 -0700 MST", csv)
	return err
}

func ReadProjects(filename string)([]*ProjectData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []*ProjectData{}, err
	}
	defer f.Close()

	var projects []*ProjectData

	if err := gocsv.UnmarshalFile(f, &projects); err != nil {
		return []*ProjectData{}, err
	}

	return projects, nil
}

func WriteModule(module ModuleData) error {
	if modulesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]ModuleData{module}, modulesFile)
	} else {
		modulesFileHeaderWritten = true
		return gocsv.Marshal([]ModuleData{module}, modulesFile)
	}
}

func WriteMatchResult(matchResult MatchResultData) error {
	if matchesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]MatchResultData{matchResult}, matchesFile)
	} else {
		matchesFileHeaderWritten = true
		return gocsv.Marshal([]MatchResultData{matchResult}, matchesFile)
	}
}

func WriteVetFinding(vetFinding VetFindingData) error {
	if vetResultsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]VetFindingData{vetFinding}, vetResultsFile)
	} else {
		vetResultsFileHeaderWritten = true
		return gocsv.Marshal([]VetFindingData{vetFinding}, vetResultsFile)
	}
}

func WriteErrorCondition(errorCondition ErrorConditionData) error {
	if errorConditionsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]ErrorConditionData{errorCondition}, errorConditionsFile)
	} else {
		errorConditionsFileHeaderWritten = true
		return gocsv.Marshal([]ErrorConditionData{errorCondition}, errorConditionsFile)
	}
}