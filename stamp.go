package main

import (
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

// Stamp is a wraper for a PDF page
type Stamp struct {
	pdf  *gofpdf.Fpdf
	dimX float64
	dimY float64
}

// NewStamp creates a new Stamp object.
func NewStamp(dimX float64, dimY float64) (s *Stamp) {
	s = new(Stamp)

	s.dimX = dimX
	s.dimY = dimY

	s.pdf = gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size:    gofpdf.SizeType{Wd: dimX, Ht: dimY},
	})
	s.pdf.AddPage() // 0,0 is top-left. To change use AddPageFormat() instead

	return s
}

// AddText adds a portion of text at the specified coordinates
func (s *Stamp) AddText() {
	s.pdf.SetFont("Arial", "B", 16)
	s.pdf.Cell(40, 10, "Hello, world")
}

// WriteToFile writes the contect of the Stamp object into a PDF file.
// The Stamp object should not be used anymore after that.
func (s *Stamp) WriteToFile(filename string) {
	err := s.pdf.OutputFileAndClose(filename)
	AssertNoError(err)
}

// Pdf returns the included gofpdf.Fpdf object.
// Function should be finally removed
func (s *Stamp) Pdf() (pdf *gofpdf.Fpdf) {
	return s.pdf
}

// CreateMeasurementCoordinates overlays the stamp with a set of lines
func (s *Stamp) CreateMeasurementCoordinates(gap float64) {
	s.pdf.SetFont("Arial", "B", 6)

	for curX := float64(0); curX < s.dimX; curX += gap {
		s.pdf.Line(curX, 0, curX, s.dimY)
		s.pdf.Text(curX, 8, strconv.Itoa(int(curX)))
	}

	for curY := float64(0); curY < s.dimY; curY += gap {
		s.pdf.Line(0, curY, s.dimY, curY)
		s.pdf.Text(2, curY-1, strconv.Itoa(int(curY)))
	}

}
