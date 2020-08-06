package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	util "github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	pdfTestDir string
)

func init() {
	util.SetIsTestEnvironment(true)
	pdfTestDir = filepath.Join(util.GetExecutableDir(), "testdata", "pdf")
}

func TestNewPdf(t *testing.T) {
	t.Run("non-existant file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "nonExistantFile.pdf")
		pdf, err := NewPdf(fileToTest)

		expectNil(t, pdf)
		expectError(t, err)
	})

	t.Run("directory instead of file", func(t *testing.T) {
		pdf, err := NewPdf(pdfTestDir)

		expectNil(t, pdf)
		expectError(t, err)
	})

	t.Run("valid one-page file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "OnePage.pdf")
		pdf, err := NewPdf(fileToTest)

		expectNotNil(t, pdf)
		expectNoError(t, err)
		expectEqual(t, pdf.numPages, 1)
	})

	t.Run("valid four-pages file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "FourPages.pdf")
		pdf, err := NewPdf(fileToTest)

		expectNotNil(t, pdf)
		expectNoError(t, err)
		expectEqual(t, pdf.numPages, 4)
	})
}

func TestExtractPage(t *testing.T) {
	t.Run("invalid output directory", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		expectNotNil(t, inPdf)
		expectNoError(t, err)

		nonExistantDir := filepath.Join(pdfTestDir, "nonExistantDir")

		extractedPdf, err := inPdf.ExtractPage(1, nonExistantDir)
		expectNil(t, extractedPdf)
		expectError(t, err)
	})

	t.Run("invalid page number", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		expectNotNil(t, inPdf)
		expectNoError(t, err)

		workDir := util.GetTempDir()
		defer os.RemoveAll(workDir)

		for _, pageIndex := range []int{-5, 0, 5} {
			extractedPdf, err := inPdf.ExtractPage(pageIndex, workDir)
			expectNil(t, extractedPdf)
			expectError(t, err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		expectNotNil(t, inPdf)
		expectNoError(t, err)

		for _, pageIndex := range []int{-4, -1, 1, 4} {
			workDir := util.GetTempDir()
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
	})
}
