package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeTextCell = "textCell"
)

// textCell is the final type to implement textCells.
// TODO switch to pointers to distinguish between unset values and zero values?
type textCell struct {
	Value    string
	X, Y     float64
	X2, Y2   float64
	Font     string
	Fontsize float64
	Align    string
	Presets  []string
}

func newTextCell() *textCell {
	var ce textCell
	ce.Presets = make([]string, 0)
	return &ce
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce *textCell) isValid(paramStore *param.Store) (err error) {
	err = utils.CheckFieldsAreSet(ce, "Value", "Font", "Fontsize")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = utils.CheckFieldsAreInRange(ce, 0.0, 100.0, "X", "Y", "X2", "Y2")
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

	return nil
}

// resolve the presets for this content object.
func (ce *textCell) resolve(ps preset.Store) (err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.Presets...); err != nil {
		err = fmt.Errorf("Error resolving content: %v", err)
		return
	}

	// apply presets
	for _, presetID := range ce.Presets {
		preset, _ := ps.Get(presetID)
		if err = preset.FillPublicFieldsFromPreset(ce, "Presets"); err != nil {
			err = fmt.Errorf("Error resolving content: %v", err)
			return
		}
	}

	return nil
}

// generateOutput generates the output for this textCell object.
func (ce *textCell) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	value := getValue(ce.Value, as)
	if value == nil {
		return nil // nothing to do here...
	}

	y2 := s.DeriveY2(ce.Y, ce.Y2, ce.Fontsize)
	s.AddTextCell(ce.X, ce.Y, ce.X2, y2, ce.Font, ce.Fontsize, ce.Align, *value, true)

	return nil
}

// deepCopy creates a deep copy of this entry.
// TODO create generic deep-copy function for public fields
func (ce *textCell) deepCopy() Entry {
	copy := textCell{
		Value:    ce.Value,
		X:        ce.X,
		Y:        ce.Y,
		X2:       ce.X2,
		Y2:       ce.Y2,
		Font:     ce.Font,
		Fontsize: ce.Fontsize,
		Align:    ce.Align,
	}
	for _, preset := range ce.Presets {
		copy.Presets = append(copy.Presets, preset)
	}

	return &copy
}
