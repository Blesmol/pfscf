package main

import (
	"github.com/jung-kurt/gofpdf"
)

// Stamp is a wraper for a PDF page
type Stamp struct {
	pdf *gofpdf.Fpdf
}

// NewStamp creates a new Stamp object.
func NewStamp(dimX float64, dimY float64) (s *Stamp) {
	s = new(Stamp)

	s.pdf = gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size:    gofpdf.SizeType{Wd: dimX, Ht: dimY},
	})
	s.pdf.AddPage() // 0,0 is top-left. To change use AddPageFormat() instead

	return s
}

// WriteToFile writes the contect of the Stamp object into a PDF file.
// The Stamp object should not be used anymore after that.
func (s *Stamp) WriteToFile(filename string) (err error) {
	return s.pdf.OutputFileAndClose(filename)
}

// Pdf returns the included gofpdf.Fpdf object.
// Function should be finally removed
func (s *Stamp) Pdf() (pdf *gofpdf.Fpdf) {
	return s.pdf
}
