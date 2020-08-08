package main

import (
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

func TestWriteFile(t *testing.T) {
	outputDir := utils.GetTempDir()
	defer os.RemoveAll(outputDir)

	data := [][]string{
		{"foo1"},
		{"bar;1", "bar,2"},
	}

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existing target dir", func(t *testing.T) {
			outfile := filepath.Join(outputDir, "nonExisting", "basic.csv")
			err := CsvWriteFile(outfile, ';', data)
			test.ExpectError(t, err)
		})

		t.Run("invalid separator", func(t *testing.T) {
			outfile := filepath.Join(outputDir, "invalid_separator.csv")
			err := CsvWriteFile(outfile, 'g', data)
			test.ExpectError(t, err, "Unsupported separator")
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("separators", func(t *testing.T) {
			t.Run("semicolon", func(t *testing.T) {
				outfile := filepath.Join(outputDir, "separator_semicolon.csv")
				err := CsvWriteFile(outfile, ';', data)
				test.ExpectNoError(t, err)

				// read csv back in for content check
				lines, err := utils.ReadFileToLines(outfile)
				test.ExpectNoError(t, err)
				test.ExpectEqual(t, len(lines), 2)
				test.ExpectStringContains(t, lines[0], "foo1")
				test.ExpectEqual(t, lines[1], "\"bar;1\";bar,2")
			})
			t.Run("comma", func(t *testing.T) {
				outfile := filepath.Join(outputDir, "separator_comma.csv")
				err := CsvWriteFile(outfile, ',', data)
				test.ExpectNoError(t, err)

				// read csv back in for content check
				lines, err := utils.ReadFileToLines(outfile)
				test.ExpectNoError(t, err)
				test.ExpectEqual(t, len(lines), 2)
				test.ExpectStringContains(t, lines[0], "foo1")
				test.ExpectEqual(t, lines[1], "bar;1,\"bar,2\"")
			})
		})
		t.Run("empty data", func(t *testing.T) {
			outfile := filepath.Join(outputDir, "empty.csv")
			err := CsvWriteFile(outfile, ',', [][]string{})
			test.ExpectNoError(t, err)

			// read csv back in for content check
			lines, err := utils.ReadFileToLines(outfile)
			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(lines), 0)
		})
	})
}
