package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
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
func (p *Pdf) Filename() (filename string) {
	return p.filename
}

// AllowsPageExtraction checks whether the permissions contained in the
// PDF file allow to extract pages from it
func (p *Pdf) AllowsPageExtraction() bool {
	return p.GetPermissionBit(4) == true && p.GetPermissionBit(11) == true
}

// ExtractPage extracts a single page from the input file and stores
// it under (but not necessarily in) the given output directory.
// Provided page number can also be negative, then page is searched from the back.
func (p *Pdf) ExtractPage(pageNumber int, outDir string) (extractedPage *Pdf, err error) {
	isDir, err := IsDir(outDir)
	if !isDir || err != nil {
		return nil, fmt.Errorf("Error extracting page from file %v: %w", p.filename, err)
	}

	// check PDF permissions
	if !p.AllowsPageExtraction() {
		return nil, fmt.Errorf("File %v does not allow page extraction", p.filename)
	}

	// Function accepts negative page numbers, thus calculate real page number
	var realPageNumber int
	if pageNumber < 0 {
		realPageNumber = p.numPages + /*negative*/ pageNumber + 1 // as -1 is the last page
	} else {
		realPageNumber = pageNumber
	}
	if realPageNumber <= 0 || realPageNumber > p.numPages {
		return nil, fmt.Errorf("Page number %v is out of bounds for file %v", realPageNumber, p.filename)
	}

	realPageNumberStr := strconv.Itoa(realPageNumber)
	err = pdfcpuapi.ExtractPagesFile(p.filename, outDir, []string{realPageNumberStr}, nil)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %w", realPageNumber, p.filename, err)
	}

	extractedPdf, err := NewPdf(getPdfPageExtractionFilename(outDir, p.filename, realPageNumberStr))
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %w", realPageNumber, p.filename, err)
	}

	return extractedPdf, nil
}

// GetDimensionsInPoints returns the width and height of the first page in
// a given PDF file
func (p *Pdf) GetDimensionsInPoints() (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(p.filename)
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
func (p *Pdf) GetPermissionBit(bit int) (bitValue bool) {
	perms, err := pdfcpuapi.ListPermissionsFile(p.filename, nil)
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
func (p *Pdf) StampIt(stampFile string, outFile string) {
	onTop := true // stamps go on top, watermarks do not
	wm, err := pdfcpu.ParsePDFWatermarkDetails(stampFile, "rot:0, sc:1", onTop)
	AssertNoError(err)
	err = pdfcpuapi.AddWatermarksFile(p.filename, outFile, nil, wm, nil)
	AssertNoError(err)

}
