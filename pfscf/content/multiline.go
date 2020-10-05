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
	typeMultiline = "multiline"
)

type multiline struct {
	Value    string
	X, Y     float64
	X2, Y2   float64
	Lines    int
	Font     string
	Fontsize float64
	Align    string
	Canvas   string
	Presets  []string
}

func newMultiline() *multiline {
	var e multiline
	e.Presets = make([]string, 0)
	return &e
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (e *multiline) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Value", "Font", "Fontsize", "Canvas", "Lines")
	if err != nil {
		return contentValErr(e, err)
	}

	err = utils.CheckFieldsAreInRange(e, 0.0, 100.0, "X", "Y", "X2", "Y2")
	if err != nil {
		return contentValErr(e, err)
	}

	if e.X == e.X2 {
		err = fmt.Errorf("Coordinates for X axis are equal: %v", e.X)
		return contentValErr(e, err)
	}

	if e.Y == e.Y2 {
		err = fmt.Errorf("Coordinates for Y axis are equal: %v", e.Y)
		return contentValErr(e, err)
	}

	if _, exists := canvasStore.Get(e.Canvas); !exists {
		err = fmt.Errorf("Canvas '%v' does not exist", e.Canvas)
		return contentValErr(e, err)
	}

	return nil
}

// resolve the presets for this content object.
func (e *multiline) resolve(ps preset.Store) (err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(e.Presets...); err != nil {
		err = fmt.Errorf("Error resolving content: %v", err)
		return
	}

	// apply presets
	for _, presetID := range e.Presets {
		preset, _ := ps.Get(presetID)
		if err = preset.FillPublicFieldsFromPreset(e, "Presets"); err != nil {
			err = fmt.Errorf("Error resolving content: %v", err)
			return
		}
	}

	// ensure coordinate sorting is correct
	if e.X > e.X2 {
		e.X, e.X2 = e.X2, e.X
	}
	if e.Y > e.Y2 {
		e.Y, e.Y2 = e.Y2, e.Y
	}

	return nil
}

// generateOutput generates the output for this object.
func (e *multiline) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	valueArray := getMultiValue(e.Value, as)
	if len(valueArray) == 0 {
		return nil // nothing to do here...
	}

	if len(valueArray) > e.Lines {
		return fmt.Errorf("Error generating content output: Current multiline content has a maxium of %v lines, but %v input lines were provided", e.Lines, len(valueArray))
	}

	for idx, text := range valueArray {
		if !utils.IsSet(text) {
			continue
		}

		x, y, x2, y2 := e.getLineCoords(idx+1)
		s.AddTextCell(e.Canvas, x, y, x2, y2, e.Font, e.Fontsize, e.Align, text, true)
	}

	return nil
}

// deepCopy creates a deep copy of this entry.
func (e *multiline) deepCopy() Entry {
	copy := *e
	copy.Presets = append(make([]string, 0), e.Presets...)

	return &copy
}

func (e *multiline) getLineCoords(line int) (x, y, x2, y2 float64) {
	utils.Assert(line > 0 && line <= e.Lines, "Should only query for valid lines")

	lineHeight := (e.Y2 - e.Y) / float64(e.Lines)
	y = e.Y + (lineHeight * float64(line-1))

	return e.X, y, e.X2, y + lineHeight
}
