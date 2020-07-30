package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ContentEntry is an interface for the content. D'oh!
type ContentEntry interface {
	ID() string
	Type() string
	ExampleValue() string
	UsageExample() string
	IsValid() (err error)
	Describe(verbose bool) (result string)
	Resolve(ps PresetStore) (resolvedCI ContentEntry, err error)
	GenerateOutput(s *Stamp, value *string) (err error)
}

// NewContentEntry creates a new content entry object for the provided ContentData object.
func NewContentEntry(id string, data ContentData) (ce ContentEntry, err error) {
	switch data.Type {
	case "textCell":
		return NewContentTextCell(id, data)
	case "societyId":
		return NewContentSocietyID(id, data)
	case "":
		return nil, fmt.Errorf("No content type provided")
	default:
		return nil, fmt.Errorf("Unknown content type: '%v'", data.Type)
	}
}

// ---------------------------------------------------------------------------------

// ContentTextCell is the final type to implement textCells.
// TODO switch to pointers to distinguish between unset values and zero values?
type ContentTextCell struct {
	id           string
	description  string
	exampleValue string
	presets      []string

	X1, Y1   float64
	X2, Y2   float64
	Font     string
	Fontsize float64
	Align    string
}

// NewContentTextCell will return a content object that represents a text cell
func NewContentTextCell(id string, data ContentData) (tc ContentTextCell, err error) {
	// TODO return error for values that are set here besides the required ones

	tc.id = id
	tc.description = data.Desc
	tc.exampleValue = data.Example
	tc.presets = data.Presets

	tc.X1 = data.X1
	tc.Y1 = data.Y1
	tc.X2 = data.X2
	tc.Y2 = data.Y2
	tc.Font = data.Font
	tc.Fontsize = data.Fontsize
	tc.Align = data.Align

	return tc, nil
}

// ID returns the content objects ID
func (tc ContentTextCell) ID() (id string) {
	return tc.id
}

// Type returns the (hardcoded) type for this type of content
func (tc ContentTextCell) Type() (contentType string) {
	return "textCell"
}

// ExampleValue returns the example value provided for this content object. If
// no example was provided, an empty string is returned.
func (tc ContentTextCell) ExampleValue() (exampleValue string) {
	return tc.exampleValue
}

// UsageExample retuns an example call on how this content can be invoked
// from the command line.
func (tc ContentTextCell) UsageExample() (result string) {
	return genericContentUsageExample(tc.id, tc.exampleValue)
}

// IsValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (tc ContentTextCell) IsValid() (err error) {
	// TODO return error for negative numbers
	// TODO after switching to percent: Return errors for values x < 0 || x > 100
	// TODO switch to not checking all fields for being set, as values of 0 might be valid as well
	return CheckThatAllExportedFieldsAreSet(tc)
}

// Describe returns a textual description of the current content object
func (tc ContentTextCell) Describe(verbose bool) (result string) {
	var sb strings.Builder

	var description string
	if IsSet(tc.description) {
		description = tc.description
	} else {
		description = "No description available"
	}

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v", tc.id, description)
	} else {
		fmt.Fprintf(&sb, "- %v\n", tc.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", description)
		fmt.Fprintf(&sb, "\tType: %v\n", tc.Type())
		fmt.Fprintf(&sb, "\tExample: %v", tc.UsageExample())
	}

	return sb.String()
}

// Resolve the presets for this content object.
func (tc ContentTextCell) Resolve(ps PresetStore) (resolvedCI ContentEntry, err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(tc.presets...); err != nil {
		err = fmt.Errorf("Error resolving content '%v': %v", tc.ID(), err)
		return
	}

	for _, presetID := range tc.presets {
		preset, _ := ps.Get(presetID)
		AddMissingValues(&tc, preset)
	}

	return tc, nil
}

// GenerateOutput generates the output for this textCell object.
func (tc ContentTextCell) GenerateOutput(s *Stamp, value *string) (err error) {
	err = tc.IsValid()
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("No input value provided")
	}

	x, y, w, h := s.getXYWHasPt(tc.X1, tc.Y1, tc.X2, tc.Y2)

	s.AddTextCell(x, y, w, h, tc.Font, tc.Fontsize, tc.Align, *value)

	return nil
}

// ---------------------------------------------------------------------------------

// ContentSocietyID is the final type to implement societyIDs.
type ContentSocietyID struct {
	id           string
	description  string
	exampleValue string
	presets      []string

	X1, Y1   float64
	X2, Y2   float64
	XPivot   float64
	Font     string
	Fontsize float64
}

// NewContentSocietyID will return a content object that represents a society ID
func NewContentSocietyID(id string, data ContentData) (si ContentSocietyID, err error) {
	si.id = id
	si.description = data.Desc
	si.exampleValue = data.Example
	si.presets = data.Presets

	si.X1 = data.X1
	si.Y1 = data.Y1
	si.X2 = data.X2
	si.Y2 = data.Y2
	si.XPivot = data.XPivot
	si.Font = data.Font
	si.Fontsize = data.Fontsize

	return si, nil
}

// ID returns the content objects ID
func (si ContentSocietyID) ID() (id string) {
	return si.id
}

// Type returns the (hardcoded) type for this type of content
func (si ContentSocietyID) Type() (contentType string) {
	return "societyId"
}

// ExampleValue returns the example value provided for this content object. If
// no example was provided, an empty string is returned.
func (si ContentSocietyID) ExampleValue() (exampleValue string) {
	return si.exampleValue
}

// UsageExample retuns an example call on how this content can be invoked
// from the command line.
func (si ContentSocietyID) UsageExample() (result string) {
	return genericContentUsageExample(si.id, si.exampleValue)
}

// IsValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (si ContentSocietyID) IsValid() (err error) {
	if err = CheckThatAllExportedFieldsAreSet(si); err != nil {
		return err
	}

	x, _, w, _ := getXYWH(si.X1, si.Y1, si.X2, si.Y2)
	if si.XPivot <= x || si.XPivot >= (x+w) {
		return fmt.Errorf("xpivot value must lie between x1 and x2: %v < %v < %v", si.X1, si.XPivot, si.X2)
	}

	return nil
}

// Describe returns a textual description of the current content object
func (si ContentSocietyID) Describe(verbose bool) (result string) {
	var sb strings.Builder

	var description string
	if IsSet(si.description) {
		description = si.description
	} else {
		description = "No description available"
	}

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v", si.id, description)
	} else {
		fmt.Fprintf(&sb, "- %v\n", si.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", description)
		fmt.Fprintf(&sb, "\tType: %v\n", si.Type())
		fmt.Fprintf(&sb, "\tExample: %v", si.UsageExample())
	}

	return sb.String()
}

// Resolve the presets for this content object.
func (si ContentSocietyID) Resolve(ps PresetStore) (resolvedCI ContentEntry, err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(si.presets...); err != nil {
		err = fmt.Errorf("Error resolving content '%v': %v", si.ID(), err)
		return
	}

	for _, presetID := range si.presets {
		preset, _ := ps.Get(presetID)
		AddMissingValues(&si, preset)
	}

	return si, nil
}

// GenerateOutput generates the output for this textCell object.
func (si ContentSocietyID) GenerateOutput(s *Stamp, value *string) (err error) {
	err = si.IsValid()
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("No input value provided")
	}

	// check and split up provided society id value
	societyID := regexSocietyID.FindStringSubmatch(*value)
	if len(societyID) == 0 {
		return fmt.Errorf("Provided society ID does not follow the pattern '<player_id>-<char_id>': '%v'", *value)
	}
	Assert(len(societyID) == 3, "Should contain the matching text plus the capturing groups")
	playerID := societyID[1]
	charID := societyID[2]

	dash := " - "
	dashWidth, _ := s.ptToPct(s.GetStringWidth(dash, si.Font, "", si.Fontsize), 0.0)

	// draw white rectangle for (nearly) whole area to blank out existing dash
	// this is currently kind of fiddly and hackish... if we blank out the
	// complete area, then the bottom line may be gone as well, which I do not like...
	x, y, w, h := s.getXYWHasPt(si.X1, si.Y1, si.X2, si.Y2)
	yOffset := 1.0
	s.DrawRectangle(x, y-yOffset, w, h-yOffset, "F", 255, 255, 255)

	// player id
	x, y, w, h = s.getXYWHasPt(si.X1, si.Y1, si.XPivot-(dashWidth/2.0), si.Y2)
	s.AddTextCell(x, y, w, h, si.Font, si.Fontsize, "RB", playerID)

	// dash
	x, y, w, h = s.getXYWHasPt(si.XPivot-(dashWidth/2), si.Y1, si.XPivot+(dashWidth/2), si.Y2)
	s.AddTextCell(x, y, w, h, si.Font, si.Fontsize, "CB", dash)

	// char id
	x, y, w, h = s.getXYWHasPt(si.XPivot+(dashWidth/2.0), si.Y1, si.X2, si.Y2)
	s.AddTextCell(x, y, w, h, si.Font, si.Fontsize, "LB", charID)

	return nil
}

// ---------------------------------------------------------------------------------

// AddMissingValues iterates over the exported fields of the source object. For each
// such fields it checks whether the target object contains a field with the same
// name. If that is the case and if the target field does not yet have a value set,
// then the value from the source object is copied over.
func AddMissingValues(target interface{}, source interface{}, ignoredFields ...string) {
	Assert(reflect.ValueOf(source).Kind() == reflect.Struct, "Can only process structs as source")
	Assert(reflect.ValueOf(target).Kind() == reflect.Ptr, "Target argument must be passed by ptr, as we modify it")
	Assert(reflect.ValueOf(target).Elem().Kind() == reflect.Struct, "Can only process structs as target")

	vSrc := reflect.ValueOf(source)
	vDst := reflect.ValueOf(target).Elem()

	for i := 0; i < vDst.NumField(); i++ {
		fieldDst := vDst.Field(i)
		fieldName := vDst.Type().Field(i).Name

		// Ignore the Presets field, as we do not want to take over values for this.
		if Contains(ignoredFields, fieldName) { // especially filter out "Presets" and "ID"
			continue
		}

		// take care to skip unexported fields
		if !fieldDst.CanSet() {
			continue
		}

		fieldSrc := vSrc.FieldByName(fieldName)

		// skip target fields that do not exist on source side side
		if !fieldSrc.IsValid() {
			continue
		}

		if fieldDst.IsZero() && !fieldSrc.IsZero() {
			switch fieldDst.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldDst.Set(fieldSrc)
			default:
				panic(fmt.Sprintf("Unsupported datat type '%v' in struct, update function 'AddMissingValuesFrom()'", fieldDst.Kind()))
			}
		}
	}
}

// CheckThatAllExportedFieldsAreSet returns an error if at least one exported field
// in the passed structure is not set.
func CheckThatAllExportedFieldsAreSet(obj interface{}) (err error) {
	oVal := reflect.ValueOf(obj)
	Assert(oVal.Kind() == reflect.Struct, "Can only work on structs")

	unsetFields := make([]string, 0)
	for idx := 0; idx < oVal.NumField(); idx++ {
		fieldVal := oVal.Field(idx)

		// skip unexported fields
		if !IsExported(fieldVal) {
			continue
		}

		if !IsSet(fieldVal.Interface()) {
			fieldName := reflect.TypeOf(obj).Field(idx).Name
			unsetFields = append(unsetFields, fieldName)
		}
	}

	if len(unsetFields) > 0 {
		return fmt.Errorf("Missing value for the following fields: %v", unsetFields)
	}

	return nil
}

// genericContentUsageExample returns an example call for the current provided
// values. If no example value was provided, then a string containing
// "Not available" is returned instead.
func genericContentUsageExample(id, exampleValue string) (result string) {
	if !IsSet(exampleValue) {
		return fmt.Sprintf("Not available")
	}
	return fmt.Sprintf("%v=%v", id, QuoteStringIfRequired(exampleValue))
}
