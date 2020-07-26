package main

import (
	"path/filepath"
	"strings"
	"testing"
)

var (
	chronicleTemplateTestDir string
)

func init() {
	SetIsTestEnvironment(true)
	chronicleTemplateTestDir = filepath.Join(GetExecutableDir(), "testdata", "ChronicleTemplate")
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
			c0, exists := ct.GetContent("c0")
			expectTrue(t, exists)
			expectEqual(t, c0.Type(), "someType")
			expectEqual(t, c0.X1(), 1.0)
			c1, exists := ct.GetContent("c1")
			expectTrue(t, exists)
			expectEqual(t, c1.Type(), "someType")
			expectEqual(t, c1.X1(), 2.0)
			_, exists = ct.GetContent("p0")
			expectFalse(t, exists)

			// presets
			expectEqual(t, len(ct.presets), 3)
			p0, exists := ct.GetPreset("p0")
			expectTrue(t, exists)
			expectEqual(t, p0.Y1, 10.0)
			p1, exists := ct.GetPreset("p1")
			expectTrue(t, exists)
			expectEqual(t, p1.Y1, 11.0)
			_, exists = ct.GetPreset("c0")
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

	t.Run("inherit from valid", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromValid.yml")

		err := ctTo.InheritFrom(ctFrom)
		expectNoError(t, err)

		expectEqual(t, len(ctTo.GetPresetIDs()), 2)
		p0, exists := ctTo.GetPreset("p0")
		expectTrue(t, exists)
		expectEqual(t, p0.Font, "base")
		p1, exists := ctTo.GetPreset("p1")
		expectTrue(t, exists)
		expectEqual(t, p1.Font, "inherited")

		expectEqual(t, len(ctTo.GetContentIDs(false)), 2)
		c0, exists := ctTo.GetContent("c0")
		expectTrue(t, exists)
		expectEqual(t, c0.Font(), "base")
		c1, exists := ctTo.GetContent("c1")
		expectTrue(t, exists)
		expectEqual(t, c1.Font(), "inherited")
	})

	t.Run("inherit from duplicate content", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromDuplicateContent.yml")

		err := ctTo.InheritFrom(ctFrom)
		expectError(t, err)
	})

}

func TestResolveContent(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.ResolveContent()
		expectNoError(t, err)

		c1, _ := ct.GetContent("c1")
		expectEqual(t, c1.X1(), 1.0)
		expectNotSet(t, c1.X2())

		c2, _ := ct.GetContent("c2")
		expectEqual(t, c2.X1(), 2.0)
		expectEqual(t, c2.X2(), 1.0)

		c3, _ := ct.GetContent("c3")
		expectEqual(t, c3.X2(), 23.0)
		expectEqual(t, c3.Y1(), 2.0)
		expectEqual(t, c3.Y2(), 3.0)
		expectEqual(t, c3.XPivot(), 4.0)

		c4, _ := ct.GetContent("c4")
		expectEqual(t, c4.X1(), 1.0)
		expectEqual(t, c4.X2(), 1.0)
		expectEqual(t, c4.Y1(), 2.0)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existing dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml")
			err := ct.ResolveContent()
			expectError(t, err)
		})

		t.Run("contradicting values", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.ResolveContent()
			expectError(t, err)
		})
	})
}
