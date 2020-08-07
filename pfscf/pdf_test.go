package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	pdfTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	pdfTestDir = filepath.Join(utils.GetExecutableDir(), "testdata", "pdf")
}

func TestNewPdf(t *testing.T) {
	t.Run("non-existant file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "nonExistantFile.pdf")
		pdf, err := NewPdf(fileToTest)

		test.ExpectNil(t, pdf)
		test.ExpectError(t, err)
	})

	t.Run("directory instead of file", func(t *testing.T) {
		pdf, err := NewPdf(pdfTestDir)

		test.ExpectNil(t, pdf)
		test.ExpectError(t, err)
	})

	t.Run("valid one-page file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "OnePage.pdf")
		pdf, err := NewPdf(fileToTest)

		test.ExpectNotNil(t, pdf)
		test.ExpectNoError(t, err)
		test.ExpectEqual(t, pdf.numPages, 1)
	})

	t.Run("valid four-pages file", func(t *testing.T) {
		fileToTest := filepath.Join(pdfTestDir, "FourPages.pdf")
		pdf, err := NewPdf(fileToTest)

		test.ExpectNotNil(t, pdf)
		test.ExpectNoError(t, err)
		test.ExpectEqual(t, pdf.numPages, 4)
	})
}

func TestExtractPage(t *testing.T) {
	t.Run("invalid output directory", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		test.ExpectNotNil(t, inPdf)
		test.ExpectNoError(t, err)

		nonExistantDir := filepath.Join(pdfTestDir, "nonExistantDir")

		extractedPdf, err := inPdf.ExtractPage(1, nonExistantDir)
		test.ExpectNil(t, extractedPdf)
		test.ExpectError(t, err)
	})

	t.Run("invalid page number", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		test.ExpectNotNil(t, inPdf)
		test.ExpectNoError(t, err)

		workDir := utils.GetTempDir()
		defer os.RemoveAll(workDir)

		for _, pageIndex := range []int{-5, 0, 5} {
			extractedPdf, err := inPdf.ExtractPage(pageIndex, workDir)
			test.ExpectNil(t, extractedPdf)
			test.ExpectError(t, err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		inputFile := filepath.Join(pdfTestDir, "FourPages.pdf")
		inPdf, err := NewPdf(inputFile)
		test.ExpectNotNil(t, inPdf)
		test.ExpectNoError(t, err)

		for _, pageIndex := range []int{-4, -1, 1, 4} {
			workDir := utils.GetTempDir()
			defer os.RemoveAll(workDir)

			extractedPdf, err := inPdf.ExtractPage(pageIndex, workDir)
			test.ExpectNotNil(t, extractedPdf)
			test.ExpectNoError(t, err)

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
