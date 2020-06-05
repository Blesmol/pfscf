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
	s.pdf.SetMargins(0, 0, 0)
	s.pdf.SetAutoPageBreak(false, 0)
	s.pdf.AddPage() // 0,0 is top-left. To change use AddPageFormat() instead

	return s
}

// AddText adds a portion of text at the specified coordinates
func (s *Stamp) AddText(x, y float64, text string, fontName string, fontSize float64) {
	s.pdf.SetFont(fontName, "", fontSize)
	s.pdf.Text(x, y, text)
}

func getXYWH(x1, y1, x2, y2 float64) (x, y, w, h float64) {
	if x1 < x2 {
		x = x1
		w = x2 - x1
	} else {
		x = x2
		w = x1 - x2
	}
	if y1 < y2 {
		y = y1
		h = y2 - y1
	} else {
		y = y2
		h = y1 - y2
	}
	return
}

// AddCellText is a better version of AddText()
func (s *Stamp) AddCellText(x1, y1, x2, y2 float64, text string, fontName string, fontSize float64) {
	x, y, w, h := getXYWH(x1, y1, x2, y2)

	s.pdf.SetFont(fontName, "", fontSize)
	s.pdf.SetXY(x, y)
	s.pdf.SetCellMargin(0)
	s.pdf.CellFormat(w, h, text, "1" /*border*/, 0, "CB", false, 0, "")
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
		textCoord := strconv.Itoa(int(curX))
		s.pdf.Line(curX, 0, curX, s.dimY)
		s.pdf.Text(curX, 8, textCoord)
		s.pdf.Text(curX, s.dimY-8, textCoord)
	}

	for curY := float64(0); curY < s.dimY; curY += gap {
		textCoord := strconv.Itoa(int(curY))
		textWidth := s.pdf.GetStringWidth(textCoord)
		s.pdf.Line(0, curY, s.dimY, curY)
		s.pdf.Text(2, curY-1, textCoord)
		s.pdf.Text(s.dimX-textWidth-2, curY-1, textCoord)
	}

}
