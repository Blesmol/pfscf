package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ContentData is a generic data-holding struct with lots of fields
// to fit all the supported tpyes of Content

// ContentData is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentData struct {
	Type     string   // the type which this entry represents
	Desc     string   // Description of this parameter
	X1       float64  `yaml:"x"` // first x coordinate
	Y1       float64  `yaml:"y"` // first y coordinate
	X2, Y2   float64  // second set of coordinates
	XPivot   float64  // pivot point on X axis
	Font     string   // the name of the font (if any) that should be used to display the content
	Fontsize float64  // size of the font in points
	Align    string   // Alignment of the content: L/C/R + T/M/B
	Example  string   // Example value to be displayed to users
	Presets  []string // List of presets that should be applied on this ContentData / ContentEntry
	//Flags    *[]string
}

// ContentEntry wraps a ContentData object and adds some additional data.
// It also provides a bunch of functions to validate and access the data
// and to perform operations on it.
type ContentEntry struct {
	id   string
	data ContentData
}

// NewContentEntry creates a new ContentEntry object.
func NewContentEntry(id string, data ContentData) (ce ContentEntry) {
	Assert(IsSet(id), "ID should always be present here")
	ce.id = id
	ce.data = data
	return
}

// ID returns the id of this ContentEntry object
func (ce *ContentEntry) ID() (result string) {
	return ce.id
}

// Type returns the type value this ContentEntry object
func (ce *ContentEntry) Type() (result string) {
	return ce.data.Type
}

// Description returns the description of this ContentEntry object
func (ce *ContentEntry) Description() (result string) {
	return ce.data.Desc
}

// X1 returns the x1 value of this ContentEntry object
func (ce *ContentEntry) X1() (result float64) {
	return ce.data.X1
}

// Y1 returns the y1 value of this ContentEntry object
func (ce *ContentEntry) Y1() (result float64) {
	return ce.data.Y1
}

// X2 returns the x2 value of this ContentEntry object
func (ce *ContentEntry) X2() (result float64) {
	return ce.data.X2
}

// Y2 returns the y2 value of this ContentEntry object
func (ce *ContentEntry) Y2() (result float64) {
	return ce.data.Y2
}

// XPivot returns the xpivot of this ContentEntry object
func (ce *ContentEntry) XPivot() (result float64) {
	return ce.data.XPivot
}

// Font returns the font of this ContentEntry object
func (ce *ContentEntry) Font() (result string) {
	return ce.data.Font
}

// Fontsize returns of this ContentEntry object
func (ce *ContentEntry) Fontsize() (result float64) {
	return ce.data.Fontsize
}

// Align returns the alignment string of this ContentEntry object
func (ce *ContentEntry) Align() (result string) {
	return ce.data.Align
}

// Example returns the example of this ContentEntry object
func (ce *ContentEntry) Example() (result string) {
	return ce.data.Example
}

// Presets returns the list of presets set for this ContentEntry object
func (ce *ContentEntry) Presets() (result []string) {
	return ce.data.Presets
}

// CheckThatValuesArePresent takes a list of field names from the included ContentData struct and checks
// that these fields neither point to a nil ptr nor that the values behind the pointers contain the
// corresponding types zero value.
func (ce ContentEntry) CheckThatValuesArePresent(names ...string) (err error) {
	// TODO name all missing entries in error message
	r := reflect.ValueOf(ce.data)

	for _, fieldName := range names {
		field := r.FieldByName(fieldName)
		Assert(field.IsValid(), fmt.Sprintf("ContentData does not contain a field with name '%v'", fieldName))

		if !IsSet(field.Interface()) {
			return fmt.Errorf("ContentEntry '%v' does not contain a value for field '%v'", ce.ID(), fieldName)
		}
	}
	return nil
}

// IsValid checks whether a ContentEntry object is valid. This means that it
// must contain type information, and depending on the type information
// a certain set of other fields must be set.
func (ce ContentEntry) IsValid() (err error) {
	// Type must be checked first, as we decide by that value on which fields to check
	err = ce.CheckThatValuesArePresent("Type")

	if err == nil {
		switch ce.Type() {
		case "textCell":
			err = ce.CheckThatValuesArePresent("X1", "Y1", "X2", "Y2", "Font", "Fontsize", "Align")
		case "societyid":
			err = ce.CheckThatValuesArePresent("X1", "Y1", "X2", "Y2", "XPivot", "Font", "Fontsize")
		default:
			err = fmt.Errorf("Content has unknown type '%v'", ce.Type())
		}
	}

	if err != nil {
		return fmt.Errorf("Error validating content '%v': %v", ce.ID(), err)
	}
	return nil
}

// Describe describes a single ContentEntry object. It returns the
// description as a multi-line string
func (ce *ContentEntry) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", ce.ID())
		if IsSet(ce.Description()) {
			fmt.Fprintf(&sb, ": %v", ce.Description())
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", ce.ID())
		fmt.Fprintf(&sb, "\tDesc: %v\n", ce.Description())
		fmt.Fprintf(&sb, "\tType: %v\n", ce.Type())
		fmt.Fprintf(&sb, "\tExample: %v", ce.UsageExample())
	}

	return sb.String()
}

// UsageExample returns an example call for the current ContentEntry object.
// If this is not possible, e.g because no example value is included or
// this is in general not possible for the given type, then a string
// containing "Not available" is returned instead.
func (ce *ContentEntry) UsageExample() (result string) {
	switch ce.Type() {
	case "textCell":
		if !IsSet(ce.Example) {
			return fmt.Sprintf("Not available")
		}
		return fmt.Sprintf("%v=%v", ce.id, QuoteStringIfRequired(ce.Example()))
	default:
		panic("Unknown ContentEntry type")
	}
}

// IsNotContradictingWith checks if the provided ContentEntry objects are
// contradicting or not. They are not contradicting all values that are set
// (i.e. contain a non-zero value) within the objects contain the same value.
// One exception to this is the "Presets" list, which is ignored here.
func (ce ContentEntry) IsNotContradictingWith(other ContentEntry) (err error) {
	vLeft := reflect.ValueOf(ce.data)
	vRight := reflect.ValueOf(other.data)

	for i := 0; i < vLeft.NumField(); i++ {
		fieldLeft := vLeft.Field(i)
		fieldRight := vRight.Field(i)
		fieldName := vLeft.Type().Field(i).Name

		// Ignore the Presets field, as differences here are acceptable.
		// To be on the safe side wrt future changes, check the name instead of checking
		// whether this field is of kind struct.
		if fieldName == "Presets" {
			continue
		}

		if fieldLeft.IsZero() || fieldRight.IsZero() {
			continue
		}
		if fieldLeft.Interface() != fieldRight.Interface() {
			return fmt.Errorf("Contradicting data for field '%v':\n- '%v': %v\n- '%v': %v", fieldName, ce.ID(), fieldLeft.Interface(), other.ID(), fieldRight.Interface())
		}
	}

	return nil
}

// AddMissingValuesFrom wants to have a documentation
func (ce *ContentEntry) AddMissingValuesFrom(other *ContentEntry) {
	// TODO convert arg from ptr to value?
	vSrc := reflect.ValueOf(other.data)
	vDst := reflect.ValueOf(&ce.data).Elem() // go over pointer instead of value as we want to modify

	for i := 0; i < vDst.NumField(); i++ {
		fieldDst := vDst.Field(i)
		fieldSrc := vSrc.Field(i)
		fieldName := vSrc.Type().Field(i).Name

		// Ignore the Presets field, as we do not want to take over values for this.
		if fieldName == "Presets" {
			continue
		}

		if fieldDst.IsZero() && !fieldSrc.IsZero() {
			Assert(fieldDst.CanSet(), fmt.Sprintf("Field with index %v must be settable", i))

			switch fieldDst.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldDst.Set(fieldSrc)
			default:
				panic(fmt.Sprintf("Unsupported struct type '%v', update function 'AddMissingValuesFrom()'", fieldDst.Kind()))
			}
		}
	}
}
