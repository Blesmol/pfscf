package pdf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
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
func (f *File) Filename() (filename string) {
	return f.filename
}

// AllowsPageExtraction checks whether the permissions contained in the
// PDF file allow to extract pages from it
func (f *File) AllowsPageExtraction() bool {
	return f.GetPermissionBit(4) == true && f.GetPermissionBit(11) == true
}

// ExtractPage extracts a single page from the input file and stores
// it under (but not necessarily in) the given output directory.
// Provided page number can also be negative, then page is searched from the back.
func (f *File) ExtractPage(pageNumber int, outDir string) (extractedPage *File, err error) {
	isDir, err := utils.IsDir(outDir)
	if !isDir || err != nil {
		return nil, fmt.Errorf("Error extracting page from file %v: %v", f.filename, err)
	}

	// check PDF permissions
	if !f.AllowsPageExtraction() {
		return nil, fmt.Errorf("File %v does not allow page extraction", f.filename)
	}

	// Function accepts negative page numbers, thus calculate real page number
	var realPageNumber int
	if pageNumber < 0 {
		realPageNumber = f.numPages + /*negative*/ pageNumber + 1 // as -1 is the last page
	} else {
		realPageNumber = pageNumber
	}
	if realPageNumber <= 0 || realPageNumber > f.numPages {
		return nil, fmt.Errorf("Page number %v is out of bounds for file %v", realPageNumber, f.filename)
	}

	// Create PDF context
	ctx, err := pdfcpuapi.ReadContextFile(f.filename)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, f.filename, err)
	}

	// Extract requested page
	ctxNew, err := ctx.ExtractPage(realPageNumber)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, f.filename, err)
	}

	// Write context to file
	outFile := filepath.Join(outDir, "extracted.pdf")
	if err := api.WriteContextFile(ctxNew, outFile); err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, f.filename, err)
	}

	extractedPdf, err := NewFile(outFile)
	if err != nil {
		return nil, fmt.Errorf("Error extracting page %v from file %v: %v", realPageNumber, f.filename, err)
	}

	return extractedPdf, nil
}

// GetDimensionsInPoints returns the width and height of the first page in
// a given PDF file
func (f *File) GetDimensionsInPoints() (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(f.filename)
	utils.AssertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return dim[0].Width, dim[0].Height
}

func isBitSet(bitfield *uint16, pos int) bool {
	val := *bitfield & (1 << pos)
	return (val > 0)
}

// GetPermissionBit checks whether the given permission bit
// is set for the given PDF file
func (f *File) GetPermissionBit(bit int) (bitValue bool) {
	perms, err := pdfcpuapi.GetPermissionsFile(f.filename, nil)
	utils.AssertNoError(err)
	if perms == nil {
		// no permissions return => assume true/allow as default
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

	// TODO find out which exact permissions are required and need to be checked.
	// As first conservatie approach assume that if permissions are present, then nothing is allowed
	return false
}

// StampIt stamps the given PDF file with the given stamp
func (f *File) StampIt(stampFile string, outFile string) (err error) {
	onTop := true     // stamps go on top, watermarks do not
	updateWM := false // should the new watermark be added or an existing one updated?

	wm, err := pdfcpuapi.PDFWatermark(stampFile, "rot:0, sc:1", onTop, updateWM)
	if err != nil {
		return err
	}

	err = pdfcpuapi.AddWatermarksFile(f.filename, outFile, nil, wm, nil)
	if err != nil {
		return err
	}

	return nil
}

// Fill is the main function used to fill a PDF file.
func (f *File) Fill(argStore *args.Store, ct *template.Chronicle, outfile string) (err error) {
	// prepare temporary working dir
	workDir := utils.GetTempDir()
	defer os.RemoveAll(workDir)

	// extract chronicle page from pdf
	extractedPage, err := f.ExtractPage(-1, workDir)
	if err != nil {
		return err
	}
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := stamp.NewStamp(width, height, cfg.Global.OffsetX, cfg.Global.OffsetY)

	if cfg.Global.DrawCellBorder {
		stamp.SetCellBorder(true)
	}

	// add content to stamp
	if err = ct.GenerateOutput(stamp, argStore); err != nil {
		return err
	}

	if cfg.Global.DrawCanvasGrid != "" {
		if err = stamp.DrawCanvasGrid(cfg.Global.DrawCanvasGrid); err != nil {
			return fmt.Errorf("Error drawing canvas grid: %v", err)
		}
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
