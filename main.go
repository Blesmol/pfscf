package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

const input = "scenario1.pdf"
const watermark = "watermark.pdf"
const output = "test/chronicle1.pdf"

func assert(cond bool, err error) {
	if cond == false {
		fmt.Printf("Error is %v\n", err)
		panic(err)
	}
}

func assertNoError(err error) {
	assert(err == nil, err)
}

func getTempDir() (name string) {
	// TODO Wait for watermarking issue to be fixed on side of pdfcpu
	// https://github.com/pdfcpu/pdfcpu/issues/195
	// Watermarking with pdfcpu currently does not work on Windows
	// when absolute paths are used.
	// So temporarily create the working dir as subdir of the local directory
	name, err := ioutil.TempDir(".", "pfsct-")
	assertNoError(err)
	return name
}

func getLastPage(file string) (page string) {
	numPages, err := pdfcpuapi.PageCountFile(file)
	assertNoError(err)
	return strconv.Itoa(numPages)
}

func getPdfPageExtractionFilename(dir string, page string) (filename string) {
	localFilename := strings.Join([]string{"page_", page, ".pdf"}, "")
	return filepath.Join(dir, localFilename)
}

func getPdfDimensionsInPoints(filename string) (x float64, y float64) {
	dim, err := pdfcpuapi.PageDimsFile(filename)
	assertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return dim[0].Width, dim[0].Height
}

func createPdfStampFile(targetDir string, width float64, height float64) (filename string) {
	filename = filepath.Join(targetDir, "stamp.pdf")

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size:    gofpdf.SizeType{Wd: width, Ht: height},
	})

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose(filename)
	assertNoError(err)
	return filename
}

func main() {
	// prepare temporary working dir
	workDir := getTempDir()
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
	assertNoError(err)
	err = pdfcpuapi.AddWatermarksFile(extractedPage, output, nil, wm, nil)
	assertNoError(err)

}
