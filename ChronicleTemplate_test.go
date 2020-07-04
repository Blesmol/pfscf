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
			expectEqual(t, c0.Type(), "textCell")
			expectEqual(t, c0.X1(), 1.0)
			c1, exists := ct.GetContent("c1")
			expectTrue(t, exists)
			expectEqual(t, c1.Type(), "textCell")
			expectEqual(t, c1.X1(), 2.0)
			_, exists = ct.GetContent("p0")
			expectFalse(t, exists)

			// presets
			expectEqual(t, len(ct.presets), 3)
			p0, exists := ct.GetPreset("p0")
			expectTrue(t, exists)
			expectEqual(t, p0.Y1(), 10.0)
			p1, exists := ct.GetPreset("p1")
			expectTrue(t, exists)
			expectEqual(t, p1.Y1(), 11.0)
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

func TestGetPresetIDs(t *testing.T) {
	ct := getCTfromYamlFile(t, "valid.yml")

	idList := ct.GetPresetIDs()

	expectEqual(t, len(idList), 3)

	// check that all elements returned by list actually exist
	for _, entry := range idList {
		_, exists := ct.GetPreset(entry)
		expectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	expectEqual(t, idList[0], "p0")
	expectEqual(t, idList[1], "p1")
	expectEqual(t, idList[2], "p2")
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
		expectEqual(t, p0.Type(), "base")
		p1, exists := ctTo.GetPreset("p1")
		expectTrue(t, exists)
		expectEqual(t, p1.Type(), "inherited")

		expectEqual(t, len(ctTo.GetContentIDs(false)), 2)
		c0, exists := ctTo.GetContent("c0")
		expectTrue(t, exists)
		expectEqual(t, c0.Type(), "base")
		c1, exists := ctTo.GetContent("c1")
		expectTrue(t, exists)
		expectEqual(t, c1.Type(), "inherited")
	})

	t.Run("inherit from duplicate content", func(t *testing.T) {
		ctTo := getCTfromYamlFile(t, "inheritTo.yml")
		ctFrom := getCTfromYamlFile(t, "inheritFromDuplicateContent.yml")

		err := ctTo.InheritFrom(ctFrom)
		expectError(t, err)
	})

}

func TestResolvePresets(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.ResolvePresets()
		expectNoError(t, err)

		p2, _ := ct.GetPreset("p2")
		expectEqual(t, p2.X1(), 1.0)
		expectEqual(t, p2.X2(), 23.0)
		expectEqual(t, p2.Y1(), 2.0)
		expectNotSet(t, p2.Y2())
		expectEqual(t, p2.XPivot(), 1.0)

		p3, _ := ct.GetPreset("p3")
		expectEqual(t, p3.X1(), 1.0)
		expectEqual(t, p3.X2(), 23.0)
		expectNotSet(t, p3.Y1())
		expectEqual(t, p3.Y2(), 3.0)
		expectEqual(t, p3.XPivot(), 1.0)

		p4, _ := ct.GetPreset("p4")
		expectEqual(t, p4.X1(), 1.0)
		expectEqual(t, p4.X2(), 23.0)
		expectEqual(t, p4.Y1(), 2.0)
		expectEqual(t, p4.Y2(), 3.0)
		expectEqual(t, p4.XPivot(), 4.0)

		p5, _ := ct.GetPreset("p5")
		expectEqual(t, p5.X1(), 1.0)
		expectEqual(t, p5.X2(), 23.0)
		expectEqual(t, p5.Y1(), 2.0)
		expectNotSet(t, p5.Y2())
		expectEqual(t, p5.XPivot(), 1.0)
	})

	t.Run("cyclic presets", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveCyclic.yml")
		err := ct.ResolvePresets()
		expectError(t, err)
	})

	t.Run("self-dependency", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveSelf.yml")
		err := ct.ResolvePresets()
		expectError(t, err)
	})

	t.Run("non-existing dependency", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveNonExisting.yml")
		err := ct.ResolvePresets()
		expectError(t, err)
	})

	t.Run("contradicting values", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveContradicting.yml")
		err := ct.ResolvePresets()
		expectError(t, err)
	})
}
