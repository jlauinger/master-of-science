package ast

import (
	"github.com/gocarina/gocsv"
	"os"
)

var astFindingsFile *os.File
var astFindingsFileHeaderWritten = false
var functionsFile *os.File
var functionsFileHeaderWritten = false
var statementsFile *os.File
var statementsFileHeaderWritten = false
var errorConditionsFile *os.File
var errorConditionsFileHeaderWritten = false

func openAstFindingsFile(astFilename string) error {
	var err error
	astFindingsFile, err = os.OpenFile(astFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

func openFunctionsFile(functionsFilename string) error {
	var err error
	functionsFile, err = os.OpenFile(functionsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

func openStatementsFile(statementsFilename string) error {
	var err error
	statementsFile, err = os.OpenFile(statementsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

func openErrorConditionsFile(errorsFilename string) error {
	var err error
	errorConditionsFile, err = os.OpenFile(errorsFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return nil
}

func closeFiles() {
	if astFindingsFile != nil {
		astFindingsFile.Close()
	}
	if functionsFile != nil {
		functionsFile.Close()
	}
	if statementsFile != nil {
		statementsFile.Close()
	}
	if errorConditionsFile != nil {
		errorConditionsFile.Close()
	}
}

func WriteAstFinding(astFinding FindingData) error {
	if astFindingsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]FindingData{astFinding}, astFindingsFile)
	} else {
		astFindingsFileHeaderWritten = true
		return gocsv.Marshal([]FindingData{astFinding}, astFindingsFile)
	}
}

func WriteFunction(function FunctionData) error {
	if functionsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]FunctionData{function}, functionsFile)
	} else {
		functionsFileHeaderWritten = true
		return gocsv.Marshal([]FunctionData{function}, functionsFile)
	}
}

func WriteStatement(statement StatementData) error {
	if statementsFileHeaderWritten {
		return gocsv.MarshalWithoutHeaders([]StatementData{statement}, statementsFile)
	} else {
		statementsFileHeaderWritten = true
		return gocsv.Marshal([]StatementData{statement}, statementsFile)
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