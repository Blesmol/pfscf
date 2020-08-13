package preset

import (
	"testing"

	"gopkg.in/yaml.v2"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

func Test_newStore(t *testing.T) {
	s := newStore()
	test.ExpectEqual(t, len(s), 0)
}

func TestStore_addGet(t *testing.T) {
	store := newStore()

	_, exists := store.Get("id")
	test.ExpectFalse(t, exists)

	entry1 := newEntry()
	store.add("id1", &entry1)
	entry2 := newEntry()
	store.add("id2", &entry2)

	entry, exists := store.Get("id1")
	test.ExpectTrue(t, exists)
	test.ExpectEqual(t, entry.id, "id1")

	entry, exists = store.Get("id2")
	test.ExpectTrue(t, exists)
	test.ExpectEqual(t, entry.id, "id2")

	_, exists = store.Get("id")
	test.ExpectFalse(t, exists)
}

func TestStore_UnmarshalYAML(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		testData := []struct{ title, yamlInput, expectedError string }{
			{"malformed yaml", "foo: bar: foo", "mapping values are not allowed"},
			//{"duplicate id", "id1:\n  foo: bar\nid1:\n  foo: bar", "duplicate ID"},
			{"invalid entry", "id1:\n  foo: [bar: foo]", "unsupported type"},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			var store Store
			err := yaml.Unmarshal([]byte(tt.yamlInput), &store)

			test.ExpectError(t, err, tt.expectedError)
		}
	})

	t.Run("valid", func(t *testing.T) {
		testData := []struct {
			title, yamlInput string
			expectedIDs      []string
		}{
			{"empty list", "", []string{}},
			{"single entry without fields", "id1:", []string{}},
			{"single entry", "id1:\n  foo: bar", []string{"id1"}},
		}

		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			var store Store
			err := yaml.Unmarshal([]byte(tt.yamlInput), &store)

			test.ExpectNoError(t, err)
			test.ExpectEqual(t, len(tt.expectedIDs), len(store))
			for _, id := range tt.expectedIDs {
				_, exists := store[id]
				test.ExpectTrue(t, exists)
			}
		}
	})
}

func TestStore_InheritFrom(t *testing.T) {
	t.Run("batch tests", func(t *testing.T) {
		testData := []struct {
			title                             string
			sourceIDs, targetIDs, expectedIDs []string
		}{
			{"both empty", []string{}, []string{}, []string{}},
			{"empty source", []string{}, []string{"id1", "id2"}, []string{"id1", "id2"}},
			{"empty target", []string{"id1", "id2"}, []string{}, []string{"id1", "id2"}},
			{"identical", []string{"id1", "id2"}, []string{"id1", "id2"}, []string{"id1", "id2"}},
			{"overlapping", []string{"id1", "id2"}, []string{"id1"}, []string{"id1", "id2"}},
			{"non-overlapping", []string{"id1"}, []string{"id2"}, []string{"id1", "id2"}},
		}
		for _, tt := range testData {
			t.Logf("Testing: %v", tt.title)

			// construct source store
			sSource := newStore()
			for _, id := range tt.sourceIDs {
				e := newEntry()
				e.values["source"] = 1
				sSource.add(id, &e)
			}

			// construct target store
			sTarget := newStore()
			for _, id := range tt.targetIDs {
				e := newEntry()
				e.values["target"] = 1
				sTarget.add(id, &e)
			}

			sTarget.InheritFrom(sSource)

			test.ExpectEqual(t, len(sTarget), len(tt.expectedIDs))
			for _, expectedID := range tt.expectedIDs {
				entry, exists := sTarget.Get(expectedID)
				test.ExpectTrue(t, exists)

				if utils.Contains(tt.targetIDs, expectedID) {
					_, exists = entry.values["target"]
				} else {
					_, exists = entry.values["source"]
				}
				test.ExpectTrue(t, exists)
			}
		}
	})
}

func TestStore_PresetsAreNotContradicting(t *testing.T) {
	store := newStore()

	someValue := 1

	e1 := newEntry()
	e1.values["val1"] = someValue
	store.add("v1_", &e1)

	e2 := newEntry()
	e2.values["val1"] = someValue + 1
	store.add("v2_", &e2)

	e3 := newEntry()
	e3.values["val1"] = someValue
	e3.values["val2"] = someValue
	store.add("v11", &e3)

	e4 := newEntry()
	e4.values["val1"] = someValue
	e4.values["val2"] = someValue + 1
	store.add("v12", &e4)

	t.Run("errors", func(t *testing.T) {
		t.Run("contradicting entries", func(t *testing.T) {
			for _, tt := range []struct {
				title string
				IDs   []string
			}{
				{"contradicting in first two elements", []string{"v1_", "v2_"}},
				{"contradicting in later elements", []string{"v1_", "v11", "v12"}},
			} {
				t.Logf("Testing: %v", tt.title)
				err := store.PresetsAreNotContradicting(tt.IDs...)
				test.ExpectError(t, err, "Contradicting values")
			}
		})

		t.Run("non-existent presets", func(t *testing.T) {
			for _, tt := range []struct {
				title string
				IDs   []string
			}{
				{"only non-existant", []string{"nonExistant"}},
				{"non-existant + existing", []string{"nonExistant", "v1_"}},
				{"existing + non-existant", []string{"v1_", "nonExistant"}},
			} {
				t.Logf("Testing: %v", tt.title)
				err := store.PresetsAreNotContradicting(tt.IDs...)
				test.ExpectError(t, err, "nonExistant", "does not exist")
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		for _, tt := range []struct {
			title string
			IDs   []string
		}{
			{"empty list", []string{}},
			{"single element", []string{"v11"}},
			{"duplicate elements", []string{"v11", "v11"}},
			{"non-contradicting elements", []string{"v1_", "v11", "v1_", "v11"}},
		} {
			t.Logf("Checking preset '%v'", tt)
			err := store.PresetsAreNotContradicting(tt.IDs...)
			test.ExpectNoError(t, err)
		}
	})
}
