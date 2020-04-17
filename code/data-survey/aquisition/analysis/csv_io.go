package analysis

import (
	"data-aquisition/base"
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
var gosecResultsFile *os.File
var gosecResultsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

func openFiles(modulesFilename, matchesFilename, vetResultsFilename, gosecResultsFilename, errorsFilename string) error {
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

	gosecResultsFile, err = os.OpenFile(gosecResultsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	gosecResultsFile.Close()
	errorConditionsFile.Close()
}

func (date *base.DateTime) UnmarshalCSV(csv string) (err error) {
	// https://yourbasic.org/golang/format-parse-string-time-date-example/
	date.Time, err = time.Parse("2006-01-02 15:04:05 -0700 MST", csv)
	return err
}

func ReadProjects(filename string)([]*base.ProjectData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []*base.ProjectData{}, err
	}
	defer f.Close()

	var projects []*base.ProjectData

	if err := gocsv.UnmarshalFile(f, &projects); err != nil {
		return []*base.ProjectData{}, err
	}

	return projects, nil
}

func WriteModule(module base.ModuleData) error {
	if modulesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]base.ModuleData{module}, modulesFile)
	} else {
		modulesFileHeaderWritten = true
		return gocsv.Marshal([]base.ModuleData{module}, modulesFile)
	}
}

func WriteMatchResult(matchResult base.MatchResultData) error {
	if matchesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]base.MatchResultData{matchResult}, matchesFile)
	} else {
		matchesFileHeaderWritten = true
		return gocsv.Marshal([]base.MatchResultData{matchResult}, matchesFile)
	}
}

func WriteVetFinding(vetFinding base.VetFindingData) error {
	if vetResultsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]base.VetFindingData{vetFinding}, vetResultsFile)
	} else {
		vetResultsFileHeaderWritten = true
		return gocsv.Marshal([]base.VetFindingData{vetFinding}, vetResultsFile)
	}
}

func WriteGosecFinding(gosecFinding base.GosecFindingData) error {
	if gosecResultsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]base.GosecFindingData{gosecFinding}, gosecResultsFile)
	} else {
		gosecResultsFileHeaderWritten = true
		return gocsv.Marshal([]base.GosecFindingData{gosecFinding}, gosecResultsFile)
	}
}

func WriteErrorCondition(errorCondition base.ErrorConditionData) error {
	if errorConditionsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]base.ErrorConditionData{errorCondition}, errorConditionsFile)
	} else {
		errorConditionsFileHeaderWritten = true
		return gocsv.Marshal([]base.ErrorConditionData{errorCondition}, errorConditionsFile)
	}
}