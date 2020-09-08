package preset

import (
	"reflect"
	"sort"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func Test_newEntry(t *testing.T) {
	entry := newEntry()

	test.ExpectNotSet(t, entry.id)
	test.ExpectEqual(t, len(entry.presets), 0)
	test.ExpectEqual(t, len(entry.values), 0)
}

func TestEntry_UnmarshalYAML(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		testData := []struct{ title, yamlInput, expectedError string }{
			{"malformed yaml", "foo: bar: foo", "mapping values are not allowed"},
			{"presets field not in list format", "presets: something", "does not have expected list format"},
			{"presets field with nested list", "presets: [ foo, [bar] ]", "does not have expected list format"},
			{"unsupported type: list", "foo: [ x ]", "unsupported type"},
			{"unsupported type: map", "foo:\n  a: b\n  b: a", "unsupported type"},
			{"unsupported type: bool", "bool: true", "unsupported type"},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			var entry Entry
			err := yaml.Unmarshal([]byte(tt.yamlInput), &entry)

			test.ExpectError(t, err, tt.expectedError)
		}
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("presets", func(t *testing.T) {
			testData := []struct {
				title, yamlInput string
				expectedList     []string
			}{
				{"missing field", "", []string{}},
				{"empty field", "presets:", []string{}},
				{"empty list", "presets: []", []string{}},
				{"list with one element", "presets: [p1]", []string{"p1"}},
				{"list with multiple elements", "presets: [p1, p1, p2]", []string{"p1", "p1", "p2"}},
			}

			for _, tt := range testData {
				t.Logf("Testing: %v", tt.title)

				entry := newEntry()
				err := yaml.Unmarshal([]byte(tt.yamlInput), &entry)

				test.ExpectNoError(t, err)

				// check list entries
				sort.Strings(entry.presets)
				sort.Strings(tt.expectedList)
				t.Logf("- Expected: '%v'", tt.expectedList)
				t.Logf("- Got:      '%v'", entry.presets)
				test.ExpectTrue(t, reflect.DeepEqual(entry.presets, tt.expectedList))

				// check that the presets field does not appear in the regular values
				_, exists := entry.values["presets"]
				test.ExpectFalse(t, exists)
			}
		})

		t.Run("values", func(t *testing.T) {
			testData := []struct {
				title, yamlInput string
				expectedMap      map[string]interface{}
			}{
				{"no fields", "", make(map[string]interface{}, 0)},
				{"single field", "foo: x", map[string]interface{}{"foo": "x"}},
				{"multiple fields", "foo: x\nbar: xx", map[string]interface{}{"foo": "x", "bar": "xx"}},
				{"empty field", "foo: x\nempty:\nbar: xx", map[string]interface{}{"foo": "x", "bar": "xx"}},
				{"different types without bool", "str: x\nint: 1\nfloat: 2.0", map[string]interface{}{"str": "x", "int": 1, "float": 2.0}},
				//{"bool values", "a: true\nb: y\nc: yes\nz: false\ny: n\nx: no", map[string]interface{}{"a": true, "b": true, "c": true, "z": false, "y": false, "x": false}},
			}

			for _, tt := range testData {
				t.Logf("Testing: %v", tt.title)

				entry := newEntry()
				err := yaml.Unmarshal([]byte(tt.yamlInput), &entry)

				test.ExpectNoError(t, err)

				// compare maps
				t.Logf("- Expected: '%v'", tt.expectedMap)
				t.Logf("- Got:      '%v'", entry.values)
				test.ExpectTrue(t, reflect.DeepEqual(entry.values, tt.expectedMap))

				// check that no preset were included here
				test.ExpectEqual(t, len(entry.presets), 0)
			}
		})
	})
}

func TestEntry_doesNotContradict(t *testing.T) {
	entry := newEntry()
	entry.id = "entry"
	entry.values = map[string]interface{}{
		"string": "some text",
		"int":    1,
		"float":  2.0,
	}

	t.Run("errors", func(t *testing.T) {
		testData := []struct {
			title         string
			inputValues   map[string]interface{}
			expectedError string
		}{
			{"conflicting type: string", map[string]interface{}{"string": 1.0}, "Contradicting types"},
			{"conflicting type: int", map[string]interface{}{"int": "foo"}, "Contradicting types"},
			{"conflicting type: float", map[string]interface{}{"float": 1}, "Contradicting types"},
			{"conflicting value: string", map[string]interface{}{"string": "other"}, "Contradicting values"},
			{"conflicting value: int", map[string]interface{}{"int": 2}, "Contradicting values"},
			{"conflicting value: float", map[string]interface{}{"float": 3.0}, "Contradicting values"},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			other := newEntry()
			other.id = "other"
			other.values = tt.inputValues

			err := entry.doesNotContradict(other)
			test.ExpectError(t, err, "entry", "other", tt.expectedError)

			err = other.doesNotContradict(entry)
			test.ExpectError(t, err, "entry", "other", tt.expectedError)
		}
	})

	t.Run("valid", func(t *testing.T) {
		testData := []struct {
			title       string
			inputValues map[string]interface{}
		}{
			{"empty", map[string]interface{}{}},
			{"string", map[string]interface{}{"string": "some text"}},
			{"int", map[string]interface{}{"int": 1}},
			{"float", map[string]interface{}{"float": 2.0}},
			{"same as input", entry.values},
			{"non-overlapping", map[string]interface{}{"stringa": "other text", "inta": 3.0}},
			{"overlapping", map[string]interface{}{"string": "some text", "foo": "bar"}},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			other := newEntry()
			other.id = "other"
			other.values = tt.inputValues

			err := entry.doesNotContradict(other)
			test.ExpectNoError(t, err)

			err = other.doesNotContradict(entry)
			test.ExpectNoError(t, err)

			err = other.doesNotContradict(other)
			test.ExpectNoError(t, err)
		}
	})
}

func TestEntry_deepCopy(t *testing.T) {
	testText := "some text"
	otherTestText := "some other text"

	e1 := newEntry()
	e1.presets = append(e1.presets, testText)
	e1.values[testText] = testText

	e2 := e1.deepCopy()

	// check id
	e1.id = otherTestText
	test.ExpectNotEqual(t, e1.id, e2.id)

	// check presets
	e1.presets[0] = otherTestText
	test.ExpectNotEqual(t, e1.presets[0], e2.presets[0])

	// check values
	e1.values[testText] = otherTestText
	test.ExpectNotEqual(t, e1.values[testText], e2.values[testText])
}
