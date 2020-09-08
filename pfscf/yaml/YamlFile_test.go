package yaml

import (
	"path/filepath"
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

func TestGetTemplateFilenamesFromDir(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "doesNotExist")
			fileNames, err := GetYamlFilenamesFromDir(dirToTest)

			test.ExpectError(t, err)
			test.ExpectNil(t, fileNames)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("nested dirs", func(t *testing.T) {
			dirToTest := filepath.Join(yamlTestDir, "nestedDirs")
			fileNames, err := GetYamlFilenamesFromDir(dirToTest)

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
