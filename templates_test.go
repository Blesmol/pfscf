package main

import (
	"path/filepath"
	"reflect"
	"runtime/debug"
	"testing"
)

const (
	printCallStackOnFailingTest = false
)

var (
	yamlTestDir string
)

func init() {
	yamlTestDir = filepath.Join(GetExecutableDir(), "testdata", "yaml")
}

func callStack() {
	if printCallStackOnFailingTest {
		debug.PrintStack()
	}
}

func expectEqual(t *testing.T, got interface{}, exp interface{}) {
	if exp == got {
		return
	}
	callStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

func expectNotEqual(t *testing.T, got interface{}, notExp interface{}) {
	typeNotExp := reflect.TypeOf(notExp)
	typeGot := reflect.TypeOf(got)

	// we always require that both types are identical.
	// Without that, testing can be a real pain
	if typeNotExp != typeGot {
		callStack()
		t.Errorf("Types do not match! Expected '%v', got '%v'", typeNotExp, typeGot)
		return
	}

	if notExp == got {
		callStack()
		t.Errorf("Expected something different than '%v' (type %v)", notExp, typeNotExp)
	}
}

func expectNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if !reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		callStack()
		t.Error("Expected an error, got nil")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		callStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func expectNotSet(t *testing.T, got interface{}) {
	if IsSet(got) {
		callStack()
		t.Errorf("Expected not set, got '%v'", got)
	}
}

func expectAllSet(t *testing.T, got interface{}) {
	vGot := reflect.ValueOf(got)

	switch vGot.Kind() {
	case reflect.Struct:
		for i := 0; i < vGot.NumField(); i++ {
			field := vGot.Field(i)
			expectAllSet(t, field.Interface())
		}
	case reflect.Ptr:
		if IsSet(got) {
			expectAllSet(t, vGot.Elem().Interface())
		} else {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	default:
		if !IsSet(got) {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	}
}

func TestGetYamlFile_NonExistantFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "nonExistantFile.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_MalformedFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "malformed.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_ValidFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "valid.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// default values
	def := &(yFile.Default)
	expectAllSet(t, def)
	expectEqual(t, def.Type, "my default Type")
	expectEqual(t, def.ID, "my default Id")
	expectEqual(t, def.Desc, "my default Desc")
	expectEqual(t, def.X1, 91.0)
	expectEqual(t, def.Y1, 92.0)
	expectEqual(t, def.X2, 93.0)
	expectEqual(t, def.Y2, 94.0)
	expectEqual(t, def.Font, "my default Font")
	expectEqual(t, def.Fontsize, 95.0)
	expectEqual(t, def.Align, "my default Align")

	// number of entries in content array
	expectEqual(t, len(yFile.Content), 2)

	c0 := &(yFile.Content[0])
	expectAllSet(t, c0)
	expectEqual(t, c0.Type, "my Type")
	expectEqual(t, c0.ID, "my Id")
	expectEqual(t, c0.Desc, "my Desc")
	expectEqual(t, c0.X1, 11.0)
	expectEqual(t, c0.Y1, 12.0)
	expectEqual(t, c0.X2, 13.0)
	expectEqual(t, c0.Y2, 14.0)
	expectEqual(t, c0.Font, "my Font")
	expectEqual(t, c0.Fontsize, 15.0)
	expectEqual(t, c0.Align, "my Align")

	c1 := &(yFile.Content[1])
	expectAllSet(t, c1)
	expectEqual(t, c1.Type, "my other type")
	expectEqual(t, c1.ID, "my other id")
	expectEqual(t, c1.Desc, "my other desc")
	expectEqual(t, c1.X1, 21.0)
	expectEqual(t, c1.Y1, 22.0)
	expectEqual(t, c1.X2, 23.0)
	expectEqual(t, c1.Y2, 24.0)
	expectEqual(t, c1.Font, "my other font")
	expectEqual(t, c1.Fontsize, 25.0)
	expectEqual(t, c1.Align, "my other align")
}

func TestGetYamlFile_EmptyFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "empty.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// empty file => no default section and no content
	expectNotSet(t, yFile.Default)
	expectEqual(t, len(yFile.Content), 0)
}

func TestGetYamlFile_EmptyContentEntry(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "emptyContentEntry.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	expectEqual(t, 1, len(yFile.Content)) // one empty entry is included
	content0 := &(yFile.Content[0])

	// check that all fields in the empty content entry are not set
	refStruct := reflect.ValueOf(content0).Elem()
	for i := 0; i < refStruct.NumField(); i++ {
		refField := refStruct.Field(i)
		if IsSet(refField.Interface()) {
			t.Errorf("Expected ContentEntry field '%v' to be not set, but has value '%v' instead", refStruct.Type().Field(i).Name, refField.Interface())
		}
	}
}

func TestGetYamlFile_UnknownFields(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "unknownFields.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_FieldTypeMismatch(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "fieldTypeMismatch.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func getContentEntryWithDummyData(ceType string, ceID string) (ce ContentEntry) {
	ce.Type = ceType
	ce.ID = ceID
	ce.Desc = "Some Description"
	ce.X1 = 12.0
	ce.Y1 = 12.0
	ce.X2 = 24.0
	ce.Y2 = 24.0
	ce.Font = "Helvetica"
	ce.Fontsize = 14.0
	ce.Align = "LB"
	return ce
}

func TestContentEntryIsValid_emptyType(t *testing.T) {
	ce := getContentEntryWithDummyData("", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestContentEntryIsValid_invalidType(t *testing.T) {
	ce := getContentEntryWithDummyData("textCellX", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestContentEntryIsValid_validTextCell(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, true)
	expectNoError(t, err)
}

func TestContentEntryIsValid_textCellWithZeroedValues(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell", "foo")
	ce.Font = ""

	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestApplyDefaults(t *testing.T) {
	var ce ContentEntry
	ce.Type = "Foo"
	ce.Fontsize = 10.0
	ce.Font = ""

	var defaults ContentEntry
	defaults.Type = "Bar"
	defaults.X1 = 5.0
	defaults.Y1 = 0.0
	defaults.Font = "Dingbats"
	defaults.Fontsize = 20.0

	ce.applyDefaults(defaults)

	expectEqual(t, ce.Type, "Foo")
	expectNotSet(t, ce.ID)
	expectNotSet(t, ce.Desc)
	expectEqual(t, ce.X1, 5.0)
	expectNotSet(t, ce.Y1)
	expectNotSet(t, ce.X2)
	expectNotSet(t, ce.Y2)
	expectEqual(t, ce.Font, "Dingbats")
	expectEqual(t, ce.Fontsize, 10.0)
	expectNotSet(t, ce.Align)
}
