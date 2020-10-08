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
	var e line
	e.Presets = make([]string, 0)
	return &e
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (e *line) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Color", "Canvas", "Linewidth")
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

	return nil
}

// resolve the presets for this content object.
func (e *line) resolve(ps preset.Store) (err error) {
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
	if !utils.IsSet(e.Color) {
		e.Color = "black"
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (e *line) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	// if the coordinates just form a point, then don't do anything and just return.
	// This also provides an easy way to have "inactive" lines that become active
	// by providing a different preset.
	if e.X == e.X2 && e.Y == e.Y2 {
		return nil
	}

	r, g, b, err := parseColor(e.Color)
	if err != nil {
		return err
	}

	style := stamp.OutputStyle{DrawR: r, DrawG: g, DrawB: b, Linewidth: e.Linewidth}
	s.DrawLine(e.Canvas, e.X, e.Y, e.X2, e.Y2, style)

	return nil
}

// deepCopy creates a deep copy of this entry.
func (e *line) deepCopy() Entry {
	copy := *e
	copy.Presets = append(make([]string, 0), e.Presets...)

	return &copy
}
