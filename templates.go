package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

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
