package lexical

import (
	"github.com/gocarina/gocsv"
	"os"
)

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
var linterFindingsFile *os.File
var linterFindingsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

func OpenPackagesFile(packagesFilename string) error {
	var err error
	packagesFile, err = os.OpenFile(packagesFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func OpenGeigerFindingsFile(geigerFilename string) error {
	var err error
	geigerFindingsFile, err = os.OpenFile(geigerFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func openGrepFindingsFile(grepFilename string) error {
	var err error
	grepFindingsFile, err = os.OpenFile(grepFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func openVetFindingsFile(vetFilename string) error {
	var err error
	vetFindingsFile, err = os.OpenFile(vetFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func openGosecFindingsFile(gosecFilename string) error {
	var err error
	gosecFindingsFile, err = os.OpenFile(gosecFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func openLinterFindingsFile(linterFilename string) error {
	var err error
	linterFindingsFile, err = os.OpenFile(linterFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func OpenErrorConditionsFile(errorsFilename string) error {
	var err error
	errorConditionsFile, err = os.OpenFile(errorsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

func CloseFiles() {
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
	if linterFindingsFile != nil {
		linterFindingsFile.Close()
	}
	if errorConditionsFile != nil {
		errorConditionsFile.Close()
	}
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

func WritePackage(module PackageData) error {
	if packagesFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]PackageData{module}, packagesFile)
	} else {
		packagesFileHeaderWritten = true
		return gocsv.Marshal([]PackageData{module}, packagesFile)
	}
}

func WriteGeigerFinding(geigerFinding GeigerFindingData) error {
	if geigerFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]GeigerFindingData{geigerFinding}, geigerFindingsFile)
	} else {
		geigerFindingsFileHeaderWritten = true
		return gocsv.Marshal([]GeigerFindingData{geigerFinding}, geigerFindingsFile)
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

func WriteLinterFinding(linterFinding LinterFindingData) error {
	if linterFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]LinterFindingData{linterFinding}, linterFindingsFile)
	} else {
		linterFindingsFileHeaderWritten = true
		return gocsv.Marshal([]LinterFindingData{linterFinding}, linterFindingsFile)
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