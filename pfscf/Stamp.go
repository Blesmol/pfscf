package main

import (
	"fmt"
	"regexp"
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

var (
	regexSocietyID = regexp.MustCompile(`^\s*(\d*)\s*-\s*(\d*)\s*$`)
)

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

// AddTextCell adds a text cell to the stamp.
func (s *Stamp) AddTextCell(x, y, w, h float64, font string, fontsize float64, align string, text string) {
	s.pdf.SetFont(font, "", fontsize)
	s.pdf.SetXY(x, y)
	s.pdf.SetCellMargin(0)
	s.pdf.CellFormat(w, h, text, s.cellBorder, 0, align, false, 0, "")
}

// DrawRectangle draws a rectangle on the stamp.
func (s *Stamp) DrawRectangle(x, y, w, h float64, style string, r, g, b int) {
	s.pdf.SetFillColor(r, g, b)
	s.pdf.Rect(x, y, w, h, style)
}

// GetStringWidth returns the width of a given string
func (s *Stamp) GetStringWidth(str string, font string, style string, fontsize float64) (result float64) {
	s.pdf.SetFont(font, style, fontsize)
	return s.pdf.GetStringWidth(str)
}

// WriteToFile writes the content of the Stamp object into a PDF file.
// The Stamp object should not be used anymore after that.
func (s *Stamp) WriteToFile(filename string) (err error) {
	if !IsSet(filename) {
		return fmt.Errorf("No filename provided")
	}
	return s.pdf.OutputFileAndClose(filename)
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
