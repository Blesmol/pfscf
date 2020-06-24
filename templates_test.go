package main

import (
	"path/filepath"
	"reflect"
	"runtime/debug"
	"testing"
)

var templateTestDir string

func init() {
	templateTestDir = filepath.Join(GetExecutableDir(), "testdata", "templates")
}

func expectEqual(t *testing.T, exp interface{}, got interface{}) {
	if exp == got {
		return
	}
	debug.PrintStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

func expectNotEqual(t *testing.T, notExp interface{}, got interface{}) {
	typeNotExp := reflect.TypeOf(notExp)
	typeGot := reflect.TypeOf(got)

	// we always require that both types are identical.
	// Without that, testing can be a real pain
	if typeNotExp != typeGot {
		debug.PrintStack()
		t.Errorf("Types do not match! Expected '%v', got '%v'", typeNotExp, typeGot)
		return
	}

	if notExp == got {
		debug.PrintStack()
		t.Errorf("Expected something different than '%v' (type %v)", notExp, typeNotExp)
	}
}

func expectNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if !reflect.ValueOf(got).IsNil() {
		debug.PrintStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if reflect.ValueOf(got).IsNil() {
		debug.PrintStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		debug.PrintStack()
		t.Error("Expected an error, got nil")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func TestGetYamlFile_NonExistantFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "nonExistantFile.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_MalformedFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "malformed.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_ValidFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "valid.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// test values contained in the yaml file
	expectEqual(t, "Helvetica", *yFile.Default.Font)
	expectEqual(t, 14.0, *yFile.Default.Fontsize)
	expectEqual(t, 2, len(yFile.Content))

	content0 := &(yFile.Content[0])
	expectEqual(t, "foo", *content0.ID)
	expectEqual(t, "textCell", *content0.Type)

	content1 := &(yFile.Content[1])
	expectEqual(t, "bar", *content1.ID)
	expectEqual(t, "textCell", *content1.Type)
}

func TestGetYamlFile_EmptyFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "empty.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// empty file => no default section and no content
	expectNil(t, yFile.Default)
	expectEqual(t, 0, len(yFile.Content))
}

func TestGetYamlFile_EmptyContentEntry(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "emptyContentEntry.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	expectEqual(t, 1, len(yFile.Content)) // one empty entry is included
	content0 := &(yFile.Content[0])

	// check that all fields in the empty content entry are nil ptrs
	refStruct := reflect.ValueOf(content0).Elem()
	for i := 0; i < refStruct.NumField(); i++ {
		refField := refStruct.Field(i)

		if refField.Kind() != reflect.Ptr {
			t.Errorf("Expected ContentEntry field '%v' to be a Ptr, but has kind '%v' instead", refStruct.Type().Field(i).Name, refField.Kind())
		} else {
			// is Ptr, check for nil
			if !refField.IsNil() {
				t.Errorf("Expected ContentEntry field '%v' to be a nil ptr, but points to value '%v' instead", refStruct.Type().Field(i).Name, reflect.Indirect(refField.Elem()))
			}
		}
	}
}

func TestGetYamlFile_UnknownFields(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "unknownFields.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_FieldTypeMismatch(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "fieldTypeMismatch.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func getContentEntryWithDummyData(ceType *string, ceId *string) (ce ContentEntry) {
	ce.Type = ceType
	ce.ID = ceId
	ce.Desc = new(string); *ce.Desc = "Some Description"
	ce.X1 = new(float64); *ce.X1 = 12.0
	ce.Y1 = new(float64); *ce.Y1 = 12.0
	ce.X2 = new(float64); *ce.X2 = 24.0
	ce.Y2 = new(float64); *ce.Y2 = 24.0
	ce.Font = new(string); *ce.Font = "Helvetica"
	ce.Fontsize = new(float64); *ce.Fontsize = 14.0
	ce.Align = new(string); *ce.Align = "LB"
	return ce
}

func TestContentEntryIsValid_emptyType(t *testing.T) {
	ceId := "foo"

	ce := getContentEntryWithDummyData(nil, &ceId)
	isValid, err := ce.IsValid()

	expectEqual(t, false, isValid)
	expectError(t, err)
}

func TestContentEntryIsValid_invalidType(t *testing.T) {
	ceType := "textCellX"
	ceId := "foo"

	ce := getContentEntryWithDummyData(&ceType, &ceId)
	isValid, err := ce.IsValid()

	expectEqual(t, false, isValid)
	expectError(t, err)
}

func TestContentEntryIsValid_validTextCell(t *testing.T) {
	ceType := "textCell"
	ceId := "foo"

	ce := getContentEntryWithDummyData(&ceType, &ceId)
	isValid, err := ce.IsValid()

	expectEqual(t, true, isValid)
	expectNoError(t, err)
}

func TestContentEntryIsValid_textCellWithMissingValues(t *testing.T) {
	ceType := "textCell"
	ceId := "foo"

	ce := getContentEntryWithDummyData(&ceType, &ceId)
	ce.Font = nil

	isValid, err := ce.IsValid()

	expectEqual(t, false, isValid)
	expectError(t, err)
}

func TestContentEntryIsValid_textCellWithZeroedValues(t *testing.T) {
	ceType := "textCell"
	ceId := "foo"

	ce := getContentEntryWithDummyData(&ceType, &ceId)
	emptyString := ""
	ce.Font = &emptyString

	isValid, err := ce.IsValid()

	expectEqual(t, false, isValid)
	expectError(t, err)
}

func TestApplyDefaults(t *testing.T) {
	var ce ContentEntry
	ce.Type = new(string); *ce.Type = "Foo"
	ce.Fontsize = new(float64); *ce.Fontsize = 10.0
	ce.Font = new(string) // intentionally left empty

	var defaults ContentEntry
	defaults.Type = new(string); *defaults.Type = "Bar"
	defaults.X1 = new(float64); *defaults.X1 = 5.0
	defaults.Y1 = new(float64); // intentionally left empty
	defaults.Font = new(string); *defaults.Font = "Dingbats"
	defaults.Fontsize = new(float64); *defaults.Fontsize = 20.0

	ce.applyDefaults(&defaults)

	expectEqual(t, *ce.Type, "Foo")
	expectNil(t, ce.ID)
	expectNil(t, ce.Desc)
	expectEqual(t, *ce.X1, 5.0)
	expectNil(t, ce.Y1)
	expectNil(t, ce.X2)
	expectNil(t, ce.Y2)
	expectEqual(t, *ce.Font, "Dingbats")
	expectEqual(t, *ce.Fontsize, 10.0)
	expectNil(t, ce.Align)
}
