package param

import (
	"fmt"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func TestEntryYAML_UnmarshalYAML(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		testData := []struct{ title, yamlInput, expectedError string }{
			{"malformed yaml", "foo: bar: foo", "mapping values are not allowed"},
			{"missing type", "example: foo", "Missing or empty 'type' field"},
			{"empty type", "type:", "Missing or empty 'type' field"},
			{"unknown type", "type: foobar", "Unknown parameter type: 'foobar'"},
			{"list instead of string 1", "type: [ foo bar ]", "cannot unmarshal"},
			{"list instead of string 2", "type: text\nexample: [ foo bar ]", "cannot unmarshal"},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			var entry entryYAML
			err := yaml.Unmarshal([]byte(tt.yamlInput), &entry)

			test.ExpectError(t, err, tt.expectedError)
		}
	})

	t.Run("valid", func(t *testing.T) {
		testData := []struct{ typeName string }{
			{typeText},
			{typeSocietyID},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.typeName)

			yamlInput := fmt.Sprintf("type: %v\nexample: some example", tt.typeName)

			var entry entryYAML
			err := yaml.Unmarshal([]byte(yamlInput), &entry)

			test.ExpectNoError(t, err)
			test.ExpectEqual(t, entry.e.Type(), tt.typeName)
			test.ExpectEqual(t, entry.e.Example(), "some example")
		}
	})
}
