package param

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func TestTextEntry(t *testing.T) {
	entry := textEntry{
		TheExample:     "some example",
		TheDescription: "some description",
	}

	test.ExpectEqual(t, entry.Example(), "some example")
	test.ExpectEqual(t, entry.Description(), "some description")

	test.ExpectNotSet(t, entry.ID())
	entry.setID("some id")
	test.ExpectEqual(t, entry.ID(), "some id")

	test.ExpectEqual(t, entry.Type(), "text")
}

func TestTextEntry_YAML(t *testing.T) {
	yamlInput := []byte(`
example: foo
description: 1.0
thedescription: bar
`)
	var entry textEntry

	err := yaml.Unmarshal(yamlInput, &entry)
	test.ExpectNoError(t, err)
	test.ExpectNotSet(t, entry.id)
	test.ExpectEqual(t, entry.TheExample, "foo")
	test.ExpectEqual(t, entry.TheDescription, "1.0")
}

func TestTextEntry_deepCopy(t *testing.T) {
	var e1 textEntry
	e1.id = "foo"
	e2 := e1.deepCopy().(*textEntry)
	e2.id = "bar"
	test.ExpectNotEqual(t, e1.id, e2.id)
}
