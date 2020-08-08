package main

import (
	"os"
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

func TestWriteToCsvFile(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(csvTestDir, "templates"))
	test.ExpectNoError(t, err)
	test.ExpectNotNil(t, ts)

	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	t.Run("errors", func(t *testing.T) {
		t.Run("error during csv writing", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "unsupportedSeparator.csv")

			err = ct.WriteToCsvFile(outfile, '.', as)
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "basic.csv")

			// write template to csv
			err = ct.WriteToCsvFile(outfile, ';', as)
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := utils.ReadFileToLines(outfile)
			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(lines), 8)
			test.ExpectStringContains(t, lines[0], "ID")
			test.ExpectStringContains(t, lines[0], "template") // ID of the template
			test.ExpectEqual(t, lines[3], "#Players;Player 1;Player 2;Player 3;Player 4;Player 5;Player 6;Player 7")
			test.ExpectEqual(t, lines[6], "player;;;;;;;")
			test.ExpectEqual(t, lines[7], "societyid;;;;;;;")
		})

		t.Run("fill with example values", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := ArgStoreFromTemplateExamples(ct)

			outfile := filepath.Join(outputDir, "argStoreExamples.csv")

			// write template to csv
			err = ct.WriteToCsvFile(outfile, ';', as)
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := utils.ReadFileToLines(outfile)
			test.ExpectNoError(t, err)
			test.ExpectStringContains(t, lines[5], "noExample;;;")
			test.ExpectStringContains(t, lines[6], "player;Bob;Bob;")
			test.ExpectStringContains(t, lines[7], "societyid;12345-678;12345-678;")
		})

		t.Run("fill with user-provided values", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			as.Set("player", "Jack")
			as.Set("noExample", "test")

			outfile := filepath.Join(outputDir, "argStoreUserProvided.csv")

			// write template to csv
			err = ct.WriteToCsvFile(outfile, ';', as)
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := utils.ReadFileToLines(outfile)
			test.ExpectNoError(t, err)
			test.ExpectStringContains(t, lines[5], "noExample;test;test;")
			test.ExpectStringContains(t, lines[6], "player;Jack;Jack;")
			test.ExpectStringContains(t, lines[7], "societyid;;;")
		})
	})
}

func writeTemplateToFileAndReadBackIn(t *testing.T, ct *ChronicleTemplate, as *ArgStore, separator rune) (argStores []*ArgStore) {
	test.ExpectNotNil(t, ct)

	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	outfile := filepath.Join(outputDir, "test.csv")

	// write template to csv
	err := ct.WriteToCsvFile(outfile, separator, as)
	test.ExpectNoError(t, err)

	// read csv back in
	argStores, err = GetArgStoresFromCsvFile(outfile)
	test.ExpectNotNil(t, argStores)
	test.ExpectNoError(t, err)

	return argStores
}

func TestCreateAndReadCsvFile(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(csvTestDir, "templates"))
	test.ExpectNoError(t, err)
	test.ExpectNotNil(t, ts)

	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	ct, err := ts.GetTemplate("template")
	test.ExpectNoError(t, err)
	test.ExpectNotNil(t, ct)

	// built up test data
	inputArgStores := make([]*ArgStore, 0)
	inputArgStores = append(inputArgStores, NewArgStore(&ArgStoreInit{})) // empty argStore
	inputArgStores = append(inputArgStores, ArgStoreFromTemplateExamples(ct))
	userProvidedArgStore := NewArgStore(&ArgStoreInit{})
	userProvidedArgStore.Set("player", "Jack")
	userProvidedArgStore.Set("noExample", "test")
	inputArgStores = append(inputArgStores, userProvidedArgStore)

	// begin tests
	for _, inputArgStore := range inputArgStores {
		for _, separator := range []rune{';', ','} {
			resultArgStores := writeTemplateToFileAndReadBackIn(t, ct, inputArgStore, separator)
			test.ExpectNotNil(t, resultArgStores)

			// only case where no result can be empty is if input was completely empty
			test.ExpectTrue(t, len(resultArgStores) > 0 || inputArgStore.NumEntries() == 0)

			// compare complete content
			for _, resultArgStore := range resultArgStores {
				test.ExpectEqual(t, resultArgStore.NumEntries(), inputArgStore.NumEntries())

				for _, key := range resultArgStore.GetKeys() {
					resultValue, _ := resultArgStore.Get(key)
					inputValue, _ := inputArgStore.Get(key)
					test.ExpectEqual(t, resultValue, inputValue)
				}
			}
		}

	}
}
