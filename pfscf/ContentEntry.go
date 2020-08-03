package main

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	regexSocietyID = regexp.MustCompile(`^\s*(\d*)\s*-\s*(\d*)\s*$`)
)

// ContentEntry is an interface for the content. D'oh!
type ContentEntry interface {
	ID() string
	Type() string
	ExampleValue() string
	UsageExample() string
	//IsValid() (err error) // Currently not required as part of interface, might change later
	Describe(verbose bool) (result string)
	Resolve(ps PresetStore) (resolvedCI ContentEntry, err error)
	GenerateOutput(s *Stamp, as *ArgStore) (err error)
}

// NewContentEntry creates a new content entry object for the provided ContentData object.
func NewContentEntry(id string, data ContentData) (ce ContentEntry, err error) {
	switch data.Type {
	case "textCell":
		return NewContentTextCell(id, data)
	case "societyId":
		return NewContentSocietyID(id, data)
	case "rectangle":
		return NewContentRectangle(id, data)
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
func NewContentTextCell(id string, data ContentData) (ce ContentTextCell, err error) {
	// TODO return error for values that are set here besides the required ones

	ce.id = id
	ce.description = data.Desc
	ce.exampleValue = data.Example
	ce.presets = data.Presets

	ce.X1 = data.X1
	ce.Y1 = data.Y1
	ce.X2 = data.X2
	ce.Y2 = data.Y2
	ce.Font = data.Font
	ce.Fontsize = data.Fontsize
	ce.Align = data.Align

	return ce, nil
}

// ID returns the content objects ID
func (ce ContentTextCell) ID() (id string) {
	return ce.id
}

// Type returns the (hardcoded) type for this type of content
func (ce ContentTextCell) Type() (contentType string) {
	return "textCell"
}

// ExampleValue returns the example value provided for this content object. If
// no example was provided, an empty string is returned.
func (ce ContentTextCell) ExampleValue() (exampleValue string) {
	return ce.exampleValue
}

// UsageExample retuns an example call on how this content can be invoked
// from the command line.
func (ce ContentTextCell) UsageExample() (result string) {
	return genericContentUsageExample(ce.id, ce.exampleValue)
}

// IsValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce ContentTextCell) IsValid() (err error) {
	err = checkFieldsAreSet(ce, "Font", "Fontsize")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = checkFieldsAreInRange(ce, "X1", "Y1", "X2", "Y2")
	if err != nil {
		return contentValErr(ce, err)
	}

	if ce.X1 == ce.X2 {
		err = fmt.Errorf("Coordinates for X axis are equal")
		return contentValErr(ce, err)
	}

	if ce.Y1 == ce.Y2 {
		err = fmt.Errorf("Coordinates for Y axis are equal")
		return contentValErr(ce, err)
	}

	return nil
}

// Describe returns a textual description of the current content object
func (ce ContentTextCell) Describe(verbose bool) (result string) {
	var sb strings.Builder

	var description string
	if IsSet(ce.description) {
		description = ce.description
	} else {
		description = "No description available"
	}

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v", ce.id, description)
	} else {
		fmt.Fprintf(&sb, "- %v\n", ce.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", description)
		fmt.Fprintf(&sb, "\tType: %v\n", ce.Type())
		fmt.Fprintf(&sb, "\tExample: %v", ce.UsageExample())
	}

	return sb.String()
}

// Resolve the presets for this content object.
func (ce ContentTextCell) Resolve(ps PresetStore) (resolvedCI ContentEntry, err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.presets...); err != nil {
		err = fmt.Errorf("Error resolving content '%v': %v", ce.ID(), err)
		return
	}

	for _, presetID := range ce.presets {
		preset, _ := ps.Get(presetID)
		AddMissingValues(&ce, preset)
	}

	return ce, nil
}

// GenerateOutput generates the output for this textCell object.
func (ce ContentTextCell) GenerateOutput(s *Stamp, as *ArgStore) (err error) {
	err = ce.IsValid()
	if err != nil {
		return err
	}

	value, hasKey := as.Get(ce.ID())
	if !hasKey {
		return nil // nothing to do here...
	}

	y2 := s.DeriveY2(ce.Y1, ce.Y2, ce.Fontsize)
	s.AddTextCell(ce.X1, ce.Y1, ce.X2, y2, ce.Font, ce.Fontsize, ce.Align, value, true)

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
func (ce ContentSocietyID) ID() (id string) {
	return ce.id
}

// Type returns the (hardcoded) type for this type of content
func (ce ContentSocietyID) Type() (contentType string) {
	return "societyId"
}

// ExampleValue returns the example value provided for this content object. If
// no example was provided, an empty string is returned.
func (ce ContentSocietyID) ExampleValue() (exampleValue string) {
	return ce.exampleValue
}

// UsageExample retuns an example call on how this content can be invoked
// from the command line.
func (ce ContentSocietyID) UsageExample() (result string) {
	return genericContentUsageExample(ce.id, ce.exampleValue)
}

// IsValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce ContentSocietyID) IsValid() (err error) {
	err = checkFieldsAreSet(ce, "Font", "Fontsize")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = checkFieldsAreInRange(ce, "X1", "Y1", "X2", "Y2", "XPivot")
	if err != nil {
		return contentValErr(ce, err)
	}

	if ce.X1 == ce.X2 {
		err = fmt.Errorf("Coordinates for X axis are equal")
		return contentValErr(ce, err)
	}

	if ce.Y1 == ce.Y2 {
		err = fmt.Errorf("Coordinates for Y axis are equal")
		return contentValErr(ce, err)
	}

	x, _, w, _ := getXYWH(ce.X1, 0.0, ce.X2, 0.0)
	if ce.XPivot <= x || ce.XPivot >= (x+w) {
		return fmt.Errorf("xpivot value must lie between x1 and x2: %v < %v < %v", ce.X1, ce.XPivot, ce.X2)
	}

	return nil
}

// Describe returns a textual description of the current content object
func (ce ContentSocietyID) Describe(verbose bool) (result string) {
	var sb strings.Builder

	var description string
	if IsSet(ce.description) {
		description = ce.description
	} else {
		description = "No description available"
	}

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v", ce.id, description)
	} else {
		fmt.Fprintf(&sb, "- %v\n", ce.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", description)
		fmt.Fprintf(&sb, "\tType: %v\n", ce.Type())
		fmt.Fprintf(&sb, "\tExample: %v", ce.UsageExample())
	}

	return sb.String()
}

// Resolve the presets for this content object.
func (ce ContentSocietyID) Resolve(ps PresetStore) (resolvedCI ContentEntry, err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.presets...); err != nil {
		err = fmt.Errorf("Error resolving content '%v': %v", ce.ID(), err)
		return
	}

	for _, presetID := range ce.presets {
		preset, _ := ps.Get(presetID)
		AddMissingValues(&ce, preset)
	}

	return ce, nil
}

// GenerateOutput generates the output for this textCell object.
func (ce ContentSocietyID) GenerateOutput(s *Stamp, as *ArgStore) (err error) {
	Assert(as != nil, "No ArgStore provided")
	err = ce.IsValid()
	if err != nil {
		return err
	}

	value, hasKey := as.Get(ce.ID())
	if !hasKey {
		return nil // nothing to do here...
	}

	// check and split up provided society id value
	societyID := regexSocietyID.FindStringSubmatch(value)
	if len(societyID) == 0 {
		return fmt.Errorf("Provided society ID does not follow the pattern '<player_id>-<char_id>': '%v'", value)
	}
	Assert(len(societyID) == 3, "Should contain the matching text plus the capturing groups")
	playerID := societyID[1]
	charID := societyID[2]

	dash := " - "
	dashWidth := s.GetStringWidth(dash, ce.Font, "", ce.Fontsize)

	// draw white rectangle for (nearly) whole area to blank out existing dash
	// this is currently kind of fiddly and hackish... if we blank out the
	// complete area, then the bottom line may be gone as well, which I do not like...

	y2 := s.DeriveY2(ce.Y1, ce.Y2, ce.Fontsize)

	//y1, y2 := SortCoords(ce.Y1, ce.Y2)
	_, yOffset := s.ptToPct(0.0, 2.0)
	s.DrawRectangle(ce.X1, ce.Y1-yOffset, ce.X2, y2+yOffset, "F", 255, 255, 255)

	// player id
	s.AddTextCell(ce.X1, ce.Y1, ce.XPivot-(dashWidth/2.0), y2, ce.Font, ce.Fontsize, "RB", playerID, false)

	// dash
	s.AddTextCell(ce.XPivot-(dashWidth/2), ce.Y1, ce.XPivot+(dashWidth/2), y2, ce.Font, ce.Fontsize, "CB", dash, false)

	// char id
	s.AddTextCell(ce.XPivot+(dashWidth/2.0), ce.Y1, ce.X2, y2, ce.Font, ce.Fontsize, "LB", charID, false)

	return nil
}

// ---------------------------------------------------------------------------------

// ContentRectangle needs a description
type ContentRectangle struct {
	id           string
	description  string
	exampleValue string
	presets      []string

	X1, Y1 float64
	X2, Y2 float64
	Color  string
}

// NewContentRectangle will return a content object that represents a rectangle
func NewContentRectangle(id string, data ContentData) (ce ContentRectangle, err error) {
	ce.id = id
	ce.description = data.Desc
	ce.exampleValue = data.Example
	ce.presets = data.Presets

	ce.X1 = data.X1
	ce.Y1 = data.Y1
	ce.X2 = data.X2
	ce.Y2 = data.Y2
	ce.Color = data.Color

	return ce, nil
}

// ID returns the content objects ID
func (ce ContentRectangle) ID() (id string) {
	return ce.id
}

// Type returns the (hardcoded) type for this type of content
func (ce ContentRectangle) Type() (contentType string) {
	return "rectangle"
}

// ExampleValue returns the example value provided for this content object. If
// no example was provided, an empty string is returned.
func (ce ContentRectangle) ExampleValue() (exampleValue string) {
	return ce.exampleValue
}

// UsageExample retuns an example call on how this content can be invoked
// from the command line.
func (ce ContentRectangle) UsageExample() (result string) {
	return genericContentUsageExample(ce.id, ce.exampleValue)
}

// IsValid checks whether the current content object is valid and returns an
// error with details if the object is not valid.
func (ce ContentRectangle) IsValid() (err error) {
	err = checkFieldsAreSet(ce, "Color")
	if err != nil {
		return contentValErr(ce, err)
	}

	err = checkFieldsAreInRange(ce, "X1", "Y1", "X2", "Y2")
	if err != nil {
		return contentValErr(ce, err)
	}

	if ce.X1 == ce.X2 {
		err = fmt.Errorf("Coordinates for X axis are equal")
		return contentValErr(ce, err)
	}

	if ce.Y1 == ce.Y2 {
		err = fmt.Errorf("Coordinates for Y axis are equal")
		return contentValErr(ce, err)
	}

	if _, _, _, err = parseColor(ce.Color); err != nil {
		return contentValErr(ce, err)
	}

	return nil
}

// Describe returns a textual description of the current content object
func (ce ContentRectangle) Describe(verbose bool) (result string) {
	var sb strings.Builder

	var description string
	if IsSet(ce.description) {
		description = ce.description
	} else {
		description = "No description available"
	}

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v", ce.id, description)
	} else {
		fmt.Fprintf(&sb, "- %v\n", ce.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", description)
		fmt.Fprintf(&sb, "\tType: %v\n", ce.Type())
		fmt.Fprintf(&sb, "\tExample: %v", ce.UsageExample())
	}

	return sb.String()
}

// Resolve the presets for this content object.
func (ce ContentRectangle) Resolve(ps PresetStore) (resolvedCI ContentEntry, err error) {
	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(ce.presets...); err != nil {
		err = fmt.Errorf("Error resolving content '%v': %v", ce.ID(), err)
		return
	}

	for _, presetID := range ce.presets {
		preset, _ := ps.Get(presetID)
		AddMissingValues(&ce, preset)
	}

	return ce, nil
}

// GenerateOutput generates the output for this textCell object.
func (ce ContentRectangle) GenerateOutput(s *Stamp, as *ArgStore) (err error) {
	err = ce.IsValid()
	if err != nil {
		return err
	}

	_, hasKey := as.Get(ce.ID())
	if !hasKey {
		return nil // nothing to do here...
	}

	r, g, b, err := parseColor(ce.Color)
	if err != nil {
		return err
	}

	s.DrawRectangle(ce.X1, ce.Y1, ce.X2, ce.Y2, "F", r, g, b)

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
	if IsSet(colorCode) {
		colorCodeBytes := []byte(colorCode)
		decoded := make([]byte, hex.DecodedLen(len(colorCodeBytes)))
		_, err := hex.Decode(decoded, colorCodeBytes)
		Assert(err == nil, fmt.Sprintf("Valid input should have been guaranteed by regexp, but instead got error: %v", err))
		Assert(len(decoded) == 3, fmt.Sprintf("Number of resultint entries should be guaranteed by regexp, was %v instead", len(decoded)))

		r, g, b = int(decoded[0]), int(decoded[1]), int(decoded[2])
		return r, g, b, nil
	}

	return 0, 0, 0, fmt.Errorf("Unknown color: '%v'", color)
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

func genericFieldsCheck(obj interface{}, isOk func(interface{}) bool, fieldNames ...string) (err error) {
	oVal := reflect.ValueOf(obj)
	Assert(oVal.Kind() == reflect.Struct, "Can only work on structs")

	errFields := make([]string, 0)
	for _, fieldName := range fieldNames {
		fieldVal := oVal.FieldByName(fieldName)
		Assert(fieldVal.IsValid(), fmt.Sprintf("No field with name '%v' found in struct of type '%T'", fieldName, obj))

		if !isOk(fieldVal.Interface()) {
			errFields = append(errFields, fieldName)
		}
	}

	if len(errFields) > 0 {
		return fmt.Errorf("%v", errFields)
	}
	return nil
}

func checkFieldsAreSet(obj interface{}, fieldNames ...string) (err error) {
	err = genericFieldsCheck(obj, IsSet, fieldNames...)
	if err != nil {
		return fmt.Errorf("Missing values for the following fields: %v", err)
	}
	return nil
}

func checkFieldsAreInRange(obj interface{}, fieldNames ...string) (err error) {
	isOk := func(obj interface{}) bool {
		fObj := obj.(float64)
		return fObj >= 0.0 && fObj <= 100.0
	}
	err = genericFieldsCheck(obj, isOk, fieldNames...)
	if err != nil {
		return fmt.Errorf("Values for the following fields are out of range: %v", err)
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

func contentValErr(ce ContentEntry, errIn error) (errOut error) {
	return fmt.Errorf("Error validating content '%v': %v", ce.ID(), errIn)
}
