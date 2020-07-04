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
	SetIsTestEnvironment(true)
	yamlTestDir = filepath.Join(GetExecutableDir(), "testdata", "YamlFile")
}

func TestGetYamlFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "nonExistantFile.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNil(t, yFile)
			expectError(t, err)
		})

		t.Run("malformed file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "malformed.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNil(t, yFile)
			expectError(t, err)
		})

		t.Run("unknown fields", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "unknownFields.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNil(t, yFile)
			expectError(t, err)
		})

		t.Run("field type mismatch", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "fieldTypeMismatch.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNil(t, yFile)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "empty.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			// empty file => no default section and no content
			expectNotSet(t, yFile.ID)
			expectNotSet(t, yFile.Description)
			expectNotSet(t, yFile.Inherit)
			expectEqual(t, len(yFile.Presets), 0)
			expectEqual(t, len(yFile.Content), 0)
		})

		t.Run("valid file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "valid.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			expectEqual(t, yFile.ID, "myID")
			expectEqual(t, yFile.Description, "my Description")
			expectEqual(t, yFile.Inherit, "myInheritId")

			// presets
			expectEqual(t, len(yFile.Presets), 2)

			expectKeyExists(t, yFile.Presets, "preset0")
			p0 := yFile.Presets["preset0"]
			expectEqual(t, p0.Y1, 150.0)
			expectEqual(t, p0.Y2, 200.0)

			expectKeyExists(t, yFile.Presets, "preset1")
			p1 := yFile.Presets["preset1"]
			expectEqual(t, p1.Font, "my font")
			expectEqual(t, p1.Fontsize, 10.0)

			// number of entries in content array
			expectEqual(t, len(yFile.Content), 2)

			expectKeyExists(t, yFile.Content, "myId")
			c0 := yFile.Content["myId"]
			expectAllExportedSet(t, c0)
			expectEqual(t, c0.Type, "my Type")
			expectEqual(t, c0.Desc, "my Desc")
			expectEqual(t, c0.X1, 11.0)
			expectEqual(t, c0.Y1, 12.0)
			expectEqual(t, c0.X2, 13.0)
			expectEqual(t, c0.Y2, 14.0)
			expectEqual(t, c0.XPivot, 15.0)
			expectEqual(t, c0.Font, "my Font")
			expectEqual(t, c0.Fontsize, 16.0)
			expectEqual(t, c0.Align, "my Align")
			expectEqual(t, c0.Example, "my Example")

			expectKeyExists(t, yFile.Content, "myOtherId")
			c1 := yFile.Content["myOtherId"]
			expectAllExportedSet(t, c1)
			expectEqual(t, c1.Type, "my other type")
			expectEqual(t, c1.Desc, "my other desc")
			expectEqual(t, c1.X1, 21.0)
			expectEqual(t, c1.Y1, 22.0)
			expectEqual(t, c1.X2, 23.0)
			expectEqual(t, c1.Y2, 24.0)
			expectEqual(t, c1.XPivot, 25.0)
			expectEqual(t, c1.Font, "my other font")
			expectEqual(t, c1.Fontsize, 26.0)
			expectEqual(t, c1.Align, "my other align")
			expectEqual(t, c1.Example, "my other example")
		})

		t.Run("empty content entry", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "emptyContentEntry.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			expectEqual(t, 1, len(yFile.Content)) // one empty entry is included
			expectKeyExists(t, yFile.Content, "myId")
			c0 := yFile.Content["myId"]

			// check that all exported fields in the empty content entry are not set
			refStruct := reflect.ValueOf(c0)
			for i := 0; i < refStruct.NumField(); i++ {
				refField := refStruct.Field(i)
				if refField.CanAddr() && IsSet(refField.Interface()) {
					t.Errorf("Expected ContentData field '%v' to be not set, but has value '%v' instead", refStruct.Type().Field(i).Name, refField.Interface())
				}
			}
		})
	})
}

func TestGetTemplateFilenamesFromDir(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
			fileNames, err := GetTemplateFilenamesFromDir(dirToTest)

			expectError(t, err)
			expectNil(t, fileNames)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("nested dirs", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "nestedDirs")
			fileNames, err := GetTemplateFilenamesFromDir(dirToTest)

			expectNoError(t, err)
			expectEqual(t, len(fileNames), 5)

			sort.Strings(fileNames) // for testing purposes, lets sort that list for easier comparison

			expectEqual(t, fileNames[0], filepath.Join(dirToTest, "BaR.YmL"))
			expectEqual(t, fileNames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
			expectEqual(t, fileNames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
			expectEqual(t, fileNames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
			expectEqual(t, fileNames[4], filepath.Join(dirToTest, "foo.yml"))
		})
	})
}

func TestGetTemplateFilesFromDir(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
			yFiles, err := GetTemplateFilesFromDir(dirToTest)

			expectError(t, err)
			expectNil(t, yFiles)
		})
	})

	t.Run("valid", func(t *testing.T) {
		dirToTest := filepath.Join(yamlTestDir, "nestedDirs")
		yFiles, err := GetTemplateFilesFromDir(dirToTest)

		expectNoError(t, err)
		expectEqual(t, len(yFiles), 5)

		// extract filenames for checking completeness
		var yFilenames []string
		for yFilename := range yFiles {
			yFilenames = append(yFilenames, yFilename)
		}
		sort.Strings(yFilenames)

		expectEqual(t, yFilenames[0], filepath.Join(dirToTest, "BaR.YmL"))
		expectEqual(t, yFilenames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
		expectEqual(t, yFilenames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
		expectEqual(t, yFilenames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
		expectEqual(t, yFilenames[4], filepath.Join(dirToTest, "foo.yml"))

		// basic check of the contents of one file
		foundTestFile := false
		for yFilename, yFile := range yFiles {
			if yFilename == filepath.Join(dirToTest, "foo.yml") {
				expectEqual(t, yFile.ID, "test")
				foundTestFile = true
			}
		}
		expectEqual(t, foundTestFile, true)
	})
}
