package analysis

import (
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

var packagesFile *os.File
var packagesFileHeaderWritten = false
var grepFindingsFile *os.File
var grepFindingsFileHeaderWritten = false
var vetFindingsFile *os.File
var vetFindingsFileHeaderWritten = false
var gosecFindingsFile *os.File
var gosecFindingsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

func openFiles(packagesFilename, grepFilename, vetFilename, gosecFilename, errorsFilename string) error {
	var err error

	packagesFile, err = os.OpenFile(packagesFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	grepFindingsFile, err = os.OpenFile(grepFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	vetFindingsFile, err = os.OpenFile(vetFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	gosecFindingsFile, err = os.OpenFile(gosecFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	packagesFile.Close()
	grepFindingsFile.Close()
	vetFindingsFile.Close()
	gosecFindingsFile.Close()
	errorConditionsFile.Close()
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	// https://yourbasic.org/golang/format-parse-string-time-date-example/
	date.Time, err = time.Parse("2006-01-02 15:04:05 -0700 MST", csv)
	return err
}

func readProjects(filename string)([]*ProjectData, error) {
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

func WritePackage(module PackageData) error {
	if packagesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]PackageData{module}, packagesFile)
	} else {
		packagesFileHeaderWritten = true
		return gocsv.Marshal([]PackageData{module}, packagesFile)
	}
}

func WriteGrepFinding(grepFinding GrepFindingData) error {
	if grepFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GrepFindingData{grepFinding}, grepFindingsFile)
	} else {
		grepFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GrepFindingData{grepFinding}, grepFindingsFile)
	}
}

func WriteVetFinding(vetFinding VetFindingData) error {
	if vetFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]VetFindingData{vetFinding}, vetFindingsFile)
	} else {
		vetFindingsFileHeaderWritten = true
		return gocsv.Marshal([]VetFindingData{vetFinding}, vetFindingsFile)
	}
}

func WriteGosecFinding(gosecFinding GosecFindingData) error {
	if gosecFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GosecFindingData{gosecFinding}, gosecFindingsFile)
	} else {
		gosecFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GosecFindingData{gosecFinding}, gosecFindingsFile)
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