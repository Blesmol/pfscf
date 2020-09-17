package content

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeRectangle = "rectangle"
)

// rectangle needs a description
type rectangle struct {
	X, Y         float64
	X2, Y2       float64
	Color        string
	Transparency float64 // TODO convert to ptr
	Presets      []string
}

func newRectangle() *rectangle {
	var ce rectangle
	ce.Presets = make([]string, 0)
	return &ce
}

// isValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce *rectangle) isValid(paramStore *param.Store) (err error) {
	err = utils.CheckFieldsAreSet(ce, "Color")
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

	if _, _, _, err = parseColor(ce.Color); err != nil {
		return contentValErr(ce, err)
	}

	if ce.Transparency < 0.0 || ce.Transparency > 1.0 {
		err = fmt.Errorf("Transparency value outside of range 0.0 to 1.0: %v", ce.Transparency)
		return contentValErr(ce, err)
	}

	return nil
}

// resolve the presets for this content object.
func (ce *rectangle) resolve(ps preset.Store) (err error) {
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

	return nil
}

// generateOutput generates the output for this textCell object.
func (ce *rectangle) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	r, g, b, err := parseColor(ce.Color)
	if err != nil {
		return err
	}

	s.DrawRectangle(ce.X, ce.Y, ce.X2, ce.Y2, "F", 0, 0, 0, r, g, b, ce.Transparency)

	return nil
}

func parseColor(color string) (r, g, b int, err error) {
	regexHexColorCode := regexp.MustCompile(`^[0-9a-f]{6}$`)

	color = strings.ToLower(strings.TrimSpace(color))

	switch color {
	case "white":
		return 255, 255, 255, nil
	case "black":
		return 0, 0, 0, nil
	case "blue":
		return 0, 0, 255, nil
	case "red":
		return 255, 0, 0, nil
	case "green":
		return 0, 255, 0, nil
	}

	colorCode := regexHexColorCode.FindString(color)
	if utils.IsSet(colorCode) {
		colorCodeBytes := []byte(colorCode)
		decoded := make([]byte, hex.DecodedLen(len(colorCodeBytes)))
		_, err := hex.Decode(decoded, colorCodeBytes)
		utils.Assert(err == nil, fmt.Sprintf("Valid input should have been guaranteed by regexp, but instead got error: %v", err))
		utils.Assert(len(decoded) == 3, fmt.Sprintf("Number of resultint entries should be guaranteed by regexp, was %v instead", len(decoded)))

		r, g, b = int(decoded[0]), int(decoded[1]), int(decoded[2])
		return r, g, b, nil
	}

	return 0, 0, 0, fmt.Errorf("Unknown color: '%v'", color)
}

// deepCopy creates a deep copy of this entry.
// TODO create generic deep-copy function for public fields
func (ce *rectangle) deepCopy() Entry {
	copy := rectangle{
		X:            ce.X,
		Y:            ce.Y,
		X2:           ce.X2,
		Y2:           ce.Y2,
		Color:        ce.Color,
		Transparency: ce.Transparency,
	}
	for _, preset := range ce.Presets {
		copy.Presets = append(copy.Presets, preset)
	}

	return &copy
}
