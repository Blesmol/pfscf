package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const input = "scenario1_nr.pdf"
const watermark = "watermark.pdf"
const output = "test/chronicle1.pdf"

func main() {
	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	pdf := NewPdf(input)

	if !pdf.AllowsPageExtraction() {
		fmt.Printf("Error: File %v does not allow page extraction, exiting", input)
		os.Exit(1)
	}

	// extract chronicle page from pdf
	extractedPage := pdf.ExtractPage(-1, workDir)
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := NewStamp(width, height)

	// demo text
	stamp.AddText(433, 107, "Grand Archive", "Helvetica", 8)
	stamp.AddCellText(227, 730, 305, 716, "05.06.2020", "Helvetica", 14)

	//stamp.CreateMeasurementCoordinates(25, 5)

	// write stamp
	stampFile := filepath.Join(workDir, "stamp.pdf")
	stamp.WriteToFile(stampFile)

	// add watermark/stamp to page
	extractedPage.StampIt(stampFile, output)

	// Configuration test run
	config := GetGlobalConfig()
	fmt.Printf("Config:\n%+v\n", *config)
	fmt.Printf("Content: %+v", *config.Content)
}
