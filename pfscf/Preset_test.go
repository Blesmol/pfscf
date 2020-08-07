package main

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

func TestPresetEntry_IsNotContradictingWith(t *testing.T) {
	var err error

	cdEmpty := yaml.ContentData{}
	peEmpty := NewPresetEntry("idEmpty", cdEmpty)

	cdAllSet := getContentDataWithDummyData(t, "type")
	peAllSet := NewPresetEntry("idAllSet", cdAllSet)

	t.Run("no self-contradiction", func(t *testing.T) {
		// a given CE with values should not contradict itself
		err = peAllSet.IsNotContradictingWith(peAllSet)
		test.ExpectNoError(t, err)
	})

	t.Run("empty contradicts nothing", func(t *testing.T) {
		// a given CE with no values should contradict nothing
		err = peEmpty.IsNotContradictingWith(peEmpty)
		test.ExpectNoError(t, err)
		err = peAllSet.IsNotContradictingWith(peEmpty)
		test.ExpectNoError(t, err)
		err = peEmpty.IsNotContradictingWith(peAllSet)
		test.ExpectNoError(t, err)
	})

	t.Run("non-overlapping", func(t *testing.T) {
		// Have two partly-set objects with non-overlapping content
		cdLeft := yaml.ContentData{X1: 1.0, Desc: "desc"}
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := yaml.ContentData{X2: 2.0, Font: "font"}
		peRight := NewPresetEntry("idRight", cdRight)
		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectNoError(t, err)
	})

	t.Run("conflicting string attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Font = cdLeft.Font + "foo" // <= conflicting data
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewPresetEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectError(t, err)
	})

	t.Run("conflicting float64 attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Fontsize = cdLeft.Fontsize + 1.0 // <= conflicting data
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewPresetEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectError(t, err)
	})
}
