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
	var ce rectangle
	ce.Presets = make([]string, 0)
	return &ce
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce *rectangle) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(ce, "Color", "Canvas", "Style")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = utils.CheckFieldsAreInRange(ce, 0.0, 100.0, "X", "Y", "X2", "Y2")
	if err != nil {
		return contentValErr(ce, err)
	}

	if _, exists := canvasStore.Get(ce.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", ce.Canvas)
		return contentValErr(ce, err)
	}

	if _, _, _, err = parseColor(ce.Color); err != nil {
		return contentValErr(ce, err)
	}

	if ce.Transparency < 0.0 || ce.Transparency > 1.0 {
		err = fmt.Errorf("Transparency value outside of range 0.0 to 1.0: %v", ce.Transparency)
		return contentValErr(ce, err)
	}

	if !utils.Contains(validStyles, ce.Style) {
		err = fmt.Errorf("Unknown style '%v'. Supported styles are %v", ce.Style, validStyles)
		return contentValErr(ce, err)
	}

	return nil
}

// resolve the presets for this content object.
func (ce *rectangle) resolve(ps preset.Store) (err error) {
	// Rectangle with zero width or height? No output!
	if ce.X == ce.X2 || ce.Y == ce.Y2 {
		return nil
	}

	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.Presets...); err != nil {
		err = fmt.Errorf("Error resolving content: %v", err)
		return
	}

	for _, presetID := range ce.Presets {
		preset, _ := ps.Get(presetID)
		if err = preset.FillPublicFieldsFromPreset(ce, "Presets"); err != nil {
			err = fmt.Errorf("Error resolving content: %v", err)
			return
		}
	}

	// defaults
	if !utils.IsSet(ce.Transparency) {
		ce.Transparency = 0.0
	}
	if !utils.IsSet(ce.Style) {
		ce.Style = validStyles[0]
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (ce *rectangle) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	r, g, b, err := parseColor(ce.Color)
	if err != nil {
		return err
	}

	switch ce.Style {
	case "filled":
		style := stamp.OutputStyle{Style: "F", FillR: r, FillG: g, FillB: b, Transparency: ce.Transparency}
		s.DrawRectangle(ce.Canvas, ce.X, ce.Y, ce.X2, ce.Y2, style)
	case "strikeout":
		style := stamp.OutputStyle{DrawR: r, DrawB: b, DrawG: g, Linewidth: 2.5}
		s.DrawLine(ce.Canvas, ce.X, ce.Y, ce.X2, ce.Y2, style)
		s.DrawLine(ce.Canvas, ce.X, ce.Y2, ce.X2, ce.Y, style)
	default:
		utils.Assert(false, "Should be unreachable, or some valid case is missing here")
	}

	return nil
}

// deepCopy creates a deep copy of this entry.
func (ce *rectangle) deepCopy() Entry {
	copy := *ce
	copy.Presets = append(make([]string, 0), ce.Presets...)

	return &copy
}
