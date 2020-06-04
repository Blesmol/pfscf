package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const input = "scenario1_nr.pdf"
const watermark = "watermark.pdf"
const output = "test/chronicle1.pdf"

func createPdfStampFile(targetDir string, width float64, height float64) (filename string) {
	filename = filepath.Join(targetDir, "stamp.pdf")

	stamp := NewStamp(width, height)

	pdf := stamp.Pdf()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")

	stamp.CreateMeasurementCoordinates(float64(25))

	err := stamp.WriteToFile(filename)
	AssertNoError(err)

	return filename
}

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

	stampFile := createPdfStampFile(workDir, width, height)

	// add watermark/stamp to page
	pdf.Stamp(stampFile, output)
}
