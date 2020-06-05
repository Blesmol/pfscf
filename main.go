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
	stamp.AddText(264, 730, "05.06.2020", "Helvetica", 14)
	stamp.AddText(433, 107, "Grand Archive", "Helvetica", 8)
	stampFile := filepath.Join(workDir, "stamp.pdf")
	stamp.WriteToFile(stampFile)

	// add watermark/stamp to page
	extractedPage.StampIt(stampFile, output)
}
