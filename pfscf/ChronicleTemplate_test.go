package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Blesmol/pfscf/pfscf/utils"
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
	yFile, err := GetYamlFile(fileToTest)
	expectNotNil(t, yFile)
	expectNoError(t, err)

	ct, err = NewChronicleTemplate(filename, yFile)
	expectNotNil(t, ct)
	expectNoError(t, err)
	return ct
}

func TestNewChronicleTemplate(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("no filename", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "basic.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			ct, err := NewChronicleTemplate("", yFile)
			expectNil(t, ct)
			expectError(t, err)
		})

		t.Run("nil file", func(t *testing.T) {
			ct, err := NewChronicleTemplate("file.yml", nil)
			expectNil(t, ct)
			expectError(t, err)
		})

		t.Run("empty fields", func(t *testing.T) {
			filenames := []string{"emptyId.yml", "emptyDescription.yml"}

			for _, filename := range filenames {
				fileToTest := filepath.Join(chronicleTemplateTestDir, filename)
				yFile, err := GetYamlFile(fileToTest)

				expectNotNil(t, yFile)
				expectNoError(t, err)

				ct, err := NewChronicleTemplate("foo", yFile)
				expectNil(t, ct)
				expectError(t, err)
			}
		})

	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic valid file", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "basic.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			ct, err := NewChronicleTemplate("foo.yml", yFile)
			expectNotNil(t, ct)
			expectNoError(t, err)

			expectEqual(t, ct.ID(), "simpleId")
			expectEqual(t, ct.Description(), "simpleDescription")
			expectNotSet(t, ct.Inherit())
			expectTrue(t, strings.HasSuffix(ct.Filename(), "foo.yml"))
			expectEqual(t, len(ct.content), 0)
			expectEqual(t, len(ct.presets), 0)
		})

		t.Run("file with content", func(t *testing.T) {
			fileToTest := filepath.Join(chronicleTemplateTestDir, "valid.yml")
			yFile, err := GetYamlFile(fileToTest)

			expectNotNil(t, yFile)
			expectNoError(t, err)

			ct, err := NewChronicleTemplate("valid.yml", yFile)
			expectNotNil(t, ct)
			expectNoError(t, err)

			expectEqual(t, ct.Inherit(), "otherId")

			// content
			expectEqual(t, len(ct.content), 3)
			ce0, exists := ct.GetContent("c0")
			expectTrue(t, exists)
			expectEqual(t, ce0.Type(), "textCell")
			expectEqual(t, ce0.ExampleValue(), "c0example")

			ce1, exists := ct.GetContent("c1")
			expectTrue(t, exists)
			expectEqual(t, ce1.Type(), "societyId")
			expectEqual(t, ce1.ExampleValue(), "c1example")
			_, exists = ct.GetContent("p0")
			expectFalse(t, exists)

			// presets
			expectEqual(t, len(ct.presets), 3)
			p0, exists := ct.presets.Get("p0")
			expectTrue(t, exists)
			expectEqual(t, p0.Y1, 10.0)
			p1, exists := ct.presets.Get("p1")
			expectTrue(t, exists)
			expectEqual(t, p1.Y1, 11.0)
			_, exists = ct.presets.Get("c0")
			expectFalse(t, exists)
		})

	})

	t.Run("", func(t *testing.T) {
	})
}

func TestGetContentIDs(t *testing.T) {
	ct := getCTfromYamlFile(t, "valid.yml")

	idList := ct.GetContentIDs(true)

	expectEqual(t, len(idList), 3)

	// as we do not yet have aliases, number of elements should be identical
	expectEqual(t, len(ct.GetContentIDs(false)), len(ct.GetContentIDs(true)))

	// check that all elements returned by list actually exist in the content list
	for _, entry := range idList {
		_, exists := ct.GetContent(entry)
		expectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	expectEqual(t, idList[0], "c0")
	expectEqual(t, idList[1], "c1")
	expectEqual(t, idList[2], "c2")
}

func TestInheritFrom(t *testing.T) {
	// TODO split and move to ContentStore and PresetStore

	t.Run("inherit from valid", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromValid.yml")

		err := ctTo.InheritFrom(ctFrom)
		expectNoError(t, err)

		expectEqual(t, len(ctTo.presets.GetIDs()), 2)
		p0, exists := ctTo.presets.Get("p0")
		expectTrue(t, exists)
		expectEqual(t, p0.Font, "base")
		p1, exists := ctTo.presets.Get("p1")
		expectTrue(t, exists)
		expectEqual(t, p1.Font, "inherited")

		expectEqual(t, len(ctTo.GetContentIDs(false)), 2)
		ce0, exists := ctTo.GetContent("c0")
		expectTrue(t, exists)
		tc0 := ce0.(ContentTextCell)
		expectEqual(t, tc0.Font, "base")
		ce1, exists := ctTo.GetContent("c1")
		expectTrue(t, exists)
		tc1 := ce1.(ContentTextCell)
		expectEqual(t, tc1.Font, "inherited")
	})

	t.Run("inherit from duplicate content", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromDuplicateContent.yml")

		err := ctTo.InheritFrom(ctFrom)
		expectError(t, err)
	})

}
