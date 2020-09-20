package stamp

import (
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	minFontSize = 4.0
)

// Stamp is a wraper for a PDF page
type Stamp struct {
	pdf         *gofpdf.Fpdf
	dimX        float64
	dimY        float64
	cellBorder  string
	canvasStore map[string]canvas
	pageCanvas  canvas
}

// RectStyle allows to provide only the required drawing parameters for a rectangle
type RectStyle struct {
	Style                     string // F=filled, D=outline, FD=both, default=D
	FillR, FillG, FillB       int
	BorderR, BorderG, BorderB int
	Transparency              float64
}

// NewStamp creates a new Stamp object.
func NewStamp(dimX float64, dimY float64) (s *Stamp) {
	s = new(Stamp)

	s.dimX = dimX
	s.dimY = dimY

	s.canvasStore = make(map[string]canvas, 0)
	s.SetPageCanvas(0.0, 0.0, 100.0, 100.0)

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

// GetDimensions returns the x/y dimensions of this stamp
func (s *Stamp) GetDimensions() (x, y float64) {
	return s.dimX, s.dimY
}

// SetPageCanvas sets a new page canvas for this stamp. This function may only be
// called as long as no entries are contained in the Stamps' internal canvas store
func (s *Stamp) SetPageCanvas(x1Pct, y1Pct, x2Pct, y2Pct float64) {
	utils.Assert(len(s.canvasStore) == 0, "May only be called before entries are added to the stamp canvas store")
	s.pageCanvas = newCanvas(0.0, 0.0, s.dimX, s.dimY).getSubCanvas(x1Pct, y1Pct, x2Pct, y2Pct)
}

// AddCanvas adds a canvas to the new internal canvas store.
func (s *Stamp) AddCanvas(id string, x1Pct, y1Pct, x2Pct, y2Pct float64) {
	// check for duplicates
	_, exists := s.canvasStore[id]
	utils.Assert(!exists, "Dulpicates should not occur here")

	s.canvasStore[id] = s.pageCanvas.getSubCanvas(x1Pct, y1Pct, x2Pct, y2Pct)
}

// getCanvas returns the canvas denoted by the provided id. At the moment it is assumed that the
// name is always valid.
func (s *Stamp) getCanvas(id string) (c canvas) {
	c, exists := s.canvasStore[id]
	utils.Assert(exists, "It should have been checked before that only valid IDs exist at this point")
	return c
}

// DeriveFontsize checks whether the provided text fits into the given width, if the current
// font and fontsize is used. If it does not fit, the size is reduced until it fits or until a
// minimum font size is reached.
func (s *Stamp) DeriveFontsize(cellWidthPt, cellHeightPt float64, font string, fontsize float64, text string) (result float64) {
	// TODO convert to percent and remove call from AddTextCell

	currentFontsize := fontsize

	if cellHeightPt < currentFontsize {
		if cellHeightPt > minFontSize {
			currentFontsize = cellHeightPt
		} else {
			currentFontsize = minFontSize
		}
	}

	for ; currentFontsize >= minFontSize; currentFontsize -= 0.25 {
		s.pdf.SetFont(font, "", currentFontsize)
		if s.pdf.GetStringWidth(text) <= cellWidthPt {
			return currentFontsize
		}
	}
	return minFontSize
}

// DeriveY2 takes two coordinates on the Y axis and the fontsize, and in case
// the second coordinate is 0.0 will calculcate a proper y2 coordinate based
// on y1 and the fontsize.
func (s *Stamp) DeriveY2(canvasID string, y1Pct, y2Pct, fontsizePt float64) (y2 float64) {
	if y2Pct != 0.0 {
		return y2Pct
	}

	_, fontsizePct := s.getCanvas(canvasID).relPtToPct(0.0, fontsizePt)
	return y1Pct - fontsizePct
}

// AddTextCell adds a text cell to the stamp.
func (s *Stamp) AddTextCell(canvasID string, x1Pct, y1Pct, x2Pct, y2Pct float64, font string, fontsize float64, align string, text string, autoShrink bool) {
	xPt, yPt, wPt, hPt := s.getCanvas(canvasID).pctToPt(x1Pct, y1Pct, x2Pct, y2Pct)

	effectiveFontsize := fontsize
	if autoShrink {
		effectiveFontsize = s.DeriveFontsize(wPt, hPt, font, fontsize, text)
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
func (s *Stamp) DrawRectangle(canvasID string, x1Pct, y1Pct, x2Pct, y2Pct float64, rs RectStyle) {
	xPt, yPt, wPt, hPt := s.getCanvas(canvasID).pctToPt(x1Pct, y1Pct, x2Pct, y2Pct)

	oldAlpha, oldBlendMode := s.pdf.GetAlpha()
	s.pdf.SetAlpha(1.0-rs.Transparency, "Normal")
	defer s.pdf.SetAlpha(oldAlpha, oldBlendMode)

	s.pdf.SetDrawColor(rs.BorderR, rs.BorderG, rs.BorderB)
	s.pdf.SetFillColor(rs.FillR, rs.FillG, rs.FillB)
	s.pdf.Rect(xPt, yPt, wPt, hPt, rs.Style)
}

// DrawCanvases draws all canvases to the stamp
func (s *Stamp) DrawCanvases() {
	fontsize := 8.0
	r, g, b := 51, 204, 51
	s.pdf.SetDrawColor(r, g, b)
	s.pdf.SetFont("Helvetica", "", fontsize)

	for canvasID, canvas := range s.canvasStore {
		xPt, yPt, wPt, hPt := canvas.pctToPt(0.0, 0.0, 100.0, 100.0)

		// rectangle
		s.pdf.Rect(xPt, yPt, wPt, hPt, "D")

		// name/ID
		s.pdf.SetXY(xPt, yPt+hPt-fontsize)
		s.pdf.SetTextColor(r, g, b)
		s.pdf.CellFormat(wPt, fontsize, canvasID, "0", 0, "RM", false, 0, "")
	}
	s.pdf.SetTextColor(0, 0, 0)
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

// DrawCanvasGrid overlays the stamp with a set of lines
func (s *Stamp) DrawCanvasGrid(canvasID string) (err error) {
	const (
		labelFont      = "Arial"
		labelFontStyle = "B"
		labelFontSize  = float64(6)
		borderAreaPt   = float64(16) // do not add lines and text if that near to the page border
		majorLineWidth = float64(0.5)
		minorLineWidth = float64(0.1)

		extraSpace = float64(1.0) // points
	)

	// find canvas for which we should draw the grid
	var canvas canvas
	if canvasID != "" {
		var exists bool
		canvas, exists = s.canvasStore[canvasID]
		if !exists {
			return fmt.Errorf("Cannot find canvas '%v'", canvasID)
		}
	} else {
		canvas = s.pageCanvas
	}
	minX := canvas.xPt
	maxX := canvas.xPt + canvas.wPt
	pctX := canvas.wPt / 100.0
	minY := canvas.yPt
	maxY := canvas.yPt + canvas.hPt
	pctY := canvas.hPt / 100.0
	height := canvas.hPt
	width := canvas.wPt

	// calculate major and minor gap for the given canvas.
	// Try to keep major gaps a multiple of minor gaps
	var majorGapVertical, minorGapVertical float64
	var majorGapHorizontal, minorGapHorizontal float64
	for _, e := range []struct {
		major, minor *float64
		dim          float64
	}{
		{&majorGapHorizontal, &minorGapHorizontal, height},
		{&majorGapVertical, &minorGapVertical, width},
	} {
		switch {
		case e.dim < 300:
			*e.major = 10.0
			*e.minor = 2.0
		default:
			*e.major = 5.0
			*e.minor = 1.0
		}
	}

	// store away old relevant settings and reset at the end (you never know...)
	formerDR, formerDB, formerDG := s.pdf.GetDrawColor()
	defer s.pdf.SetDrawColor(formerDR, formerDG, formerDB)
	formerFR, formerFB, formerFG := s.pdf.GetFillColor()
	defer s.pdf.SetFillColor(formerFR, formerFG, formerFB)
	formerLineWidth := s.pdf.GetLineWidth()
	defer s.pdf.SetLineWidth(formerLineWidth)

	s.pdf.SetFont(labelFont, labelFontStyle, labelFontSize) // Used for the labels at the borders

	maxLabelWidth := s.pdf.GetStringWidth("x:99") + extraSpace // 100% won't be reached
	maxLabelHeight := labelFontSize + extraSpace

	// settings for minor gap drawing
	s.pdf.SetDrawColor(196, 196, 196) // something lightgrayish
	s.pdf.SetLineWidth(minorLineWidth)

	// minor vertical lines
	for curPercent := 0.0; curPercent <= 100.0; curPercent += minorGapVertical {
		// dont't draw anything directly on the border
		if curPercent == 0.0 || curPercent == 100.0 {
			continue
		}

		curX := minX + (curPercent * pctX)
		s.pdf.Line(curX, minY, curX, maxY)
	}

	// minor horizontal lines
	for curPercent := 0.0; curPercent <= 100.0; curPercent += minorGapHorizontal {
		// dont't draw anything directly on the border
		if curPercent == 0.0 || curPercent == 100.0 {
			continue
		}

		curY := minY + (curPercent * pctY)
		s.pdf.Line(minX, curY, maxX, curY)
	}

	// settings for major gap drawing
	s.pdf.SetDrawColor(64, 64, 255)   // something blueish
	s.pdf.SetFillColor(255, 255, 255) // white
	s.pdf.SetCellMargin(0)
	s.pdf.SetLineWidth(majorLineWidth)

	// major vertical lines
	for curPercent := 0.0; curPercent <= 100.0; curPercent += majorGapVertical {
		// dont't draw anything directly on the border
		if curPercent == 0.0 || curPercent == 100.0 {
			continue
		}

		// line
		curX := minX + (curPercent * pctX)
		s.pdf.Line(curX, minY, curX, maxY)

		// label
		labelText := fmt.Sprintf("x:%v", strconv.Itoa(int(curPercent)))
		labelWidth := s.pdf.GetStringWidth(labelText)
		labelX := curX + extraSpace
		rotX := curX

		// top label
		rotYTop := minY
		labelYTop := rotYTop - (maxLabelHeight / 2.0)
		if labelYTop-labelWidth < 0.0 { // yes, that looks weird, but: rotation
			rotYTop = 0.0 + labelWidth + extraSpace
			labelYTop = rotYTop - (maxLabelHeight / 2.0)
		}

		s.pdf.TransformBegin()
		s.pdf.TransformRotate(90, rotX, rotYTop)
		s.pdf.SetXY(labelX, labelYTop)
		s.pdf.CellFormat(labelWidth, maxLabelHeight, labelText, "0", 0, "LM", true, 0, "")
		s.pdf.TransformEnd()

		// bottom label
		rotYBottom := maxY + maxLabelWidth + 2*extraSpace
		labelYBottom := rotYBottom - (maxLabelHeight / 2.0)
		if rotYBottom > s.dimY {
			rotYBottom = s.dimY - extraSpace
			labelYBottom = rotYBottom - (maxLabelHeight / 2.0)
		}

		s.pdf.TransformBegin()
		s.pdf.TransformRotate(90, rotX, rotYBottom)
		s.pdf.SetXY(labelX, labelYBottom)
		s.pdf.CellFormat(maxLabelWidth, maxLabelHeight, labelText, "0", 0, "RM", true, 0, "")
		s.pdf.TransformEnd()
	}

	// major horizontal lines + labels
	for curPercent := 0.0; curPercent <= 100.0; curPercent += majorGapHorizontal {
		// dont't draw anything directly on the border
		if curPercent == 0.0 || curPercent == 100.0 {
			continue
		}

		curY := minY + (curPercent * pctY)

		// line
		s.pdf.Line(minX, curY, maxX, curY)

		// label
		labelText := fmt.Sprintf("y:%v", strconv.Itoa(int(curPercent)))
		labelWidth := s.pdf.GetStringWidth(labelText)
		labelY := curY - (maxLabelHeight / 2.0)

		labelXLeft := minX - (labelWidth + extraSpace)
		if labelXLeft < 0.0 {
			labelXLeft = 0.0
		}

		labelXRight := maxX + extraSpace
		if (labelXRight + labelWidth) > s.dimX {
			labelXRight = s.dimX - labelWidth
		}

		// left label
		s.pdf.SetXY(labelXLeft, labelY)
		s.pdf.CellFormat(labelWidth, maxLabelHeight, labelText, "0", 0, "RM", true, 0, "")

		// right label
		s.pdf.SetXY(labelXRight, labelY)
		s.pdf.CellFormat(labelWidth, maxLabelHeight, labelText, "0", 0, "LM", true, 0, "")
	}

	return nil
}
