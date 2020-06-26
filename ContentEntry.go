package main

import (
	"fmt"
	"reflect"
)

// ContentEntry is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentEntry struct {
	Type     string  // the type which this entry represents
	Desc     string  // Description of this parameter
	X1, Y1   float64 // first set of coordinates
	X2, Y2   float64 // second set of coordinates
	Font     string  // the name of the font (if any) that should be used to display the content
	Fontsize float64 // size of the font in points
	Align    string  // Alignment of the content: L/C/R + T/M/B
	//Flags    *[]string
}

// applyDefaults takes another ContentEntry object and then sets each field in the first
// CE object that is nil or zero to the value that is provided in the second CE object.
func (ce *ContentEntry) applyDefaults(other ContentEntry) {
	// function has to be called on ptr, as we modify the original object
	Assert(ce != nil, "Provided ContentEntry object should always be valid")

	vCE := reflect.ValueOf(ce).Elem()
	vOther := reflect.ValueOf(other)
	for i := 0; i < vCE.NumField(); i++ {
		fieldCE := vCE.Field(i)
		fieldOther := vOther.Field(i)

		// only proceed with current field if there is a need to use the default value
		if !IsSet(fieldCE.Interface()) && IsSet(fieldOther.Interface()) {
			Assert(fieldCE.CanSet(), fmt.Sprintf("Field with index %v must be settable", i))
			switch fieldCE.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldCE.Set(fieldOther)
			default:
				panic(fmt.Sprintf("Unsupported struct type '%v', update function", fieldCE.Kind()))
			}
		}
	}
}

// checkThatValuesArePresent takes a list of field names from the ContentEntry struct and checks
// that these fields neither point to a nil ptr nor that the values behind the pointers contain the
// corresponding types zero value.
func (ce ContentEntry) checkThatValuesArePresent(names ...string) (isValid bool, err error) {
	r := reflect.ValueOf(ce)

	for _, name := range names {
		field := r.FieldByName(name)
		Assert(field.IsValid(), fmt.Sprintf("ContentEntry does not contain a field with name '%v'", name))

		if !IsSet(field.Interface()) {
			return false, fmt.Errorf("ContentEntry object does not contain a value for field '%v'", name)
		}
	}
	return true, nil
}

// IsValid checks whether a ContentEntry is valid. This means that it
// must contain type information, and depending on the type information
// a certain set of other fields must be set.
func (ce ContentEntry) IsValid() (isValid bool, err error) {
	// Type must be checked first, as we decide by that value on which fields to check
	isValid, err = ce.checkThatValuesArePresent("Type")
	if !isValid {
		return isValid, err
	}

	switch ce.Type {
	case "textCell":
		return ce.checkThatValuesArePresent("X1", "Y1", "X2", "Y2", "Font", "Fontsize", "Align")
	default:
		return false, fmt.Errorf("ContentEntry object contains unknown content type '%v'", ce.Type)
	}
}
