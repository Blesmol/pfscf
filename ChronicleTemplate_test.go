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
	SetIsTestEnvironment()
	chronicleTemplateTestDir = filepath.Join(GetExecutableDir(), "testdata", "ChronicleTemplate")
}

func Test_NewChronicleTemplate_NoFilename(t *testing.T) {
	fileToTest := filepath.Join(chronicleTemplateTestDir, "basic.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	ct, err := NewChronicleTemplate("", yFile)
	expectNil(t, ct)
	expectError(t, err)
}

func Test_NewChronicleTemplate_NilFile(t *testing.T) {
	ct, err := NewChronicleTemplate("file.yml", nil)
	expectNil(t, ct)
	expectError(t, err)
}

func Test_NewChronicleTemplate_EmptyFields(t *testing.T) {
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
}

func Test_NewChronicleTemplate_BasicValidFile(t *testing.T) {
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
}

func Test_NewChronicleTemplate_FileWithContent(t *testing.T) {
	fileToTest := filepath.Join(chronicleTemplateTestDir, "valid.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	ct, err := NewChronicleTemplate("valid.yml", yFile)
	expectNotNil(t, ct)
	expectNoError(t, err)

	expectNotSet(t, ct.Inherit()) // will sooner or later be filled...

	// content
	expectEqual(t, len(ct.content), 2)
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
	expectEqual(t, len(ct.presets), 2)
	p0, exists := ct.GetPreset("p0")
	expectTrue(t, exists)
	expectEqual(t, p0.Y1(), 10.0)
	p1, exists := ct.GetPreset("p1")
	expectTrue(t, exists)
	expectEqual(t, p1.Y1(), 11.0)
	_, exists = ct.GetPreset("c0")
	expectFalse(t, exists)
}
