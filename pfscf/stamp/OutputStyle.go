package stamp

import "github.com/jung-kurt/gofpdf"

// OutputStyle allows to provide only the required drawing parameters
type OutputStyle struct {
	Style               string // F=filled, D=outline, FD=both, default=D
	FillR, FillG, FillB int
	DrawR, DrawG, DrawB int
	TextR, TextG, TextB int
	Transparency        float64
	BlendMode           string
	Linewidth           float64
	Fontsize            float64
	CellMargin          float64
}

func getOutputStyle(p *gofpdf.Fpdf) (s OutputStyle) {
	// s.Style cannot be retrieved, as this is a parameter set for the concrete call. Nevertheless, it may be used elsewhere
	// Font cannot be retrieved, no proper API provided

	s.FillR, s.FillG, s.FillB = p.GetFillColor()
	s.DrawR, s.DrawG, s.DrawB = p.GetDrawColor()
	s.TextR, s.TextG, s.TextB = p.GetTextColor()
	s.Transparency, s.BlendMode = p.GetAlpha()
	s.Linewidth = p.GetLineWidth()
	s.Fontsize, _ = p.GetFontSize()
	s.CellMargin = p.GetCellMargin()

	return s
}

// setOutputStyle sets all pdf parameters that are not equal to the desired values.
// TODO this currently sets *everything* in here... even if the values are currently not set.
func (s OutputStyle) setOutputStyle(p *gofpdf.Fpdf) {
	if r, g, b := p.GetFillColor(); s.FillR != r || s.FillG != g || s.FillB != b {
		p.SetFillColor(s.FillR, s.FillG, s.FillB)
	}

	if r, g, b := p.GetDrawColor(); s.DrawR != r || s.DrawG != g || s.DrawB != b {
		p.SetDrawColor(s.DrawR, s.DrawG, s.DrawB)
	}

	if r, g, b := p.GetTextColor(); s.TextR != r || s.TextG != g || s.TextB != b {
		p.SetTextColor(s.TextR, s.TextG, s.TextB)
	}

	if tr, blend := p.GetAlpha(); tr != s.Transparency || blend != s.BlendMode {
		p.SetAlpha(s.Transparency, s.BlendMode)
	}

	if s.Linewidth != p.GetLineWidth() {
		p.SetLineWidth(s.Linewidth)
	}

	if fs, _ := p.GetFontSize(); fs != s.Fontsize {
		p.SetFontSize(fs)
	}

	if s.CellMargin != p.GetCellMargin() {
		p.SetCellMargin(s.CellMargin)
	}
}
