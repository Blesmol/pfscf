package main

import (
	"fmt"
	"os"
	"path/filepath"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

const input = "scenario1_nr.pdf"
const watermark = "watermark.pdf"
const output = "test/chronicle1.pdf"

func getPdfDimensionsInPoints(filename string) (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(filename)
	AssertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return dim[0].Width, dim[0].Height
}

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
	if !DoesPdfFileAllowPageExtraction(input) {
		fmt.Printf("Error: File %v does not allow page extraction, exiting", input)
		os.Exit(1)
	}

	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	chroniclePage := GetPdfLastPageNumber(input)
	pdfcpuapi.ExtractPagesFile(input, workDir, []string{chroniclePage}, nil)
	extractedPage := GetPdfPageExtractionFilename(workDir, chroniclePage)

	width, height := getPdfDimensionsInPoints(extractedPage)

	// add demo watermark to page
	onTop := true
	stampFile := createPdfStampFile(workDir, width, height)
	wm, err := pdfcpu.ParsePDFWatermarkDetails(stampFile, "rot:0, sc:1", onTop)
	AssertNoError(err)
	err = pdfcpuapi.AddWatermarksFile(extractedPage, output, nil, wm, nil)
	AssertNoError(err)

}
