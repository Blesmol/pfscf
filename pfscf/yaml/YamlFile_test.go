package yaml

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	yamlTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	yamlTestDir = filepath.Join(utils.GetExecutableDir(), "testdata")
}

func TestGetYamlFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "nonExistantFile.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNil(t, yFile)
			test.ExpectError(t, err)
		})

		t.Run("malformed file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "malformed.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNil(t, yFile)
			test.ExpectError(t, err)
		})

		t.Run("unknown fields", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "unknownFields.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNil(t, yFile)
			test.ExpectError(t, err)
		})

		t.Run("field type mismatch", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "fieldTypeMismatch.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNil(t, yFile)
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "empty.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			// empty file => no default section and no content
			test.ExpectNotSet(t, yFile.ID)
			test.ExpectNotSet(t, yFile.Description)
			test.ExpectNotSet(t, yFile.Inherit)
			test.ExpectEqual(t, len(yFile.Presets), 0)
			test.ExpectEqual(t, len(yFile.Content), 0)
		})

		t.Run("valid file", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "valid.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, yFile.ID, "myID")
			test.ExpectEqual(t, yFile.Description, "my Description")
			test.ExpectEqual(t, yFile.Inherit, "myInheritId")

			// presets
			test.ExpectEqual(t, len(yFile.Presets), 2)

			test.ExpectKeyExists(t, yFile.Presets, "preset0")
			p0 := yFile.Presets["preset0"]
			test.ExpectEqual(t, p0.Y1, 150.0)
			test.ExpectEqual(t, p0.Y2, 200.0)

			test.ExpectKeyExists(t, yFile.Presets, "preset1")
			p1 := yFile.Presets["preset1"]
			test.ExpectEqual(t, p1.Font, "my font")
			test.ExpectEqual(t, p1.Fontsize, 10.0)

			// number of entries in content array
			test.ExpectEqual(t, len(yFile.Content), 2)

			test.ExpectKeyExists(t, yFile.Content, "myId")
			c0 := yFile.Content["myId"]
			test.ExpectAllExportedSet(t, c0)
			test.ExpectEqual(t, c0.Type, "my Type")
			test.ExpectEqual(t, c0.Desc, "my Desc")
			test.ExpectEqual(t, c0.X1, 11.0)
			test.ExpectEqual(t, c0.Y1, 12.0)
			test.ExpectEqual(t, c0.X2, 13.0)
			test.ExpectEqual(t, c0.Y2, 14.0)
			test.ExpectEqual(t, c0.XPivot, 15.0)
			test.ExpectEqual(t, c0.Font, "my Font")
			test.ExpectEqual(t, c0.Fontsize, 16.0)
			test.ExpectEqual(t, c0.Align, "my Align")
			test.ExpectEqual(t, c0.Color, "my Color")
			test.ExpectEqual(t, c0.Example, "my Example")

			test.ExpectKeyExists(t, yFile.Content, "myOtherId")
			c1 := yFile.Content["myOtherId"]
			test.ExpectAllExportedSet(t, c1)
			test.ExpectEqual(t, c1.Type, "my other type")
			test.ExpectEqual(t, c1.Desc, "my other desc")
			test.ExpectEqual(t, c1.X1, 21.0)
			test.ExpectEqual(t, c1.Y1, 22.0)
			test.ExpectEqual(t, c1.X2, 23.0)
			test.ExpectEqual(t, c1.Y2, 24.0)
			test.ExpectEqual(t, c1.XPivot, 25.0)
			test.ExpectEqual(t, c1.Font, "my other font")
			test.ExpectEqual(t, c1.Fontsize, 26.0)
			test.ExpectEqual(t, c1.Align, "my other align")
			test.ExpectEqual(t, c1.Color, "my other color")
			test.ExpectEqual(t, c1.Example, "my other example")
		})

		t.Run("empty content entry", func(t *testing.T) {
			fileToTest := filepath.Join(yamlTestDir, "emptyContentEntry.yml")
			yFile, err := GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, 1, len(yFile.Content)) // one empty entry is included
			test.ExpectKeyExists(t, yFile.Content, "myId")
			c0 := yFile.Content["myId"]

			// check that all exported fields in the empty content entry are not set
			refStruct := reflect.ValueOf(c0)
			for i := 0; i < refStruct.NumField(); i++ {
				refField := refStruct.Field(i)
				if refField.CanAddr() && utils.IsSet(refField.Interface()) {
					t.Errorf("test.Expected ContentData field '%v' to be not set, but has value '%v' instead", refStruct.Type().Field(i).Name, refField.Interface())
				}
			}
		})
	})
}

func TestGetTemplateFilenamesFromDir(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
			fileNames, err := getTemplateFilenamesFromDir(dirToTest)

			test.ExpectError(t, err)
			test.ExpectNil(t, fileNames)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("nested dirs", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "nestedDirs")
			fileNames, err := getTemplateFilenamesFromDir(dirToTest)

			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(fileNames), 5)

			sort.Strings(fileNames) // for testing purposes, lets sort that list for easier comparison

			test.ExpectEqual(t, fileNames[0], filepath.Join(dirToTest, "BaR.YmL"))
			test.ExpectEqual(t, fileNames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
			test.ExpectEqual(t, fileNames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
			test.ExpectEqual(t, fileNames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
			test.ExpectEqual(t, fileNames[4], filepath.Join(dirToTest, "foo.yml"))
		})
	})
}

func TestGetTemplateFilesFromDir(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
			yFiles, err := GetTemplateFilesFromDir(dirToTest)

			test.ExpectError(t, err)
			test.ExpectNil(t, yFiles)
		})
	})

	t.Run("valid", func(t *testing.T) {
		dirToTest := filepath.Join(yamlTestDir, "nestedDirs")
		yFiles, err := GetTemplateFilesFromDir(dirToTest)

		test.ExpectNoError(t, err)
		test.ExpectEqual(t, len(yFiles), 5)

		// extract filenames for checking completeness
		var yFilenames []string
		for yFilename := range yFiles {
			yFilenames = append(yFilenames, yFilename)
		}
		sort.Strings(yFilenames)

		test.ExpectEqual(t, yFilenames[0], filepath.Join(dirToTest, "BaR.YmL"))
		test.ExpectEqual(t, yFilenames[1], filepath.Join(dirToTest, "dir1", "foo.yml"))
		test.ExpectEqual(t, yFilenames[2], filepath.Join(dirToTest, "dir2", "bar.yml"))
		test.ExpectEqual(t, yFilenames[3], filepath.Join(dirToTest, "dir3.yml", "foobar.yml"))
		test.ExpectEqual(t, yFilenames[4], filepath.Join(dirToTest, "foo.yml"))

		// basic check of the contents of one file
		foundTestFile := false
		for yFilename, yFile := range yFiles {
			if yFilename == filepath.Join(dirToTest, "foo.yml") {
				test.ExpectEqual(t, yFile.ID, "test")
				foundTestFile = true
			}
		}
		test.ExpectEqual(t, foundTestFile, true)
	})
}
