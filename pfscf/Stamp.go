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

// pctToPt converts the provided percent coordinates into absolute
// point coordinates for the current stamp object.
// A value of, e.g. 10% should be passed as 10.0, not as 0.10
func (s *Stamp) pctToPt(x, y float64) (xPt, yPt float64) {
	return s.dimX * (x / 100.0), s.dimY * (y / 100.0)
}

// ptToPct converts the provided point coordinates into percent
// coordinates for the current stamp object.
// A value of, e.g. 10% will be returned as 10.0, not as 0.10
func (s *Stamp) ptToPct(x, y float64) (xPct, yPct float64) {
	return (100.0 / s.dimX) * x, (100.0 / s.dimY) * y
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

// getXYWHasPt transforms two sets of x/y coordinates in percent into a single
// set of x/y coodinates and a pair of width/height values in points unit
func (s *Stamp) getXYWHasPt(x1, y1, x2, y2 float64) (x, y, w, h float64) {
	xPct, yPct, wPct, hPct := getXYWH(x1, y1, x2, y2)
	x, y = s.pctToPt(xPct, yPct)
	w, h = s.pctToPt(wPct, hPct)
	return x, y, w, h
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

	const (
		labelFont      = "Arial"
		labelFontStyle = "B"
		labelFontSize  = float64(6)
		borderAreaPt   = float64(16) // do not add lines and text if that near to the page border
		majorLineWidth = float64(0.5)
		minorLineWidth = float64(0.1)

		extraSpace = float64(1.0) // points
	)

	// store away old settings and reset at the end (you never know...)
	formerR, formerB, formerG := s.pdf.GetDrawColor()
	defer s.pdf.SetDrawColor(formerR, formerG, formerB)
	formerLineWidth := s.pdf.GetLineWidth()
	defer s.pdf.SetLineWidth(formerLineWidth)

	s.pdf.SetFont(labelFont, labelFontStyle, labelFontSize) // Used for the labels at the borders

	maxLabelWidth := s.pdf.GetStringWidth("x:99%") + extraSpace // 100% won't be reached
	maxLabelHeight := labelFontSize + extraSpace

	minX := 0.0 + maxLabelWidth
	maxX := s.dimX - maxLabelWidth
	minY := 0.0 + maxLabelHeight
	maxY := s.dimY - maxLabelHeight

	// ignore minor gap if 0 (or below)
	if minorGap > 0 {
		// settings for minor gap drawing
		s.pdf.SetDrawColor(196, 196, 196) // something lightgrayish
		s.pdf.SetLineWidth(minorLineWidth)

		for curPercent := 0.0; curPercent <= 100.0; curPercent += minorGap {
			curX, curY := s.pctToPt(curPercent, curPercent)

			if curX >= minX && curX <= maxX {
				s.pdf.Line(curX, minY, curX, maxY)
			}

			if curY >= minY && curY <= maxY {
				s.pdf.Line(minX, curY, maxX, curY)
			}
		}
	}

	// settings for major gap drawing
	s.pdf.SetDrawColor(64, 64, 255) // something blueish
	s.pdf.SetLineWidth(majorLineWidth)

	// draw major gap X lines with labels
	for curPercent := 0.0; curPercent <= 100.0; curPercent += majorGap {
		curX, curY := s.pctToPt(curPercent, curPercent)

		if curX >= minX && curX <= maxX {
			s.pdf.Line(curX, minY, curX, maxY)

			labelText := fmt.Sprintf("x:%v%%", strconv.Itoa(int(curPercent)))
			labelWidth := s.pdf.GetStringWidth(labelText)

			labelXPos := curX - (labelWidth / 2.0) // place in middle of line
			labelYTopPos := 0.0 + maxLabelHeight - extraSpace
			labelYBottomPos := s.dimY - extraSpace

			s.pdf.Text(labelXPos, labelYTopPos, labelText)
			s.pdf.Text(labelXPos, labelYBottomPos, labelText)
		}

		if curY >= minY && curY <= maxY {
			s.pdf.Line(minX, curY, maxX, curY)

			labelText := fmt.Sprintf("y:%v%%", strconv.Itoa(int(curPercent)))
			labelWidth := s.pdf.GetStringWidth(labelText)
			labelXLeft := 0.0 + extraSpace
			labelXRight := s.dimX - labelWidth
			labelYPos := curY + (maxLabelHeight / 2.0) - extraSpace

			s.pdf.Text(labelXLeft, labelYPos, labelText)
			s.pdf.Text(labelXRight, labelYPos, labelText)
		}
	}

	/*
		for curX := float64(0); curX < s.dimX; curX += majorGap {
			if curX < (0+borderAreaPt) || curX > (s.dimX-borderAreaPt) {
				continue
			}

			coordString := fmt.Sprintf("x:%v", strconv.Itoa(int(curX)))
			textWidth := s.pdf.GetStringWidth(coordString)
			textOffset := textWidth / 2 // place in the middle of the line
			textTopBorderMargin := labelFontSize + 2
			textBottomBorderMargin := float64(2)
			lineTopBorderMargin := textTopBorderMargin + 2
			lineBottomBorderMargin := textBottomBorderMargin + labelFontSize + 2

			s.pdf.Line(curX, 0+lineTopBorderMargin, curX, s.dimY-lineBottomBorderMargin)
			s.pdf.Text(curX-textOffset, textTopBorderMargin, coordString)
			s.pdf.Text(curX-textOffset, s.dimY-textBottomBorderMargin, coordString)
		}

		// draw major gap Y lines with labels
		for curY := float64(0); curY < s.dimY; curY += majorGap {
			if curY < (0+borderAreaPt) || curY > (s.dimY-borderAreaPt) {
				continue
			}

			coordString := fmt.Sprintf("y:%v", strconv.Itoa(int(curY)))
			textWidth := s.pdf.GetStringWidth(coordString)
			textPosY := curY + (labelFontSize / 2) - 1
			lineBorderMargin := textWidth + 4 // enough space for the text plus a little

			s.pdf.Line(0+lineBorderMargin, curY, s.dimX-lineBorderMargin, curY)
			s.pdf.Text(2, textPosY, coordString)
			s.pdf.Text(s.dimX-textWidth-2, textPosY, coordString)
		}
	*/
}
