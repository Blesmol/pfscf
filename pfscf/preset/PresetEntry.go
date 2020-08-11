package preset

import (
	"fmt"
	"reflect"

	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

// Entry represents an entry in the 'preset' section
type Entry struct {
	id      string
	presets []string

	X1, Y1   float64
	X2, Y2   float64
	XPivot   float64
	Font     string
	Fontsize float64
	Align    string
	Color    string
}

// NewEntry create a new PresetEntry object.
// TODO Throw error in case of unused fields from ContentData that are set.
func NewEntry(id string, data yaml.ContentData) (pe Entry) {
	utils.Assert(utils.IsSet(id), "ID should always be present here")
	pe = Entry{
		id:      id,
		presets: data.Presets,

		X1:       data.X1,
		Y1:       data.Y1,
		X2:       data.X2,
		Y2:       data.Y2,
		XPivot:   data.XPivot,
		Font:     data.Font,
		Fontsize: data.Fontsize,
		Align:    data.Align,
		Color:    data.Color,
	}
	return
}

// IsNotContradictingWith checks if the provided ContentEntry objects are
// contradicting or not. They are not contradicting if all values that are set
// (i.e. contain a non-zero value) within the objects contain the same value.
// One exception to this is the "Presets" list, which is ignored here.
func (pe Entry) IsNotContradictingWith(other Entry) (err error) {
	vLeft := reflect.ValueOf(pe)
	vRight := reflect.ValueOf(other)

	for i := 0; i < vLeft.NumField(); i++ {
		fieldLeft := vLeft.Field(i)
		fieldRight := vRight.Field(i)
		fieldName := vLeft.Type().Field(i).Name

		if !utils.IsExported(fieldLeft) {
			continue // skip non-exported fields
		}

		if fieldLeft.IsZero() || fieldRight.IsZero() {
			continue
		}
		if fieldLeft.Interface() != fieldRight.Interface() {
			return fmt.Errorf("Contradicting data for field '%v':\n- '%v': %v\n- '%v': %v", fieldName, pe.id, fieldLeft.Interface(), other.id, fieldRight.Interface())
		}
	}

	return nil
}
