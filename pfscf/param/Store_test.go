package param

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func TestStore_AddGet(t *testing.T) {
	store := newStore()

	_, exists := store.Get("id")
	test.ExpectFalse(t, exists)

	store.add("id1", &textEntry{TheExample: "example1"})
	store.add("id2", &societyidEntry{TheExample: "example2"})

	entry, exists := store.Get("id1")
	test.ExpectTrue(t, exists)
	test.ExpectEqual(t, entry.ID(), "id1")
	test.ExpectEqual(t, entry.Example(), "example1")

	entry, exists = store.Get("id2")
	test.ExpectTrue(t, exists)
	test.ExpectEqual(t, entry.ID(), "id2")
	test.ExpectEqual(t, entry.Example(), "example2")

	_, exists = store.Get("id")
	test.ExpectFalse(t, exists)
}

func TestStore_UnmarshalYAML(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		yamlInput := []byte("id1:\n  type: invalid")

		var store Store
		err := yaml.Unmarshal(yamlInput, &store)

		test.ExpectError(t, err, "Unknown parameter type")
	})

	t.Run("valid", func(t *testing.T) {
		yamlInput := []byte(`
id1:
  type: text
  example: example1
id2:
  type: societyid
  example: example2
`)

		var store Store
		err := yaml.Unmarshal(yamlInput, &store)

		test.ExpectNoError(t, err)
		test.ExpectEqual(t, len(store), 2)

		entry, exists := store.Get("id1")
		test.ExpectTrue(t, exists)
		test.ExpectEqual(t, entry.Type(), typeText)
		test.ExpectEqual(t, entry.Example(), "example1")

		entry, exists = store.Get("id2")
		test.ExpectTrue(t, exists)
		test.ExpectEqual(t, entry.Type(), typeSocietyID)
		test.ExpectEqual(t, entry.Example(), "example2")
	})
}
