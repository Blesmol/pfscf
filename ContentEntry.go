package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ContentData is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentData struct {
	Type     string  // the type which this entry represents
	Desc     string  // Description of this parameter
	X1, Y1   float64 // first set of coordinates
	X2, Y2   float64 // second set of coordinates
	XPivot   float64 // pivot point on X axis
	Font     string  // the name of the font (if any) that should be used to display the content
	Fontsize float64 // size of the font in points
	Align    string  // Alignment of the content: L/C/R + T/M/B
	Example  string  // Example value to be displayed to users
	//Flags    *[]string
	id string // not read directly from the yaml file
}

// applyDefaults takes another ContentData object and then sets each field in the first
// CD object that is nil or zero to the value that is provided in the second CD object.
func (cd *ContentData) applyDefaults(other ContentData) {
	// function has to be called on ptr, as we modify the original object
	Assert(cd != nil, "Provided ContentData object should always be valid")

	vCD := reflect.ValueOf(cd).Elem()
	vOther := reflect.ValueOf(other)
	for i := 0; i < vCD.NumField(); i++ {
		fieldCD := vCD.Field(i)
		fieldOther := vOther.Field(i)

		// private fields cannot be set, so skip in such cases
		if !fieldCD.CanSet() {
			continue
		}

		// only proceed with current field if there is a need to use the default value
		if !IsSet(fieldCD.Interface()) && IsSet(fieldOther.Interface()) {
			Assert(fieldCD.CanSet(), fmt.Sprintf("Field with index %v must be settable", i))
			switch fieldCD.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldCD.Set(fieldOther)
			default:
				panic(fmt.Sprintf("Unsupported struct type '%v', update function", fieldCD.Kind()))
			}
		}
	}
}

// checkThatValuesArePresent takes a list of field names from the ContentData struct and checks
// that these fields neither point to a nil ptr nor that the values behind the pointers contain the
// corresponding types zero value.
func (cd ContentData) checkThatValuesArePresent(names ...string) (isValid bool, err error) {
	r := reflect.ValueOf(cd)

	for _, name := range names {
		field := r.FieldByName(name)
		Assert(field.IsValid(), fmt.Sprintf("ContentData does not contain a field with name '%v'", name))

		if !IsSet(field.Interface()) {
			return false, fmt.Errorf("ContentData object does not contain a value for field '%v'", name)
		}
	}
	return true, nil
}

// IsValid checks whether a ContentData object is valid. This means that it
// must contain type information, and depending on the type information
// a certain set of other fields must be set.
func (cd ContentData) IsValid() (isValid bool, err error) {
	// Type must be checked first, as we decide by that value on which fields to check
	isValid, err = cd.checkThatValuesArePresent("Type")
	if !isValid {
		return isValid, err
	}

	switch cd.Type {
	case "textCell":
		return cd.checkThatValuesArePresent("X1", "Y1", "X2", "Y2", "Font", "Fontsize", "Align")
	default:
		return false, fmt.Errorf("ContentData object contains unknown content type '%v'", cd.Type)
	}
}

// Describe describes a single ContentData object. It returns the
// description as a multi-line string
func (cd *ContentData) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", cd.id)
		if IsSet(cd.Desc) {
			fmt.Fprintf(&sb, ": %v", cd.Desc)
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", cd.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", cd.Desc)
		fmt.Fprintf(&sb, "\tType: %v\n", cd.Type)
		fmt.Fprintf(&sb, "\tExample: %v", cd.getExample())
	}

	return sb.String()
}

// getExample returns an example call for the current ContentData object.
// If this is not possible, e.g because no example value is included or
// this is in general not possible for the given type, then a string
// containing "Not available" is returned instead.
func (cd *ContentData) getExample() (result string) {
	switch cd.Type {
	case "textCell":
		if !IsSet(cd.Example) {
			return fmt.Sprintf("Not available")
		}
		return fmt.Sprintf("%v=%v", cd.id, QuoteStringIfRequired(cd.Example))
	default:
		panic("Unknown ContentData type")
	}
}
