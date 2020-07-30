package main

import "testing"

func TestPresetEntry_IsNotContradictingWith(t *testing.T) {
	var err error

	cdEmpty := ContentData{}
	peEmpty := NewPresetEntry("idEmpty", cdEmpty)

	cdAllSet := getContentDataWithDummyData(t, "type")
	peAllSet := NewPresetEntry("idAllSet", cdAllSet)

	t.Run("no self-contradiction", func(t *testing.T) {
		// a given CE with values should not contradict itself
		err = peAllSet.IsNotContradictingWith(peAllSet)
		expectNoError(t, err)
	})

	t.Run("empty contradicts nothing", func(t *testing.T) {
		// a given CE with no values should contradict nothing
		err = peEmpty.IsNotContradictingWith(peEmpty)
		expectNoError(t, err)
		err = peAllSet.IsNotContradictingWith(peEmpty)
		expectNoError(t, err)
		err = peEmpty.IsNotContradictingWith(peAllSet)
		expectNoError(t, err)
	})

	t.Run("non-overlapping", func(t *testing.T) {
		// Have two partly-set objects with non-overlapping content
		cdLeft := ContentData{X1: 1.0, Desc: "desc"}
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := ContentData{X2: 2.0, Font: "font"}
		peRight := NewPresetEntry("idRight", cdRight)
		err = peLeft.IsNotContradictingWith(peRight)
		expectNoError(t, err)
	})

	t.Run("conflicting string attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Font = cdLeft.Font + "foo" // <= conflicting data
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewPresetEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		expectError(t, err)
	})

	t.Run("conflicting float64 attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Fontsize = cdLeft.Fontsize + 1.0 // <= conflicting data
		peLeft := NewPresetEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewPresetEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		expectError(t, err)
	})
}
