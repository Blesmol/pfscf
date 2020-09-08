package stamp

import (
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Stamp is a wraper for a PDF page
type Stamp struct {
	pdf        *gofpdf.Fpdf
	dimX       float64
	dimY       float64
	cellBorder string
	canvas     canvas
}

const (
	minFontSize = 4.0
)

// NewStamp creates a new Stamp object.
func NewStamp(dimX float64, dimY float64) (s *Stamp) {
	s = new(Stamp)

	s.dimX = dimX
	s.dimY = dimY

	s.canvas = newCanvas(0.0, 0.0, dimX, dimY)

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

func (s *Stamp) shouldDrawCellBorder() bool {
	return s.cellBorder == "1"
}

// AddCanvas adds another canvas to set a smaller canvas on this stamp.
func (s *Stamp) AddCanvas(x1Pct, y1Pct, x2Pct, y2Pct float64) {
	if s.shouldDrawCellBorder() {
		s.DrawRectangle(x1Pct, y1Pct, x2Pct, y2Pct, "D", 0, 255, 0, 0, 0, 0)
	}
	s.canvas = s.canvas.getSubCanvas(x1Pct, y1Pct, x2Pct, y2Pct)
}

// RemoveCanvas removes the newest canvas from this stamp.
func (s *Stamp) RemoveCanvas() {
	s.canvas = s.canvas.getParentCanvas()
}

// DeriveFontsize checks whether the provided text fits into the given width, if the current
// font and fontsize is used. If it does not fit, the size is reduced until it fits or until a
// minimum font size is reached.
func (s *Stamp) DeriveFontsize(ptWidth float64, font string, fontsize float64, text string) (result float64) {
	// TODO extend to also take height into account?
	// TODO convert to percent and remove call from AddTextCell
	for autoFontsize := fontsize; autoFontsize >= minFontSize; autoFontsize -= 0.25 {
		s.pdf.SetFont(font, "", autoFontsize)
		if s.pdf.GetStringWidth(text) <= ptWidth {
			return autoFontsize
		}
	}
	return minFontSize
}

// DeriveY2 takes two coordinates on the Y axis and the fontsize, and in case
// the second coordinate is 0.0 will calculcate a proper y2 coordinate based
// on y1 and the fontsize.
func (s *Stamp) DeriveY2(y1Pct, y2Pct, fontsizePt float64) (y2 float64) {
	if y2Pct != 0.0 {
		return y2Pct
	}

	_, fontsizePct := s.canvas.relPtToPct(0.0, fontsizePt)
	return y1Pct - fontsizePct
}

// AddTextCell adds a text cell to the stamp.
func (s *Stamp) AddTextCell(x1Pct, y1Pct, x2Pct, y2Pct float64, font string, fontsize float64, align string, text string, autoShrink bool) {
	xPt, yPt, wPt, hPt := s.canvas.pctToPt(x1Pct, y1Pct, x2Pct, y2Pct)

	effectiveFontsize := fontsize
	if autoShrink {
		effectiveFontsize = s.DeriveFontsize(wPt, font, fontsize, text)
	}

	s.pdf.SetFont(font, "", effectiveFontsize)
	s.pdf.SetXY(xPt, yPt)
	s.pdf.SetCellMargin(0)
	if s.shouldDrawCellBorder() {
		s.pdf.SetDrawColor(0, 0, 0)
	}
	s.pdf.CellFormat(wPt, hPt, text, s.cellBorder, 0, align, false, 0, "")
}

// DrawRectangle draws a rectangle on the stamp.
func (s *Stamp) DrawRectangle(x1Pct, y1Pct, x2Pct, y2Pct float64, style string, dr, dg, db int, fr, fg, fb int) {
	xPt, yPt, wPt, hPt := s.canvas.pctToPt(x1Pct, y1Pct, x2Pct, y2Pct)

	s.pdf.SetDrawColor(dr, dg, db)
	s.pdf.SetFillColor(fr, fg, fb)
	s.pdf.Rect(xPt, yPt, wPt, hPt, style)
}

// WriteToFile writes the content of the Stamp object into a PDF file.
// The Stamp object should not be used anymore after that.
func (s *Stamp) WriteToFile(filename string) (err error) {
	if !utils.IsSet(filename) {
		return fmt.Errorf("No filename provided")
	}
	// TODO invalidate object as part of this call
	return s.pdf.OutputFileAndClose(filename)
}

// CreateMeasurementCoordinates overlays the stamp with a set of lines
func (s *Stamp) CreateMeasurementCoordinates(majorGap, minorGap float64) {
	utils.Assert(majorGap > 0, "Provided gap should be greater than 0")

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

	maxLabelWidth := s.pdf.GetStringWidth("x:99") + extraSpace // 100% won't be reached
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
			curX, curY := s.canvas.relPctToPt(curPercent, curPercent)

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
		curX, curY := s.canvas.relPctToPt(curPercent, curPercent)

		if curX >= minX && curX <= maxX {
			s.pdf.Line(curX, minY, curX, maxY)

			labelText := fmt.Sprintf("x:%v", strconv.Itoa(int(curPercent)))
			labelWidth := s.pdf.GetStringWidth(labelText)

			labelXPos := curX - (labelWidth / 2.0) // place in middle of line
			labelYTopPos := 0.0 + maxLabelHeight - extraSpace
			labelYBottomPos := s.dimY - extraSpace

			s.pdf.Text(labelXPos, labelYTopPos, labelText)
			s.pdf.Text(labelXPos, labelYBottomPos, labelText)
		}

		if curY >= minY && curY <= maxY {
			s.pdf.Line(minX, curY, maxX, curY)

			labelText := fmt.Sprintf("y:%v", strconv.Itoa(int(curPercent)))
			labelWidth := s.pdf.GetStringWidth(labelText)
			labelXLeft := 0.0 + extraSpace
			labelXRight := s.dimX - labelWidth
			labelYPos := curY + (maxLabelHeight / 2.0) - extraSpace

			s.pdf.Text(labelXLeft, labelYPos, labelText)
			s.pdf.Text(labelXRight, labelYPos, labelText)
		}
	}
}
