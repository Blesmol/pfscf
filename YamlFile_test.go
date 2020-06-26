package main

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

var (
	yamlTestDir string
)

func init() {
	SetIsTestEnvironment()
	yamlTestDir = filepath.Join(GetExecutableDir(), "testdata", "yaml")
}

func Test_GetYamlFile_NonExistantFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "nonExistantFile.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func Test_GetYamlFile_MalformedFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "malformed.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func Test_GetYamlFile_ValidFile(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "valid.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	expectEqual(t, yFile.fileName, fileToTest)

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

func Test_GetYamlFile_EmptyFile(t *testing.T) {
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

func Test_GetYamlFile_EmptyContentEntry(t *testing.T) {
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

func Test_GetYamlFile_UnknownFields(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "unknownFields.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func Test_GetYamlFile_FieldTypeMismatch(t *testing.T) {
	fileToTest := filepath.Join(yamlTestDir, "fieldTypeMismatch.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func Test_GetTemplateFilenamesFromDir_ValidDir(t *testing.T) {
	dirToTest := filepath.Join(yamlTestDir, "nested")
	fileNames, err := GetTemplateFilenamesFromDir(dirToTest)

	expectNoError(t, err)
	expectEqual(t, len(fileNames), 5)

	sort.Strings(fileNames) // for testing purposes, lets sort that list for easier comparison

	expectEqual(t, fileNames[0], filepath.Join(dirToTest, "BaR.YmL"))
	expectEqual(t, fileNames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
	expectEqual(t, fileNames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
	expectEqual(t, fileNames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
	expectEqual(t, fileNames[4], filepath.Join(dirToTest, "foo.yml"))
}

func Test_GetTemplateFilenamesFromDir_NonExistantDir(t *testing.T) {
	dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
	fileNames, err := GetTemplateFilenamesFromDir(dirToTest)

	expectError(t, err)
	expectNil(t, fileNames)
}

func Test_GetTemplateFilesFromDir_ValidDir(t *testing.T) {
	dirToTest := filepath.Join(yamlTestDir, "nested")
	yFiles, err := GetTemplateFilesFromDir(dirToTest)

	expectNoError(t, err)
	expectEqual(t, len(yFiles), 5)

	// extract filenames for checking completeness
	var fileNames []string
	for _, yFile := range yFiles {
		fileNames = append(fileNames, yFile.fileName)
	}
	sort.Strings(fileNames)

	expectEqual(t, fileNames[0], filepath.Join(dirToTest, "BaR.YmL"))
	expectEqual(t, fileNames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
	expectEqual(t, fileNames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
	expectEqual(t, fileNames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
	expectEqual(t, fileNames[4], filepath.Join(dirToTest, "foo.yml"))

	// basic check of the contents of one file
	foundTestFile := false
	for _, yFile := range yFiles {
		if yFile.fileName == filepath.Join(dirToTest, "foo.yml") {
			expectEqual(t, yFile.ID, "test")
			foundTestFile = true
		}
	}
	expectEqual(t, foundTestFile, true)
}

func Test_GetTemplateFilesFromDir_NonExistantDir(t *testing.T) {
	dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
	yFiles, err := GetTemplateFilesFromDir(dirToTest)

	expectError(t, err)
	expectNil(t, yFiles)
}
