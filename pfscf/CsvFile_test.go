package main

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

var (
	csvTestDir string
)

func init() {
	SetIsTestEnvironment(true)
	csvTestDir = filepath.Join(GetExecutableDir(), "testdata", "CsvFile")
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
		expectError(t, err)
		expectNil(t, records)
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic", func(t *testing.T) {
			for _, filename := range []string{"validBasicSemicolon.csv", "validBasicComma.csv"} {
				t.Logf("Processing file '%v'", filename)
				records, err := ReadCsvFile(filepath.Join(csvTestDir, filename))
				expectNotNil(t, records)
				expectNoError(t, err)

				expectEqual(t, len(records), 3)
				expectEqual(t, len(records[0]), 5)

				expectEqual(t, records[1][0], "societyid")
				expectEqual(t, records[1][2], "1233-123")
				expectEqual(t, records[2][4], "Fire")
			}
		})

		t.Run("equal number of commas and semicolons", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "equalNumberOfCommasAndSemicolons.csv"))
			expectNotNil(t, records)
			expectNoError(t, err)

			expectEqual(t, len(records), 3)
			expectEqual(t, len(records[0]), 5)

			expectEqual(t, records[1][0], "societyid,")
			expectEqual(t, records[2][3], ",Air")
		})

		t.Run("more commas than semicolons", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "moreCommas.csv"))
			expectNotNil(t, records)
			expectNoError(t, err)

			expectEqual(t, len(records), 3)
			expectEqual(t, len(records[0]), 5)

			expectEqual(t, records[1][0], "societyid")
			expectEqual(t, records[2][3], "Air;")
		})

		t.Run("more semicolons than commas", func(t *testing.T) {
			records, err := ReadCsvFile(filepath.Join(csvTestDir, "moreSemicolons.csv"))
			expectNotNil(t, records)
			expectNoError(t, err)

			expectEqual(t, len(records), 3)
			expectEqual(t, len(records[0]), 5)

			expectEqual(t, records[1][0], "societyid,")
			expectEqual(t, records[2][3], ",Air")
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
		expectEqual(t, len(record), len(records)-1)
		for idx2, entry := range record {
			if idx2 < idx {
				expectEqual(t, entry, "foo")
			} else {
				expectNotSet(t, entry)
			}
		}
	}
}

func TestWriteTemplateToCsvFile(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(csvTestDir, "templates"))
	expectNoError(t, err)
	expectNotNil(t, ts)

	outputDir := GetTempDir()
	defer os.RemoveAll(outputDir)

	t.Run("errors", func(t *testing.T) {
		t.Run("Non-existing target dir", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := make(ArgStore, 0)
			outfile := filepath.Join(outputDir, "nonExisting", "basic.csv")
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			expectError(t, err)
		})

		t.Run("unsupported separator", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := make(ArgStore, 0)
			outfile := filepath.Join(outputDir, "unsupportedSeparator.csv")

			err = ct.WriteTemplateToCsvFile(outfile, as, '.')
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("basic with semicolon", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := make(ArgStore, 0)
			outfile := filepath.Join(outputDir, "basic_with_semicolon.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			expectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			expectNoError(t, err)
			expectEqual(t, len(lines), 8)
			expectStringContains(t, lines[0], "ID")
			expectStringContains(t, lines[0], "template") // ID of the template
			expectEqual(t, lines[3], "#Players;Player 1;Player 2;Player 3;Player 4;Player 5;Player 6;Player 7")
			expectEqual(t, lines[6], "player;;;;;;;")
			expectEqual(t, lines[7], "societyid;;;;;;;")
		})

		t.Run("basic with comma", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := make(ArgStore, 0)
			outfile := filepath.Join(outputDir, "basic_with_comma.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ',')
			expectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			expectNoError(t, err)
			expectEqual(t, len(lines), 8)
			expectStringContains(t, lines[0], "ID")
			expectStringContains(t, lines[0], "template") // ID of the template
			expectEqual(t, lines[3], "#Players,Player 1,Player 2,Player 3,Player 4,Player 5,Player 6,Player 7")
			expectEqual(t, lines[6], "player,,,,,,,")
			expectEqual(t, lines[7], "societyid,,,,,,,")
		})

		t.Run("fill with example values", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := ArgStoreFromTemplateExamples(ct)

			outfile := filepath.Join(outputDir, "argStoreExamples.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			expectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			expectNoError(t, err)
			expectStringContains(t, lines[5], "noExample;;;")
			expectStringContains(t, lines[6], "player;Bob;Bob;")
			expectStringContains(t, lines[7], "societyid;12345-678;12345-678;")
		})

		t.Run("fill with user-provided values", func(t *testing.T) {
			ct, err := ts.GetTemplate("template")
			expectNoError(t, err)
			expectNotNil(t, ct)

			as := make(ArgStore, 0)
			as["player"] = "Jack"
			as["noExample"] = "test"

			outfile := filepath.Join(outputDir, "argStoreUserProvided.csv")

			// write template to csv
			err = ct.WriteTemplateToCsvFile(outfile, as, ';')
			expectNoError(t, err)

			// try to read the same csv back in (as text) for some basic content check
			lines, err := readFileToLines(t, outfile)
			expectNoError(t, err)
			expectStringContains(t, lines[5], "noExample;test;test;")
			expectStringContains(t, lines[6], "player;Jack;Jack;")
			expectStringContains(t, lines[7], "societyid;;;")
		})
	})
}

func TestGetFillInformationFromCsvFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("non-existing file", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "nonExisting.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			expectNil(t, as)
			expectError(t, err)
		})

		t.Run("content without ID", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "contentWithoutId.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			expectNil(t, as)
			expectError(t, err)
		})

		t.Run("duplicate content id", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "duplicateContent.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			expectNil(t, as)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty file", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "emptyFile.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			expectNotNil(t, argStores)
			expectNoError(t, err)
		})

		t.Run("basic file", func(t *testing.T) {
			for _, baseFilename := range []string{"validBasicSemicolon.csv", "validBasicComma.csv"} {
				t.Logf("Filename is '%v'", baseFilename)
				filename := filepath.Join(csvTestDir, baseFilename)
				argStores, err := GetFillInformationFromCsvFile(filename)
				expectNotNil(t, argStores)
				expectNoError(t, err)

				expectEqual(t, len(argStores), 4)

				expectEqual(t, argStores[0]["player"], "John")
				expectEqual(t, argStores[0]["societyid"], "123456-789")
				expectEqual(t, argStores[0]["char"], "Earth")

				expectEqual(t, argStores[3]["player"], "Hanna")
				expectEqual(t, argStores[3]["societyid"], "7435-432")
				expectEqual(t, argStores[3]["char"], "Fire")
			}
		})

		t.Run("empty lines and comment lines", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "emptyLines.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			expectNotNil(t, argStores)
			expectNoError(t, err)

			expectEqual(t, len(argStores), 1)
			expectEqual(t, len(argStores[0]), 3)

			expectEqual(t, argStores[0]["player"], "John")
			expectEqual(t, argStores[0]["societyid"], "123456-789")
			expectEqual(t, argStores[0]["char"], "Earth")
		})

		t.Run("file without players", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "noPlayers.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			expectNotNil(t, as)
			expectNoError(t, err)
		})

		t.Run("file with missing values", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "validWithSomeMissingValues.csv")
			argStores, err := GetFillInformationFromCsvFile(filename)
			expectNotNil(t, argStores)
			expectNoError(t, err)

			expectEqual(t, len(argStores), 4)

			expectEqual(t, len(argStores[0]), 2)
			expectIsSet(t, argStores[0]["societyid"])
			expectIsSet(t, argStores[0]["char"])

			expectEqual(t, len(argStores[1]), 2)
			expectIsSet(t, argStores[1]["player"])
			expectIsSet(t, argStores[1]["char"])

			expectEqual(t, len(argStores[2]), 2)
			expectIsSet(t, argStores[2]["player"])
			expectIsSet(t, argStores[2]["societyid"])
		})

		// currently this is only checked while stamping, so reading this in is currently not an error
		t.Run("invalid society id", func(t *testing.T) {
			filename := filepath.Join(csvTestDir, "invalidSocietyId.csv")
			as, err := GetFillInformationFromCsvFile(filename)
			expectNotNil(t, as)
			expectNoError(t, err)
		})
	})
}

func writeTemplateToFileAndReadBackIn(t *testing.T, ct *ChronicleTemplate, as ArgStore, separator rune) (argStores []ArgStore) {
	expectNotNil(t, ct)

	outputDir := GetTempDir()
	defer os.RemoveAll(outputDir)

	outfile := filepath.Join(outputDir, "test.csv")

	// write template to csv
	err := ct.WriteTemplateToCsvFile(outfile, as, separator)
	expectNoError(t, err)

	// read csv back in
	argStores, err = GetFillInformationFromCsvFile(outfile)
	expectNotNil(t, argStores)
	expectNoError(t, err)

	return argStores
}

func TestCreateAndReadCsvFile(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(csvTestDir, "templates"))
	expectNoError(t, err)
	expectNotNil(t, ts)

	outputDir := GetTempDir()
	defer os.RemoveAll(outputDir)

	ct, err := ts.GetTemplate("template")
	expectNoError(t, err)
	expectNotNil(t, ct)

	// built up test data
	inputArgStores := make([]ArgStore, 0)
	inputArgStores = append(inputArgStores, make(ArgStore, 0)) // empty argStore
	inputArgStores = append(inputArgStores, ArgStoreFromTemplateExamples(ct))
	userProvidedArgStore := make(ArgStore, 0)
	userProvidedArgStore["player"] = "Jack"
	userProvidedArgStore["noExample"] = "test"
	inputArgStores = append(inputArgStores, userProvidedArgStore)

	// begin tests
	for _, inputArgStore := range inputArgStores {
		for _, separator := range []rune{';', ','} {
			resultArgStores := writeTemplateToFileAndReadBackIn(t, ct, inputArgStore, separator)
			expectNotNil(t, resultArgStores)

			// only case where no result can be empty is if input was completely empty
			expectTrue(t, len(resultArgStores) > 0 || len(inputArgStore) == 0)

			// compare complete content
			for _, resultArgStore := range resultArgStores {
				expectEqual(t, len(resultArgStore), len(inputArgStore))

				for key, value := range resultArgStore {
					expectEqual(t, value, inputArgStore[key])
				}
			}
		}

	}
}
