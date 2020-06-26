package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

var (
	yamlTestDir string
)

func init() {
	yamlTestDir = filepath.Join(GetExecutableDir(), "testdata", "yaml")
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

	expectEqual(t, yFile.ID, "myID")
	expectEqual(t, yFile.Description, "my Description")

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
	expectNotSet(t, yFile.ID)
	expectNotSet(t, yFile.Description)
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
