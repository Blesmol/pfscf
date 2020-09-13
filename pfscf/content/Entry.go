package content

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	regexParamValue = regexp.MustCompile(`^\s*param:\s*(\S*)$`)
)

// Entry is an interface for the content. D'oh!
type Entry interface {
	isValid(*param.Store) (err error)
	resolve(ps preset.Store) (err error)
	generateOutput(s *stamp.Stamp, as *args.Store) (err error)
	deepCopy() Entry
}

// TODO now with no ID we should print all fields of the respective entry instead (don't forget the type)
func contentValErr(ce Entry, errIn error) (errOut error) {
	return fmt.Errorf("Error validating content: %v; complete content entry is: %v", errIn, ce)
}

func genericFieldsCheck(obj interface{}, isOk func(interface{}) bool, fieldNames ...string) (err error) {
	oVal := reflect.ValueOf(obj)
	if oVal.Kind() == reflect.Ptr {
		oVal = oVal.Elem()
	}
	utils.Assert(oVal.Kind() == reflect.Struct, "Can only work on structs or pointers to structs")

	errFields := make([]string, 0)
	for _, fieldName := range fieldNames {
		fieldVal := oVal.FieldByName(fieldName)
		utils.Assert(fieldVal.IsValid(), fmt.Sprintf("No field with name '%v' found in struct of type '%T'", fieldName, obj))

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
	err = genericFieldsCheck(obj, utils.IsSet, fieldNames...)
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

func fillPublicFieldsFromPreset(target interface{}, pe *preset.Entry, ignoredFields ...string) (err error) {
	// assumption: target is pointer to struct
	utils.Assert(reflect.ValueOf(target).Kind() == reflect.Ptr, "Target argument must be passed by ptr, as we modify it")
	utils.Assert(reflect.ValueOf(target).Elem().Kind() == reflect.Struct, "Can only process structs as target")

	vDst := reflect.ValueOf(target).Elem()

	for i := 0; i < vDst.NumField(); i++ {
		fieldDst := vDst.Field(i)
		fieldName := vDst.Type().Field(i).Name

		// Ignore the Presets field, as we do not want to take over values for this.
		if utils.Contains(ignoredFields, fieldName) { // useful for filtering out things like "Presets"
			continue
		}

		// take care to skip unexported fields or fields that already have a value
		if !fieldDst.CanSet() || !fieldDst.IsZero() {
			continue
		}

		lowerFieldName := strings.ToLower(fieldName)
		presetVal, exists := pe.Get(lowerFieldName) // field names in presets map should be all lowercase
		if !exists {
			continue
		}
		vPreset := reflect.ValueOf(presetVal)

		// trivial case: equal kinds.
		if fieldDst.Kind() == vPreset.Kind() {
			fieldDst.Set(reflect.ValueOf(presetVal))
			continue
		}

		// try several conversions... hooray
		switch fieldDst.Kind() {
		case reflect.Float64:
			switch vPreset.Kind() {
			case reflect.Int:
				fieldDst.SetFloat(float64(vPreset.Int()))
				continue
			}
		default:
		}

		return fmt.Errorf("Error while applying preset '%v:%v' to content, types do not match: Preset has '%v', content wants '%v'", pe.ID(), lowerFieldName, vPreset.Kind(), fieldDst.Kind())
	}

	return nil
}

// getValue returns the value that should be used for the current content.
func getValue(valueField string, as *args.Store) (result *string) {
	// No input? No result!
	if !utils.IsSet(valueField) {
		return nil
	}

	// check whether a parameter reference was provided
	paramName := regexParamValue.FindStringSubmatch(valueField)
	if len(paramName) > 0 {
		utils.Assert(len(paramName) == 2, "Should contain the matching text plus a single capturing group")
		argValue, exists := as.Get(paramName[1])
		if exists {
			return &argValue
		}
		return nil
	}

	// else assume that provided value was a static text
	return &valueField
}
