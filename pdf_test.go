package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	pdfTestDir string
)

func init() {
	SetIsTestEnvironment(true)
	pdfTestDir = filepath.Join(GetExecutableDir(), "testdata", "pdf")
}

func Test_NewPdf_NonExistantFile(t *testing.T) {
	fileToTest := filepath.Join(pdfTestDir, "nonExistantFile.pdf")
	pdf, err := NewPdf(fileToTest)

	expectNil(t, pdf)
	expectError(t, err)
}

func Test_NewPdf_ProvideDirectory(t *testing.T) {
	pdf, err := NewPdf(pdfTestDir)

	expectNil(t, pdf)
	expectError(t, err)
}

func Test_NewPdf_ValidFileOnePage(t *testing.T) {
	fileToTest := filepath.Join(pdfTestDir, "OnePage.pdf")
	pdf, err := NewPdf(fileToTest)

	expectNotNil(t, pdf)
	expectNoError(t, err)
	expectEqual(t, pdf.numPages, 1)
}

func Test_NewPdf_ValidFileFourPages(t *testing.T) {
	fileToTest := filepath.Join(pdfTestDir, "FourPages.pdf")
	pdf, err := NewPdf(fileToTest)

	expectNotNil(t, pdf)
	expectNoError(t, err)
	expectEqual(t, pdf.numPages, 4)
}

func Test_ExtractPage_InvalidOutputDirectory(t *testing.T) {
	inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
	inPdf, err := NewPdf(inputFile)
	expectNotNil(t, inPdf)
	expectNoError(t, err)

	nonExistantDir := filepath.Join(pdfTestDir, "nonExistantDir")

	extractedPdf, err := inPdf.ExtractPage(1, nonExistantDir)
	expectNil(t, extractedPdf)
	expectError(t, err)
}

func Test_ExtractPage_InvalidPage(t *testing.T) {
	inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
	inPdf, err := NewPdf(inputFile)
	expectNotNil(t, inPdf)
	expectNoError(t, err)

	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	for _, pageIndex := range []int{-5, 0, 5} {
		extractedPdf, err := inPdf.ExtractPage(pageIndex, workDir)
		expectNil(t, extractedPdf)
		expectError(t, err)
	}
}

func Test_ExtractPage_Valid(t *testing.T) {
	inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
	inPdf, err := NewPdf(inputFile)
	expectNotNil(t, inPdf)
	expectNoError(t, err)

	for _, pageIndex := range []int{-4, -1, 1, 4} {
		workDir := GetTempDir()
		defer os.RemoveAll(workDir)

		extractedPdf, err := inPdf.ExtractPage(pageIndex, workDir)
		expectNotNil(t, extractedPdf)
		expectNoError(t, err)

		// debug output
		if err != nil {
			files, err := ioutil.ReadDir(workDir)
			if err != nil {
				t.Logf("Cannot read dir '%v' to analyze issue: %v", workDir, err)
				continue
			}
			t.Logf("Files in directory %v:", workDir)
			for _, file := range files {
				t.Logf("- %v\n", file.Name())
			}
		}

	}
}
