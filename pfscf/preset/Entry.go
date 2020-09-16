package preset

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Entry represents an entry in the 'preset' section
type Entry struct {
	id      string
	presets []string
	values  map[string]interface{}
}

func newEntry() (e Entry) {
	e = Entry{
		presets: make([]string, 0),
		values:  make(map[string]interface{}, 0),
	}
	return
}

// ID returns the ID of this entry
func (entry *Entry) ID() string {
	return entry.id
}

// Get returns the value matching the provided key.
func (entry *Entry) Get(key string) (val interface{}, exists bool) {
	val, exists = entry.values[key]
	return
}

// UnmarshalYAML unmarshals a Parameter Store
func (entry *Entry) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	// contruct target
	*entry = newEntry()

	// unmarshal list of presets first
	type entryPresetsYAML struct{ Presets *[]string }
	epy := entryPresetsYAML{Presets: &entry.presets}
	if err = unmarshal(&epy); err != nil {
		return fmt.Errorf("Field 'presets' does not have expected list format")
	}

	// unmarshal everything else, remove the "presets" entry afterwards
	if err = unmarshal(&entry.values); err != nil {
		return err
	}
	delete(entry.values, "presets") // already included in separate "presets" array

	// remove empty values
	for key, value := range entry.values {
		// do not use IsSet() here. This will not work with any preset that is "", 0, 0.0, false, etc.
		// Instead check whether the value has a type. If it has a type, then we know that something
		// was provided, i.e. the field was not just empty.

		if reflect.TypeOf(value) == nil {
			delete(entry.values, key)
		}
	}

	// ensure that we only have supported types in here
	for key, value := range entry.values {
		switch kind := reflect.TypeOf(value).Kind(); kind {
		case reflect.String:
		case reflect.Float64:
		case reflect.Int:
		//case reflect.Bool: // comment in when bool is required somewhere
		default:
			return fmt.Errorf("Field '%v' in preset entry has unsupported type '%v'", key, kind)
		}

	}

	return nil
}

// doesNotContradict checks whether the provided objects are contradicting
// or not. They are not contradicting if all fields that they have in common
// have the same type and value.
// One exception to this is the "presets" list which is ignored here.
func (entry *Entry) doesNotContradict(other Entry) (err error) {
	for id, eValue := range entry.values {
		oValue, exists := other.values[id]
		if !exists {
			continue
		}

		eKind := reflect.TypeOf(eValue).Kind()
		oKind := reflect.TypeOf(oValue).Kind()
		if eKind != oKind {
			return fmt.Errorf("Contradicting types for field '%v':\n- '%v': %v\n- '%v': %v", id, entry.id, eKind, other.id, oKind)
		}

		if eValue != oValue {
			return fmt.Errorf("Contradicting values for field '%v':\n- '%v': %v\n- '%v': %v", id, entry.id, eValue, other.id, oValue)
		}
	}

	return nil
}

func (entry *Entry) deepCopy() *Entry {
	copy := newEntry()
	copy.id = entry.id
	for _, preset := range entry.presets {
		copy.presets = append(copy.presets, preset)
	}
	for key, value := range entry.values {
		copy.values[key] = value
	}

	return &copy
}

func (entry *Entry) inheritFrom(other Entry) {
	for otherKey, otherValue := range other.values {
		if _, exists := entry.values[otherKey]; !exists {
			entry.values[otherKey] = otherValue
		}
	}
}

// FillPublicFieldsFromPreset is a function that will do exactly as the name says
func (entry *Entry) FillPublicFieldsFromPreset(target interface{}, ignoredFields ...string) (err error) {
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
		presetVal, exists := entry.Get(lowerFieldName) // field names in presets map should be all lowercase
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

		return fmt.Errorf("Error while applying preset '%v:%v' to content, types do not match: Preset has '%v', content wants '%v'", entry.ID(), lowerFieldName, vPreset.Kind(), fieldDst.Kind())
	}

	return nil
}
