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
	tmplFilename := filepath.Join(GetTemplatesDir(), tmplBaseFilename)

	yFile, err = GetYamlFile(tmplFilename)

	return yFile, err
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
