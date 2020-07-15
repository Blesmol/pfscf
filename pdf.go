package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

// Pdf is a wraper for a PDF file
type Pdf struct {
	filename string
	numPages int
}

// NewPdf creates a new Pdf object.
func NewPdf(filename string) (p *Pdf, err error) {
	if exists, err := IsFile(filename); !exists {
		return nil, err
	}
	p = new(Pdf)
	p.filename = filename
	p.numPages, err = pdfcpuapi.PageCountFile(p.filename)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Filename returns the filename of the given PDF
func (pdf *Pdf) Filename() (filename string) {
	return pdf.filename
}

// AllowsPageExtraction checks whether the permissions contained in the
// PDF file allow to extract pages from it
func (pdf *Pdf) AllowsPageExtraction() bool {
	return pdf.GetPermissionBit(4) == true && pdf.GetPermissionBit(11) == true
}

// ExtractPage extracts a single page from the input file and stores
// it under (but not necessarily in) the given output directory.
// Provided page number can also be negative, then page is searched from the back.
func (pdf *Pdf) ExtractPage(pageNumber int, outDir string) (extractedPage *Pdf, err error) {
	isDir, err := IsDir(outDir)
	if !isDir || err != nil {
		return nil, fmt.Errorf("Error extracting page from file %v: %v", pdf.filename, err)
	}

	// check PDF permissions
	if !pdf.AllowsPageExtraction() {
		return nil, fmt.Errorf("File %v does not allow page extraction", pdf.filename)
	}

	// Function accepts negative page numbers, thus calculate real page number
	var realPageNumber int
	if pageNumber < 0 {
		realPageNumber = pdf.numPages + /*negative*/ pageNumber + 1 // as -1 is the last page
	} else {
		realPageNumber = pageNumber
	}
	if realPageNumber <= 0 || realPageNumber > pdf.numPages {
		return nil, fmt.Errorf("Page number %v is out of bounds for file %v", realPageNumber, pdf.filename)
	}

	realPageNumberStr := strconv.Itoa(realPageNumber)
	err = pdfcpuapi.ExtractPagesFile(pdf.filename, outDir, []string{realPageNumberStr}, nil)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, pdf.filename, err)
	}

	extractedPdf, err := NewPdf(getPdfPageExtractionFilename(outDir, pdf.filename, realPageNumberStr))
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, pdf.filename, err)
	}

	return extractedPdf, nil
}

// GetDimensionsInPoints returns the width and height of the first page in
// a given PDF file
func (pdf *Pdf) GetDimensionsInPoints() (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(pdf.filename)
	AssertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return dim[0].Width, dim[0].Height
}

// getPdfPageExtractionFilename returns the path and filename of the target
// file if a single page was extracted.
func getPdfPageExtractionFilename(dirname, inFile, page string) (outFile string) {
	inFileWithoutDir := filepath.Base(inFile)
	inFileBase := strings.TrimSuffix(inFileWithoutDir, filepath.Ext(inFileWithoutDir))
	localFilename := strings.Join([]string{inFileBase, "_", page, ".pdf"}, "")
	return filepath.Join(dirname, localFilename)
}

// GetPermissionBit checks whether the given permission bit
// is set for the given PDF file
func (pdf *Pdf) GetPermissionBit(bit int) (bitValue bool) {
	perms, err := pdfcpuapi.ListPermissionsFile(pdf.filename, nil)
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

// StampIt stamps the given PDF file with the given stamp
func (pdf *Pdf) StampIt(stampFile string, outFile string) (err error) {
	onTop := true     // stamps go on top, watermarks do not
	updateWM := false // should the new watermark be added or an existing one updated?

	wm, err := pdfcpuapi.PDFWatermark(stampFile, "rot:0, sc:1", onTop, updateWM)
	if err != nil {
		return err
	}

	err = pdfcpuapi.AddWatermarksFile(pdf.filename, outFile, nil, wm, nil)
	if err != nil {
		return err
	}

	return nil
}

// Fill is the main function used to fill a PDF file.
func (pdf *Pdf) Fill(argStore ArgStore, ct *ChronicleTemplate, outfile string) (err error) {
	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	extractedPage, err := pdf.ExtractPage(-1, workDir)
	if err != nil {
		return err
	}
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := NewStamp(width, height)

	if drawCellBorder { // TODO use argument rather than global variable
		stamp.SetCellBorder(true)
	}

	// add content to stamp
	for key, value := range argStore {
		//fmt.Printf("Processing Key='%v', value='%v'\n", key, value)

		content, exists := ct.GetContent(key)
		if !exists {
			return fmt.Errorf("Found no content with key '%v'", key)
		}

		err := stamp.AddContent(content, &value)
		if err != nil {
			return err
		}
	}

	if drawGrid { // TODO use argument rather than global variable
		stamp.CreateMeasurementCoordinates(25, 5)
	}

	// write stamp
	stampFile := filepath.Join(workDir, "stamp.pdf")
	err = stamp.WriteToFile(stampFile)
	if err != nil {
		return err
	}

	// add watermark/stamp to page
	err = extractedPage.StampIt(stampFile, outfile)
	if err != nil {
		return err
	}

	return nil
}
