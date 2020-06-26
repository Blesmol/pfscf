package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

// ContentEntry is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentEntry struct {
	Type     string  // the type which this entry represents
	ID       string  // the ID or name of that concrete content entry
	Desc     string  // Description of this parameter
	X1, Y1   float64 // first set of coordinates
	X2, Y2   float64 // second set of coordinates
	Font     string  // the name of the font (if any) that should be used to display the content
	Fontsize float64 // size of the font in points
	Align    string  // Alignment of the content: L/C/R + T/M/B
	//Flags    *[]string
}

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	name    string
	content map[string]ContentEntry
}

// NewChronicleTemplate returns a new ChronicleTemplate object
func NewChronicleTemplate(name string) (c *ChronicleTemplate) {
	c = new(ChronicleTemplate)
	c.name = name
	c.content = make(map[string]ContentEntry)
	return c
}

// TODO #6 generic function for checking required fields in struct
// also output warnings for non-required fields

// GetTemplateByName returns the template object for the given name, or nil and
// an error object if no template with that name could be found. The template
// name is case-insensitive.
func GetTemplateByName(tmplName string) (yFile *YamlFile, err error) {
	// Keep it simple for the moment. Search in 'templates' subdir
	// for a file with cfgName as basename and 'yml' as file extension

	tmplBaseFilename := strings.ToLower(tmplName) + ".yml"
	tmplFilename := filepath.Join(GetExecutableDir(), "templates", tmplBaseFilename)

	yFile, err = GetYamlFile(tmplFilename)

	return yFile, err
}

// GetChronicleTemplate extracts, processes, and prepares the template
// information from a YamlFile object and puts it into a form
// that can be worked with.
func (yFile *YamlFile) GetChronicleTemplate() (cTmpl *ChronicleTemplate, err error) {
	cTmpl = NewChronicleTemplate("pfs2") // TODO remove hardcoded name

	// add content entries from yamlFile with name mapping into chronicleTemplate
	for _, val := range yFile.Content {
		Assert(IsSet(val.ID), "No ID provided!")
		id := val.ID
		if _, exists := cTmpl.content[id]; exists {
			return nil, fmt.Errorf("Duplicate ID found: '%v'", id)
		}

		val.applyDefaults(yFile.Default)
		cTmpl.content[id] = val
	}

	return cTmpl, nil
}

// applyDefaults takes another ContentEntry object and then sets each field in the first
// CE object that is nil or zero to the value that is provided in the second CE object.
func (ce *ContentEntry) applyDefaults(defaults ContentEntry) {
	// function has to be called on ptr, as we modify the original object
	Assert(ce != nil, "Provided ContentEntry object should always be valid")

	vCE := reflect.ValueOf(ce).Elem()
	vDef := reflect.ValueOf(defaults)
	for i := 0; i < vCE.NumField(); i++ {
		fieldCE := vCE.Field(i)
		fieldDef := vDef.Field(i)

		// only proceed with current field if there is a need to use the default value
		if !IsSet(fieldCE.Interface()) && IsSet(fieldDef.Interface()) {
			Assert(fieldCE.CanSet(), fmt.Sprintf("Field with index %v must be settable", i))
			switch fieldCE.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldCE.Set(fieldDef)
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

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleTemplate
func (cTmpl *ChronicleTemplate) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = cTmpl.content[key]
	return
}

// CheckValuesArePresent assumes that all provided arguments are pointers and
// checks that all of them have a value, i.e. do not equal nil. It will return
// an error for each nil argument.
func CheckValuesArePresent(args ...interface{}) (err error) {
	for _, arg := range args {
		Assert(reflect.TypeOf(arg).Kind() == reflect.Ptr, "Argument should be a pointer")

		if reflect.ValueOf(arg).IsNil() {
			return fmt.Errorf("Missing ")
		}
	}

	return nil
}
