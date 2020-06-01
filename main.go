package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	name, err := ioutil.TempDir("", "pfsct")
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

func getPdfDimensionsInPoints(filename string) (x int, y int) {
	// TODO change return type to whatever is accepted by gofpdf
	dim, err := pdfcpuapi.PageDimsFile(filename)
	assertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return int(dim[0].Width), int(dim[0].Height)
}

func main() {
	// prepare temporary working dir
	workDir := getTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	chroniclePage := getLastPage(input)
	pdfcpuapi.ExtractPagesFile(input, workDir, []string{chroniclePage}, nil)
	extractedPage := getPdfPageExtractionFilename(workDir, chroniclePage)

	getPdfDimensionsInPoints(extractedPage)

	// add demo watermark do page
	onTop := true
	wm, err := pdfcpu.ParsePDFWatermarkDetails(watermark, "rot:0, sc:1", onTop)
	assertNoError(err)
	err = pdfcpuapi.AddWatermarksFile(extractedPage, output, nil, wm, nil)
	assertNoError(err)

}
