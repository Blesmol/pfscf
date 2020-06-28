package main

import (
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

// Stamp is a wraper for a PDF page
type Stamp struct {
	pdf        *gofpdf.Fpdf
	dimX       float64
	dimY       float64
	cellBorder string
}

// NewStamp creates a new Stamp object.
func NewStamp(dimX float64, dimY float64) (s *Stamp) {
	s = new(Stamp)

	s.dimX = dimX
	s.dimY = dimY

	s.SetCellBorder(false)

	s.pdf = gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size:    gofpdf.SizeType{Wd: dimX, Ht: dimY},
	})
	s.pdf.SetMargins(0, 0, 0)
	s.pdf.SetAutoPageBreak(false, 0)
	s.pdf.AddPage() // 0,0 is top-left. To change use AddPageFormat() instead

	return s
}

// SetCellBorder sets whether the border around cells should be drawn.
func (s *Stamp) SetCellBorder(shouldDrawBorder bool) {
	if shouldDrawBorder {
		s.cellBorder = "1"
	} else {
		s.cellBorder = "0"
	}
}

// AddText adds a portion of text at the specified coordinates
func (s *Stamp) AddText(x, y float64, text string, fontName string, fontSize float64) {
	s.pdf.SetFont(fontName, "", fontSize)
	s.pdf.Text(x, y, text)
}

// getXYWH transforms two sets of x/y coordinates into a single set of
// x/y coordinates and a pair of width/height values.
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

// AddContent is a generic function to add content to a stamp. It will
// internally check the content type and call the appropriate subfunction.
func (s *Stamp) AddContent(ce ContentData, value *string) {
	switch ce.Type {
	case "textCell":
		Assert(value != nil, "Content type 'textCell' needs an input value")
		isValid, err := ce.IsValid()
		Assert(isValid, fmt.Sprintf("Received error: %v\n", err))
		s.AddCellText(ce.X1, ce.Y1, ce.X2, ce.Y2, *value, ce.Font, ce.Fontsize, ce.Align)
	default:
		panic("Unknown content type: " + ce.Type)
	}
}

// AddCellText is a better version of AddText()
func (s *Stamp) AddCellText(x1, y1, x2, y2 float64, text string, fontName string, fontSize float64, alignString string) {
	x, y, w, h := getXYWH(x1, y1, x2, y2)

	s.pdf.SetFont(fontName, "", fontSize)
	s.pdf.SetXY(x, y)
	s.pdf.SetCellMargin(0)
	s.pdf.CellFormat(w, h, text, s.cellBorder, 0, alignString, false, 0, "")
}

// WriteToFile writes the contect of the Stamp object into a PDF file.
// The Stamp object should not be used anymore after that.
func (s *Stamp) WriteToFile(filename string) {
	err := s.pdf.OutputFileAndClose(filename)
	AssertNoError(err)
}

// CreateMeasurementCoordinates overlays the stamp with a set of lines
func (s *Stamp) CreateMeasurementCoordinates(majorGap, minorGap float64) {
	Assert(majorGap > 0, "Provided gap should be greater than 0")

	const fontSize = float64(6)
	const borderArea = float64(16) // do not add lines and text if that near to the page border
	const majorLineWidth = float64(0.5)
	const minorLineWidth = float64(0.1)

	// store away old settings and reset at the end (you never know...)
	formerR, formerB, formerG := s.pdf.GetDrawColor()
	formerLineWidth := s.pdf.GetLineWidth()
	defer s.pdf.SetDrawColor(formerR, formerG, formerB)
	defer s.pdf.SetLineWidth(formerLineWidth)

	// ignore minor gap if 0 or below
	if minorGap > 0 {
		// settings for minor gap drawing
		s.pdf.SetDrawColor(196, 196, 196) // something lightgrayish
		s.pdf.SetLineWidth(minorLineWidth)

		// draw minor gap X lines
		for curX := float64(0); curX < s.dimX; curX += minorGap {
			if curX < (0+borderArea) || curX > (s.dimX-borderArea) {
				continue
			}
			s.pdf.Line(curX, 0+borderArea, curX, s.dimY-borderArea)
		}

		// draw minor gap Y
		for curY := float64(0); curY < s.dimY; curY += minorGap {
			if curY < (0+borderArea) || curY > (s.dimY-borderArea) {
				continue
			}
			s.pdf.Line(0+borderArea, curY, s.dimX-borderArea, curY)
		}
	}

	// settings for major gap drawing
	s.pdf.SetFont("Arial", "B", fontSize)
	s.pdf.SetDrawColor(64, 64, 255) // something blueish
	s.pdf.SetLineWidth(majorLineWidth)

	// draw major gap X lines with labels
	for curX := float64(0); curX < s.dimX; curX += majorGap {
		if curX < (0+borderArea) || curX > (s.dimX-borderArea) {
			continue
		}

		coordString := fmt.Sprintf("x:%v", strconv.Itoa(int(curX)))
		textWidth := s.pdf.GetStringWidth(coordString)
		textOffset := textWidth / 2 // place in the middle of the line
		textTopBorderMargin := fontSize + 2
		textBottomBorderMargin := float64(2)
		lineTopBorderMargin := textTopBorderMargin + 2
		lineBottomBorderMargin := textBottomBorderMargin + fontSize + 2

		s.pdf.Line(curX, 0+lineTopBorderMargin, curX, s.dimY-lineBottomBorderMargin)
		s.pdf.Text(curX-textOffset, textTopBorderMargin, coordString)
		s.pdf.Text(curX-textOffset, s.dimY-textBottomBorderMargin, coordString)
	}

	// draw major gap Y lines with labels
	for curY := float64(0); curY < s.dimY; curY += majorGap {
		if curY < (0+borderArea) || curY > (s.dimY-borderArea) {
			continue
		}

		coordString := fmt.Sprintf("y:%v", strconv.Itoa(int(curY)))
		textWidth := s.pdf.GetStringWidth(coordString)
		textPosY := curY + (fontSize / 2) - 1
		lineBorderMargin := textWidth + 4 // enough space for the text plus a little

		s.pdf.Line(0+lineBorderMargin, curY, s.dimX-lineBorderMargin, curY)
		s.pdf.Text(2, textPosY, coordString)
		s.pdf.Text(s.dimX-textWidth-2, textPosY, coordString)
	}
}
