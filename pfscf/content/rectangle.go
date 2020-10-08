package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/canvas"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeRectangle = "rectangle"
)

var (
	validStyles = []string{"filled", "strikeout"}
)

// rectangle needs a description
type rectangle struct {
	X, Y         float64
	X2, Y2       float64
	Color        string
	Transparency float64 // TODO convert to ptr
	Style        string
	Canvas       string
	Presets      []string
}

func newRectangle() *rectangle {
	var e rectangle
	e.Presets = make([]string, 0)
	return &e
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (e *rectangle) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Color", "Canvas", "Style")
	if err != nil {
		return contentValErr(e, err)
	}

	err = utils.CheckFieldsAreInRange(e, 0.0, 100.0, "X", "Y", "X2", "Y2")
	if err != nil {
		return contentValErr(e, err)
	}

	if _, exists := canvasStore.Get(e.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", e.Canvas)
		return contentValErr(e, err)
	}

	if _, _, _, err = parseColor(e.Color); err != nil {
		return contentValErr(e, err)
	}

	if e.Transparency < 0.0 || e.Transparency > 1.0 {
		err = fmt.Errorf("Transparency value outside of range 0.0 to 1.0: %v", e.Transparency)
		return contentValErr(e, err)
	}

	if !utils.Contains(validStyles, e.Style) {
		err = fmt.Errorf("Unknown style '%v'. Supported styles are %v", e.Style, validStyles)
		return contentValErr(e, err)
	}

	return nil
}

// resolve the presets for this content object.
func (e *rectangle) resolve(ps preset.Store) (err error) {
	// Rectangle with zero width or height? No output!
	if e.X == e.X2 || e.Y == e.Y2 {
		return nil
	}

	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(e.Presets...); err != nil {
		err = fmt.Errorf("Error resolving content: %v", err)
		return
	}

	for _, presetID := range e.Presets {
		preset, _ := ps.Get(presetID)
		if err = preset.FillPublicFieldsFromPreset(e, "Presets"); err != nil {
			err = fmt.Errorf("Error resolving content: %v", err)
			return
		}
	}

	// defaults
	if !utils.IsSet(e.Transparency) {
		e.Transparency = 0.0
	}
	if !utils.IsSet(e.Style) {
		e.Style = validStyles[0]
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (e *rectangle) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	r, g, b, err := parseColor(e.Color)
	if err != nil {
		return err
	}

	switch e.Style {
	case "filled":
		style := stamp.OutputStyle{Style: "F", FillR: r, FillG: g, FillB: b, Transparency: e.Transparency}
		s.DrawRectangle(e.Canvas, e.X, e.Y, e.X2, e.Y2, style)
	case "strikeout":
		style := stamp.OutputStyle{DrawR: r, DrawB: b, DrawG: g, Linewidth: 2.5}
		s.DrawLine(e.Canvas, e.X, e.Y, e.X2, e.Y2, style)
		s.DrawLine(e.Canvas, e.X, e.Y2, e.X2, e.Y, style)
	default:
		utils.Assert(false, "Should be unreachable, or some valid case is missing here")
	}

	return nil
}

// deepCopy creates a deep copy of this entry.
func (e *rectangle) deepCopy() Entry {
	copy := *e
	copy.Presets = append(make([]string, 0), e.Presets...)

	return &copy
}
