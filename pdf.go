package main

import (
	"path/filepath"
	"strconv"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

// DoesPdfFileAllowPageExtraction checks whether the permissions contained in the
// PDF file allow to extract pages from it
func DoesPdfFileAllowPageExtraction(filename string) bool {
	return GetPdfPermissionBit(filename, 4) == true && GetPdfPermissionBit(filename, 11) == true
}

// GetPdfDimensionsInPoints returns the width and height of the first page in
// a given PDF file
func GetPdfDimensionsInPoints(filename string) (width float64, height float64) {
	dim, err := pdfcpuapi.PageDimsFile(filename)
	AssertNoError(err)
	if len(dim) != 1 {
		panic(dim)
	}
	return dim[0].Width, dim[0].Height
}

// GetPdfLastPageNumber returns the number of the last page
// in the given PDF file as string.
func GetPdfLastPageNumber(file string) (page string) {
	numPages, err := pdfcpuapi.PageCountFile(file)
	AssertNoError(err)
	return strconv.Itoa(numPages)
}

// GetPdfPageExtractionFilename returns the path and filename of the target
// file if a single page was extracted.
func GetPdfPageExtractionFilename(dir string, page string) (filename string) {
	localFilename := strings.Join([]string{"page_", page, ".pdf"}, "")
	return filepath.Join(dir, localFilename)
}

// GetPdfPermissionBit checks whether the given permission bit
// is set for the given PDF file
func GetPdfPermissionBit(filename string, bit int) (bitValue bool) {
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
