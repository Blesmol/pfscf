package args

import (
	"path/filepath"
	"testing"

	"github.com/Blesmol/pfscf/pfscf/csv"
	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	argStoreTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	argStoreTestDir = filepath.Join(utils.GetExecutableDir(), "testdata")
}

func getCsvRecords(t *testing.T, filename string) (records [][]string) {
	t.Helper()

	records, err := csv.ReadCsvFile(filename)
	if err != nil {
		t.Errorf("Problems reading CSV file '%v': %v", filename, err)
	}
	return records
}

func TestGetArgStoresFromCsvFile(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("content without ID", func(t *testing.T) {
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "contentWithoutId.csv"))
			as, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNil(t, as)
			test.ExpectError(t, err)
		})

		t.Run("duplicate content id", func(t *testing.T) {
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "duplicateContent.csv"))
			as, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNil(t, as)
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty file", func(t *testing.T) {
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "emptyFile.csv"))
			argStores, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)
		})

		t.Run("basic file", func(t *testing.T) {
			for _, baseFilename := range []string{"validBasicSemicolon.csv", "validBasicComma.csv"} {
				t.Logf("Filename is '%v'", baseFilename)
				records := getCsvRecords(t, filepath.Join(argStoreTestDir, baseFilename))
				argStores, err := GetArgStoresFromCsvRecords(records)
				test.ExpectNotNil(t, argStores)
				test.ExpectNoError(t, err)

				test.ExpectEqual(t, len(argStores), 4)

				for _, data := range []struct {
					argStore *Store
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
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "emptyLines.csv"))
			argStores, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(argStores), 1)
			test.ExpectEqual(t, argStores[0].numEntries(), 3)

			for _, data := range []struct {
				argStore *Store
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
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "noPlayers.csv"))
			as, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNotNil(t, as)
			test.ExpectNoError(t, err)
		})

		t.Run("file with missing values", func(t *testing.T) {
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "validWithSomeMissingValues.csv"))
			argStores, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNotNil(t, argStores)
			test.ExpectNoError(t, err)

			test.ExpectEqual(t, len(argStores), 4)

			for _, data := range []struct {
				argStore   *Store
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
				test.ExpectEqual(t, data.argStore.numEntries(), data.expEntries)

				argEntry, exists := data.argStore.Get(data.key)
				test.ExpectTrue(t, exists)
				test.ExpectIsSet(t, argEntry)
			}
		})

		// currently this is only checked while stamping, so reading this in is currently not an error
		t.Run("invalid society id", func(t *testing.T) {
			records := getCsvRecords(t, filepath.Join(argStoreTestDir, "invalidSocietyId.csv"))
			as, err := GetArgStoresFromCsvRecords(records)
			test.ExpectNotNil(t, as)
			test.ExpectNoError(t, err)
		})
	})
}
