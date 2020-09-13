package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
)

const (
	typeCanvas = "canvas"
)

type canvas struct {
	X, Y    float64
	X2, Y2  float64
	Content ListStore
	Presets []string
}

func newCanvas() *canvas {
	var ce canvas
	ce.Content = NewListStore()
	ce.Presets = make([]string, 0)
	return &ce
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce *canvas) isValid(paramStore *param.Store) (err error) {
	err = checkFieldsAreInRange(ce, "X", "Y", "X2", "Y2")
	if err != nil {
		return contentValErr(ce, err)
	}

	if ce.X == ce.X2 {
		err = fmt.Errorf("Coordinates for X axis are equal: %v", ce.X)
		return contentValErr(ce, err)
	}

	if ce.Y == ce.Y2 {
		err = fmt.Errorf("Coordinates for Y axis are equal: %v", ce.Y)
		return contentValErr(ce, err)
	}

	return ce.Content.IsValid(paramStore)
}

// resolve the presets for this content object.
func (ce *canvas) resolve(ps preset.Store) (err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.Presets...); err != nil {
		err = fmt.Errorf("Error resolving content: %v", err)
		return
	}

	// apply presets
	for _, presetID := range ce.Presets {
		preset, _ := ps.Get(presetID)
		if err = fillPublicFieldsFromPreset(ce, &preset, "Presets"); err != nil {
			err = fmt.Errorf("Error resolving content: %v", err)
			return
		}
	}

	return ce.Content.Resolve(ps)
}

// generateOutput generates the output for this canvas object.
func (ce *canvas) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	s.AddCanvas(ce.X, ce.Y, ce.X2, ce.Y2)
	defer s.RemoveCanvas()

	return ce.Content.GenerateOutput(s, as)
}

// deepCopy creates a deep copy of this entry.
func (ce *canvas) deepCopy() Entry {
	copy := canvas{
		X:       ce.X,
		Y:       ce.Y,
		X2:      ce.X2,
		Y2:      ce.Y2,
		Content: ce.Content.deepCopy(),
	}
	for _, preset := range ce.Presets {
		copy.Presets = append(copy.Presets, preset)
	}

	return &copy
}
