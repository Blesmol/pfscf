package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

const input = "scenario1_nr.pdf"
const watermark = "watermark.pdf"
const output = "test/chronicle1.pdf"

func getLastPage(file string) (page string) {
	numPages, err := pdfcpuapi.PageCountFile(file)
	AssertNoError(err)
	return strconv.Itoa(numPages)
}

func getPdfPermissionBit(filename string, bit int) (bitValue bool) {
	perms, err := pdfcpuapi.ListPermissionsFile(filename, nil)
	AssertNoError(err)
	if len(perms) <= 1 {
		// no permissions return => assume true/allow as default
		// TODO should check whether text is "Full access"
		return true
	}

	// If permissions are set, then first entry of perms should have
	// the following format in pdfcpu:
	//   permission bits: 101100110100
	// Else text would be "Full access"
	//
	// - Bit 1 is on the right side.
	// - "1" is true, "0" is false.
	// - "True" means the permission is granted
	//
	// - Bit  3: print(rev2), print quality(rev>=3)
	// - Bit  4: modify other than controlled by bits 6,9,11
	// - Bit  5: extract(rev2), extract other than controlled by bit 10(rev>=3)
	// - Bit  6: add or modify annotations
	// - Bit  9: fill in form fields(rev>=3)
	// - Bit 10: extract(rev>=3)
	// - Bit 11: modify(rev>=3)
	// - Bit 12: print high-level(rev>=3)

	/*
		fmt.Printf("Permissions:\n")
		for _, val := range perms {
			fmt.Printf("- %v\n", val)
		}
	*/

	// TODO add proper permission check. As first conservatie approach assume
	// that if permissions are present, then nothing is allowed
	return false
}

func doesPdfFileAllowPageExtraction(filename string) bool {
	return getPdfPermissionBit(filename, 4) == true && getPdfPermissionBit(filename, 11) == true
}

func getPdfPageExtractionFilename(dir string, page string) (filename string) {
	localFilename := strings.Join([]string{"page_", page, ".pdf"}, "")
	return filepath.Join(dir, localFilename)
}

func getPdfDimensionsInPoints(filename string) (x float64, y float64) {
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
	if !doesPdfFileAllowPageExtraction(input) {
		fmt.Printf("Error: File %v does not allow page extraction, exiting", input)
		os.Exit(1)
	}

	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	chroniclePage := getLastPage(input)
	pdfcpuapi.ExtractPagesFile(input, workDir, []string{chroniclePage}, nil)
	extractedPage := getPdfPageExtractionFilename(workDir, chroniclePage)

	width, height := getPdfDimensionsInPoints(extractedPage)

	// add demo watermark to page
	onTop := true
	stampFile := createPdfStampFile(workDir, width, height)
	wm, err := pdfcpu.ParsePDFWatermarkDetails(stampFile, "rot:0, sc:1", onTop)
	AssertNoError(err)
	err = pdfcpuapi.AddWatermarksFile(extractedPage, output, nil, wm, nil)
	AssertNoError(err)

}
