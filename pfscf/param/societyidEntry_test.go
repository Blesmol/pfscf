package param

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func TestSocietyidEntry(t *testing.T) {
	entry := societyidEntry{
		TheExample:     "some example",
		TheDescription: "some description",
	}

	test.ExpectEqual(t, entry.Example(), "some example")
	test.ExpectEqual(t, entry.Description(), "some description")

	test.ExpectNotSet(t, entry.ID())
	entry.setID("some id")
	test.ExpectEqual(t, entry.ID(), "some id")

	test.ExpectEqual(t, entry.Type(), "societyid")
}

func TestSocietyidEntry_YAML(t *testing.T) {
	yamlInput := `
example: foo
description: 1.0
thedescription: bar
`
	var entry societyidEntry

	err := yaml.Unmarshal([]byte(yamlInput), &entry)
	test.ExpectNoError(t, err)
	test.ExpectNotSet(t, entry.id)
	test.ExpectEqual(t, entry.TheExample, "foo")
	test.ExpectEqual(t, entry.TheDescription, "1.0")
}

func TestSocietyidEntry_deepCopy(t *testing.T) {
	var e1 societyidEntry
	e1.id = "foo"
	e2 := e1.deepCopy().(*societyidEntry)
	e2.id = "bar"
	test.ExpectNotEqual(t, e1.id, e2.id)
}
