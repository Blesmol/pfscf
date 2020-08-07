package main

import (
	"path/filepath"
	"strings"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

var (
	chronicleTemplateTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	chronicleTemplateTestDir = filepath.Join(utils.GetExecutableDir(), "testdata", "ChronicleTemplate")
}

func getCTfromYamlFile(t *testing.T, filename string) (ct *ChronicleTemplate) {
	t.Helper()

	fileToTest := filepath.Join(chronicleTemplateTestDir, filename)
	yFile, err := yaml.GetYamlFile(fileToTest)
	test.ExpectNotNil(t, yFile)
	test.ExpectNoError(t, err)

	ct, err = NewChronicleTemplate(filename, yFile)
	test.ExpectNotNil(t, ct)
	test.ExpectNoError(t, err)
	return ct
}

func TestNewChronicleTemplate(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("no filename", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "basic.yml")
			yFile, err := yaml.GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			ct, err := NewChronicleTemplate("", yFile)
			test.ExpectNil(t, ct)
			test.ExpectError(t, err)
		})

		t.Run("nil file", func(t *testing.T) {
			ct, err := NewChronicleTemplate("file.yml", nil)
			test.ExpectNil(t, ct)
			test.ExpectError(t, err)
		})

		t.Run("empty fields", func(t *testing.T) {
			filenames := []string{"emptyId.yml", "emptyDescription.yml"}

			for _, filename := range filenames {
				fileToTest := filepath.Join(chronicleTemplateTestDir, filename)
				yFile, err := yaml.GetYamlFile(fileToTest)

				test.ExpectNotNil(t, yFile)
				test.ExpectNoError(t, err)

				ct, err := NewChronicleTemplate("foo", yFile)
				test.ExpectNil(t, ct)
				test.ExpectError(t, err)
			}
		})

	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic valid file", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "basic.yml")
			yFile, err := yaml.GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			ct, err := NewChronicleTemplate("foo.yml", yFile)
			test.ExpectNotNil(t, ct)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, ct.ID(), "simpleId")
			test.ExpectEqual(t, ct.Description(), "simpleDescription")
			test.ExpectNotSet(t, ct.Inherit())
			test.ExpectTrue(t, strings.HasSuffix(ct.Filename(), "foo.yml"))
			test.ExpectEqual(t, len(ct.content), 0)
			test.ExpectEqual(t, len(ct.presets), 0)
		})

		t.Run("file with content", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "valid.yml")
			yFile, err := yaml.GetYamlFile(fileToTest)

			test.ExpectNotNil(t, yFile)
			test.ExpectNoError(t, err)

			ct, err := NewChronicleTemplate("valid.yml", yFile)
			test.ExpectNotNil(t, ct)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, ct.Inherit(), "otherId")

			// content
			test.ExpectEqual(t, len(ct.content), 3)
			ce0, exists := ct.GetContent("c0")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, ce0.Type(), "textCell")
			test.ExpectEqual(t, ce0.ExampleValue(), "c0example")

			ce1, exists := ct.GetContent("c1")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, ce1.Type(), "societyId")
			test.ExpectEqual(t, ce1.ExampleValue(), "c1example")
			_, exists = ct.GetContent("p0")
			test.ExpectFalse(t, exists)

			// presets
			test.ExpectEqual(t, len(ct.presets), 3)
			p0, exists := ct.presets.Get("p0")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, p0.Y1, 10.0)
			p1, exists := ct.presets.Get("p1")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, p1.Y1, 11.0)
			_, exists = ct.presets.Get("c0")
			test.ExpectFalse(t, exists)
		})

	})

	t.Run("", func(t *testing.T) {
	})
}

func TestGetContentIDs(t *testing.T) {
	ct := getCTfromYamlFile(t, "valid.yml")

	idList := ct.GetContentIDs(true)

	test.ExpectEqual(t, len(idList), 3)

	// as we do not yet have aliases, number of elements should be identical
	test.ExpectEqual(t, len(ct.GetContentIDs(false)), len(ct.GetContentIDs(true)))

	// check that all elements returned by list actually exist in the content list
	for _, entry := range idList {
		_, exists := ct.GetContent(entry)
		test.ExpectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	test.ExpectEqual(t, idList[0], "c0")
	test.ExpectEqual(t, idList[1], "c1")
	test.ExpectEqual(t, idList[2], "c2")
}

func TestInheritFrom(t *testing.T) {
	// TODO split and move to ContentStore and PresetStore

	t.Run("inherit from valid", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromValid.yml")

		err := ctTo.InheritFrom(ctFrom)
		test.ExpectNoError(t, err)

		test.ExpectEqual(t, len(ctTo.presets.GetIDs()), 2)
		p0, exists := ctTo.presets.Get("p0")
		test.ExpectTrue(t, exists)
		test.ExpectEqual(t, p0.Font, "base")
		p1, exists := ctTo.presets.Get("p1")
		test.ExpectTrue(t, exists)
		test.ExpectEqual(t, p1.Font, "inherited")

		test.ExpectEqual(t, len(ctTo.GetContentIDs(false)), 2)
		ce0, exists := ctTo.GetContent("c0")
		test.ExpectTrue(t, exists)
		tc0 := ce0.(ContentTextCell)
		test.ExpectEqual(t, tc0.Font, "base")
		ce1, exists := ctTo.GetContent("c1")
		test.ExpectTrue(t, exists)
		tc1 := ce1.(ContentTextCell)
		test.ExpectEqual(t, tc1.Font, "inherited")
	})

	t.Run("inherit from duplicate content", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromDuplicateContent.yml")

		err := ctTo.InheritFrom(ctFrom)
		test.ExpectError(t, err)
	})

}
