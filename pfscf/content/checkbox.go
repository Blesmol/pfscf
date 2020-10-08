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
	typeCheckbox = "checkbox"
)

// rectangle needs a description
type checkbox struct {
	X, Y    float64
	Size    float64
	Canvas  string
	Presets []string
}

func newCheckbox() *checkbox {
	var e checkbox
	e.Presets = make([]string, 0)
	return &e
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (e *checkbox) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Canvas")
	if err != nil {
		return contentValErr(e, err)
	}

	err = utils.CheckFieldsAreInRange(e, 0.0, 100.0, "X", "Y", "Size")
	if err != nil {
		return contentValErr(e, err)
	}

	if _, exists := canvasStore.Get(e.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", e.Canvas)
		return contentValErr(e, err)
	}

	return nil
}

// resolve the presets for this content object.
func (e *checkbox) resolve(ps preset.Store) (err error) {
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

	return nil
}

// generateOutput generates the output for this textCell object.
func (e *checkbox) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	if e.Size == 0.0 { // No size? No output! Also a way to have this disabled per default
		return nil
	}

	style := stamp.OutputStyle{DrawR: 0, DrawB: 0, DrawG: 0, Linewidth: 0.5}
	s.DrawStrikeoutCentered(e.Canvas, e.X, e.Y, e.Size, style)

	return nil
}

// deepCopy creates a deep copy of this entry.
func (e *checkbox) deepCopy() Entry {
	copy := *e
	copy.Presets = append(make([]string, 0), e.Presets...)

	return &copy
}
