package csv

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

// WriteFile creates a CSV file with the provided 2-dimensional array as content.
func WriteFile(filename string, separator rune, data [][]string) (err error) {

	if separator != ';' && separator != ',' {
		return fmt.Errorf("Unsupported separator provided; only ';' and ',' are currently supported")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	csvw := csv.NewWriter(file)
	csvw.Comma = separator
	csvw.UseCRLF = false // TODO do we need to adapt this based on the OS?

	for _, record := range data {
		err = csvw.Write(record)
		if err != nil {
			return err
		}
	}
	csvw.Flush()

	return nil
}
