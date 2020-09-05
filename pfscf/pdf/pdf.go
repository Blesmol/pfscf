package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// File is a wraper for a PDF file
type File struct {
	filename string
	numPages int
}

// NewFile creates a new Pdf object.
func NewFile(filename string) (p *File, err error) {
	if exists, err := utils.IsFile(filename); !exists {
		return nil, err
	}

	p = new(File)
	p.filename = filename
	p.numPages, err = pdfcpuapi.PageCountFile(p.filename)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Filename returns the filename of the given PDF
func (pf *File) Filename() (filename string) {
	return pf.filename
}

// AllowsPageExtraction checks whether the permissions contained in the
// PDF file allow to extract pages from it
func (pf *File) AllowsPageExtraction() bool {
	return pf.GetPermissionBit(4) == true && pf.GetPermissionBit(11) == true
}

// ExtractPage extracts a single page from the input file and stores
// it under (but not necessarily in) the given output directory.
// Provided page number can also be negative, then page is searched from the back.
func (pf *File) ExtractPage(pageNumber int, outDir string) (extractedPage *File, err error) {
	isDir, err := utils.IsDir(outDir)
	if !isDir || err != nil {
		return nil, fmt.Errorf("Error extracting page from file %v: %v", pf.filename, err)
	}

	// check PDF permissions
	if !pf.AllowsPageExtraction() {
		return nil, fmt.Errorf("File %v does not allow page extraction", pf.filename)
	}

	// Function accepts negative page numbers, thus calculate real page number
	var realPageNumber int
	if pageNumber < 0 {
		realPageNumber = pf.numPages + /*negative*/ pageNumber + 1 // as -1 is the last page
	} else {
		realPageNumber = pageNumber
	}
	if realPageNumber <= 0 || realPageNumber > pf.numPages {
		return nil, fmt.Errorf("Page number %v is out of bounds for file %v", realPageNumber, pf.filename)
	}

	realPageNumberStr := strconv.Itoa(realPageNumber)
	err = pdfcpuapi.ExtractPagesFile(pf.filename, outDir, []string{realPageNumberStr}, nil)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, pf.filename, err)
	}

	extractedPdf, err := NewFile(getPdfPageExtractionFilename(outDir, pf.filename, realPageNumberStr))
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, pf.filename, err)
	}

	return extractedPdf, nil
}

// GetDimensionsInPoints returns the width and height of the first page in
// a given PDF file
func (pf *File) GetDimensionsInPoints() (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(pf.filename)
	utils.AssertNoError(err)
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
	localFilename := strings.Join([]string{inFileBase, "_page_", page, ".pdf"}, "")
	return filepath.Join(dirname, localFilename)
}

// GetPermissionBit checks whether the given permission bit
// is set for the given PDF file
func (pf *File) GetPermissionBit(bit int) (bitValue bool) {
	perms, err := pdfcpuapi.ListPermissionsFile(pf.filename, nil)
	utils.AssertNoError(err)
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
func (pf *File) StampIt(stampFile string, outFile string) (err error) {
	onTop := true     // stamps go on top, watermarks do not
	updateWM := false // should the new watermark be added or an existing one updated?

	wm, err := pdfcpuapi.PDFWatermark(stampFile, "rot:0, sc:1", onTop, updateWM)
	if err != nil {
		return err
	}

	err = pdfcpuapi.AddWatermarksFile(pf.filename, outFile, nil, wm, nil)
	if err != nil {
		return err
	}

	return nil
}

// Fill is the main function used to fill a PDF file.
func (pf *File) Fill(argStore *args.Store, ct *template.ChronicleTemplate, outfile string) (err error) {
	// prepare temporary working dir
	workDir := utils.GetTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	extractedPage, err := pf.ExtractPage(-1, workDir)
	if err != nil {
		return err
	}
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := stamp.NewStamp(width, height)

	if cfg.Global.DrawCellBorder {
		stamp.SetCellBorder(true)
	}

	// add content to stamp
	ct.GenerateOutput(stamp, argStore)

	if cfg.Global.DrawGrid {
		stamp.CreateMeasurementCoordinates(5, 1)
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
