package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
)

// ReadCsvFile reads the csv file from the provided location.
func ReadCsvFile(filename string) (records [][]string, err error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading csv file '%v': %v", filename, err)
	}

	// TODO check number of contains commas vs semicolons to determine separator

	r := csv.NewReader(bytes.NewReader(fileData))
	r.Comma = detectSeparator(fileData)
	r.Comment = '#'

	records, err = r.ReadAll()
	if err != nil {
		return nil, err
	}

	// ensure that we need less range checks
	alignRecordLength(&records)

	return records, nil
}

// detectSeparator detects whether a slice of bytes contains more commas or
// more semicolons. In case of a tie, semicolons win.
func detectSeparator(content []byte) (separator rune) {
	runes := bytes.Runes(content)
	var commas, semicolons int
	for _, r := range runes {
		switch r {
		case ',':
			commas++
		case ';':
			semicolons++
		}
	}

	if commas > semicolons {
		return ','
	}
	return ';'
}

// alignRecordLength takes a two-layered string array as input and ensures
// that each included array has the same length.
func alignRecordLength(records *[][]string) {
	var max int = 0
	for _, record := range *records {
		if len(record) > max {
			max = len(record)
		}
	}

	var empty string
	for idx, record := range *records {
		for len(record) < max {
			record = append(record, empty)
		}
		(*records)[idx] = record
	}
}

// WriteTemplateToCsvFile takes a chronicle template and creates a CSV file that can be used
// as input for the "batch fill" command
func (ct *ChronicleTemplate) WriteTemplateToCsvFile(filename string, as ArgStore, separator rune) (err error) {
	const numPlayers = 7

	// TODO if no file extension is added, add ".csv" automatically

	if separator != ';' && separator != ',' {
		return fmt.Errorf("Unsupported separator provided; only ';' and ',' are currently supported")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	records := [][]string{
		{"#ID", ct.ID()},
		{"#Description", ct.Description()},
		{"#"},
		{"#Players"}, // will be filled below with labels
	}
	for idx := 1; idx <= numPlayers; idx++ {
		outerIdx := len(records) - 1
		records[outerIdx] = append(records[outerIdx], fmt.Sprintf("Player %d", idx))
	}

	for _, contentID := range ct.GetContentIDs(false) {
		// entry should be large enough for id column + 7 players
		entry := make([]string, numPlayers+1)

		entry[0] = contentID

		// check if some value was provided on the cmd line that should be filled in everywhere
		if val, exists := as[contentID]; exists {
			for colIdx := 1; colIdx <= numPlayers; colIdx++ {
				entry[colIdx] = val
			}
		}

		records = append(records, entry)
	}

	csvw := csv.NewWriter(file)
	csvw.Comma = separator
	csvw.UseCRLF = false // TODO do we need to adapt this based on the OS?

	for _, record := range records {
		err = csvw.Write(record)
		if err != nil {
			return err
		}
	}
	csvw.Flush()

	return nil
}

// GetFillInformationFromCsvFile reads a csv file and returns a list of ArgStores that
// contain the required arguments to fill out a chronicle.
func GetFillInformationFromCsvFile(filename string) (argStores []ArgStore, err error) {
	records, err := ReadCsvFile(filename)
	if err != nil {
		return nil, err
	}

	argStores = make([]ArgStore, 0)

	if len(records) == 0 {
		return argStores, nil
	}

	numPlayers := len(records[0]) - 1

	for idx := 1; idx <= numPlayers; idx++ {
		as := make(ArgStore, len(records))

		for _, record := range records {
			key := record[0]
			value := record[idx]
			if _, exists := as[key]; exists {
				return nil, fmt.Errorf("File '%v' contains multiple lines for content ID '%v'", filename, key)
			}

			// only store if there is an actual value
			if IsSet(value) {
				if !IsSet(key) {
					return nil, fmt.Errorf("CSV Line has content value '%v', but is missing content ID in first column", value)
				}
				as[key] = value
			}
		}

		// only add if we have at least one entry here
		if len(as) >= 1 {
			argStores = append(argStores, as)
		}
	}

	return argStores, nil
}
