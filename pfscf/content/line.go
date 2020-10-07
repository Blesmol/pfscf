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
	typeLine = "line"
)

type line struct {
	X, Y    float64
	X2, Y2  float64
	Linewidth float64
	Color   string
	Canvas  string
	Presets []string
}

func newLine() *line {
	var ce line
	ce.Presets = make([]string, 0)
	return &ce
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce *line) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(ce, "Color", "Canvas", "Linewidth")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = utils.CheckFieldsAreInRange(ce, 0.0, 100.0, "X", "Y", "X2", "Y2")
	if err != nil {
		return contentValErr(ce, err)
	}

	if ce.X == ce.X2 && ce.Y == ce.Y2 {
		err = fmt.Errorf("Line should have different coordinates either for x and x2 or for y and y2: x/x2: %v, y/y2: %v", ce.X, ce.Y)
		return contentValErr(ce, err)
	}

	if _, exists := canvasStore.Get(ce.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", ce.Canvas)
		return contentValErr(ce, err)
	}

	if _, _, _, err = parseColor(ce.Color); err != nil {
		return contentValErr(ce, err)
	}

	return nil
}

// resolve the presets for this content object.
func (ce *line) resolve(ps preset.Store) (err error) {
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
	if !utils.IsSet(ce.Color) {
		ce.Color = "black"
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (ce *line) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	r, g, b, err := parseColor(ce.Color)
	if err != nil {
		return err
	}

	style := stamp.OutputStyle{DrawR: r, DrawG: g, DrawB: b, Linewidth: ce.Linewidth}
	s.DrawLine(ce.Canvas, ce.X, ce.Y, ce.X2, ce.Y2, style)

	return nil
}

// deepCopy creates a deep copy of this entry.
func (ce *line) deepCopy() Entry {
	copy := *ce
	copy.Presets = append(make([]string, 0), ce.Presets...)

	return &copy
}
