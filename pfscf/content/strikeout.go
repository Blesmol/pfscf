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
	typeStrikeout = "strikeout"
)

// rectangle needs a description
type strikeout struct {
	X, Y      float64
	X2, Y2    float64
	Size      float64
	Linewidth float64
	Color     string
	Canvas    string
	Presets   []string
}

func newStrikeout() *strikeout {
	var e strikeout
	e.Presets = make([]string, 0)
	return &e
}

func (e *strikeout) shouldDrawCentered() bool {
	return utils.IsSet(e.Size) && utils.IsSet(e.X) && utils.IsSet(e.Y)
}

func (e *strikeout) shouldDrawArea() bool {
	return utils.IsSet(e.X2) && utils.IsSet(e.Y2)
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (e *strikeout) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Canvas")
	if err != nil {
		return contentValErr(e, err)
	}

	err = utils.CheckFieldsAreInRange(e, 0.0, 100.0, "X", "Y", "X2", "Y2", "Size", "Linewidth")
	if err != nil {
		return contentValErr(e, err)
	}

	if _, _, _, err = parseColor(e.Color); err != nil {
		return contentValErr(e, err)
	}

	if e.shouldDrawArea() && e.shouldDrawCentered() {
		err := fmt.Errorf("Can only have either a 'size' value or 'x2'&'y2' values, but not both at the same time")
		return contentValErr(e, err)
	}

	if _, exists := canvasStore.Get(e.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", e.Canvas)
		return contentValErr(e, err)
	}

	return nil
}

// resolve the presets for this content object.
func (e *strikeout) resolve(ps preset.Store) (err error) {
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

	// set defaults
	if !utils.IsSet(e.Linewidth) {
		e.Linewidth = 1.0 // TODO set dynamically based on the size of the area
	}
	if !utils.IsSet(e.Color) {
		e.Color = "black"
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (e *strikeout) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	if !e.shouldDrawArea() && !e.shouldDrawCentered() { // Nothing to do? No output! Also a way to have this disabled per default
		return nil
	}

	r, g, b, err := parseColor(e.Color)
	if err != nil {
		return err
	}
	style := stamp.OutputStyle{DrawR: r, DrawB: g, DrawG: b, Linewidth: e.Linewidth}

	switch {
	case e.shouldDrawArea():
		s.DrawStrikeoutArea(e.Canvas, e.X, e.Y, e.X2, e.Y2, style)
	case e.shouldDrawCentered():
		s.DrawStrikeoutCentered(e.Canvas, e.X, e.Y, e.Size, style)
	}

	return nil
}

// deepCopy creates a deep copy of this entry.
func (e *strikeout) deepCopy() Entry {
	copy := *e
	copy.Presets = append(make([]string, 0), e.Presets...)

	return &copy
}
