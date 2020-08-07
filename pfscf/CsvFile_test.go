package main

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	csvTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	csvTestDir = filepath.Join(utils.GetExecutableDir(), "testdata", "CsvFile")
}

func readFileToLines(t *testing.T, filename string) (lines []string, err error) {
	t.Helper()

	lines = make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineScanner := bufio.NewScanner(file)
	for lineScanner.Scan() {
		lines = append(lines, lineScanner.Text())
	}

	if lineScanner.Err() != nil {
		return nil, lineScanner.Err()
	}

	return lines, nil
}

func TestReadCsvFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		records, err := ReadCsvFile(filepath.Join(csvTestDir, "nonExisting.csv"))
		test.ExpectError(t, err)
		test.ExpectNil(t, records)
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			for _, filename := range []string{"validBasicSemicolon.csv", "validBasicComma.csv"} {
				t.Logf("Processing file '%v'", filename)
				records, err := ReadCsvFile(filepath.Join(csvTestDir, filename))
				test.ExpectNotNil(t, records)
				test.ExpectNoError(t, err)

				test.ExpectEqual(t, len(records), 3)
				test.ExpectEqual(t, len(records[0]), 5)

				test.ExpectEqual(t, records[1][0], "societyid")
				test.ExpectEqual(t, records[1][2], "1233-123")
				test.ExpectEqual(t, records[2][4], "Fire")
			}
		})

		t.Run("equal number of commas and semicolons", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "equalNumberOfCommasAndSemicolons.csv"))
			test.ExpectNotNil(t, records)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(records), 3)
			test.ExpectEqual(t, len(records[0]), 5)

			test.ExpectEqual(t, records[1][0], "societyid,")
			test.ExpectEqual(t, records[2][3], ",Air")
		})

		t.Run("more commas than semicolons", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "moreCommas.csv"))
			test.ExpectNotNil(t, records)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(records), 3)
			test.ExpectEqual(t, len(records[0]), 5)

			test.ExpectEqual(t, records[1][0], "societyid")
			test.ExpectEqual(t, records[2][3], "Air;")
		})

		t.Run("more semicolons than commas", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "moreSemicolons.csv"))
			test.ExpectNotNil(t, records)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(records), 3)
			test.ExpectEqual(t, len(records[0]), 5)

			test.ExpectEqual(t, records[1][0], "societyid,")
			test.ExpectEqual(t, records[2][3], ",Air")
		})

	})
}

func TestAlignRecordLength(t *testing.T) {
	// build testdata
	records := make([][]string, 10)
	for idx, record := range records {
		record = make([]string, idx)
		for idx2 := range record {
			record[idx2] = "foo"
		}
		records[idx] = record
	}

	alignRecordLength(&records)

	for idx, record := range records {
		test.ExpectEqual(t, len(record), len(records)-1)
		for idx2, entry := range record {
			if idx2 < idx {
				test.ExpectEqual(t, entry, "foo")
			} else {
				test.ExpectNotSet(t, entry)
			}
		}
	}
}

func TestWriteTemplateToCsvFile(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(csvTestDir, "templates"))
	test.ExpectNoError(t, err)
	test.ExpectNotNil(t, ts)

	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	t.Run("errors", func(t *testing.T) {
		t.Run("Non-existing target dir", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "nonExisting", "basic.csv")
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			test.ExpectError(t, err)
		})

		t.Run("unsupported separator", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "unsupportedSeparator.csv")

			err = ct.WriteTemplateToCsvFile(outfile, as, '.')
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic with semicolon", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "basic_with_semicolon.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(lines), 8)
			test.ExpectStringContains(t, lines[0], "ID")
			test.ExpectStringContains(t, lines[0], "template") // ID of the template
			test.ExpectEqual(t, lines[3], "#Players;Player 1;Player 2;Player 3;Player 4;Player 5;Player 6;Player 7")
			test.ExpectEqual(t, lines[6], "player;;;;;;;")
			test.ExpectEqual(t, lines[7], "societyid;;;;;;;")
		})

		t.Run("basic with comma", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := NewArgStore(&ArgStoreInit{})
			outfile := filepath.Join(outputDir, "basic_with_comma.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ',')
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(lines), 8)
			test.ExpectStringContains(t, lines[0], "ID")
			test.ExpectStringContains(t, lines[0], "template") // ID of the template
			test.ExpectEqual(t, lines[3], "#Players,Player 1,Player 2,Player 3,Player 4,Player 5,Player 6,Player 7")
			test.ExpectEqual(t, lines[6], "player,,,,,,,")
			test.ExpectEqual(t, lines[7], "societyid,,,,,,,")
		})

		t.Run("fill with example values", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			test.ExpectNoError(t, err)
			test.ExpectNotNil(t, ct)

			as := ArgStoreFromTemplateExamples(ct)

			outfile := filepath.Join(outputDir, "argStoreExamples.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
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
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			test.ExpectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			test.ExpectNoError(t, err)
			test.ExpectStringContains(t, lines[5], "noExample;test;test;")
			test.ExpectStringContains(t, lines[6], "player;Jack;Jack;")
			test.ExpectStringContains(t, lines[7], "societyid;;;")
		})
	})
}

func TestGetFillInformationFromCsvFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existing file", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "nonExisting.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNil(t, as)
			test.ExpectError(t, err)
		})

		t.Run("content without ID", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "contentWithoutId.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNil(t, as)
			test.ExpectError(t, err)
		})

		t.Run("duplicate content id", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "duplicateContent.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNil(t, as)
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty file", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "emptyFile.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)
		})

		t.Run("basic file", func(t *testing.T) {
			for _, baseFilename := range []string{"validBasicSemicolon.csv", "validBasicComma.csv"} {
				t.Logf("Filename is '%v'", baseFilename)
				filename := filepath.Join(csvTestDir, baseFilename)
				argStores, err := GetFillInformationFromCsvFile(filename)
				test.ExpectNotNil(t, argStores)
				test.ExpectNoError(t, err)

				test.ExpectEqual(t, len(argStores), 4)

				for _, data := range []struct {
					argStore *ArgStore
					key      string
					expValue string
				}{
					{argStores[0], "player", "John"},
					{argStores[0], "societyid", "123456-789"},
					{argStores[0], "char", "Earth"},
					{argStores[3], "player", "Hanna"},
					{argStores[3], "societyid", "7435-432"},
					{argStores[3], "char", "Fire"},
				} {
					argEntry, exists := data.argStore.Get(data.key)
					test.ExpectTrue(t, exists)
					test.ExpectEqual(t, argEntry, data.expValue)
				}
			}
		})

		t.Run("empty lines and comment lines", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "emptyLines.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(argStores), 1)
			test.ExpectEqual(t, argStores[0].NumEntries(), 3)

			for _, data := range []struct {
				argStore *ArgStore
				key      string
				expValue string
			}{
				{argStores[0], "player", "John"},
				{argStores[0], "societyid", "123456-789"},
				{argStores[0], "char", "Earth"},
			} {
				argEntry, exists := data.argStore.Get(data.key)
				test.ExpectTrue(t, exists)
				test.ExpectEqual(t, argEntry, data.expValue)
			}
		})

		t.Run("file without players", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "noPlayers.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNotNil(t, as)
			test.ExpectNoError(t, err)
		})

		t.Run("file with missing values", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "validWithSomeMissingValues.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(argStores), 4)

			for _, data := range []struct {
				argStore   *ArgStore
				expEntries int
				key        string
			}{
				{argStores[0], 2, "societyid"},
				{argStores[0], 2, "char"},
				{argStores[1], 2, "player"},
				{argStores[1], 2, "char"},
				{argStores[2], 2, "player"},
				{argStores[2], 2, "societyid"},
			} {
				test.ExpectEqual(t, data.argStore.NumEntries(), data.expEntries)

				argEntry, exists := data.argStore.Get(data.key)
				test.ExpectTrue(t, exists)
				test.ExpectIsSet(t, argEntry)
			}
		})

		// currently this is only checked while stamping, so reading this in is currently not an error
		t.Run("invalid society id", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "invalidSocietyId.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			test.ExpectNotNil(t, as)
			test.ExpectNoError(t, err)
		})
	})
}

func writeTemplateToFileAndReadBackIn(t *testing.T, ct *ChronicleTemplate, as *ArgStore, separator rune) (argStores []*ArgStore) {
	test.ExpectNotNil(t, ct)

	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	outfile := filepath.Join(outputDir, "test.csv")

	// write template to csv
	err := ct.WriteTemplateToCsvFile(outfile, as, separator)
	test.ExpectNoError(t, err)

	// read csv back in
	argStores, err = GetFillInformationFromCsvFile(outfile)
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
