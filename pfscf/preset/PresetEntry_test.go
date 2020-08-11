package preset

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

func getContentDataWithDummyData(t *testing.T, cdType string) (cd yaml.ContentData) {
	cd.Type = "Dummy, replaced below"
	cd.Desc = "Some Description"
	cd.X1 = 12.0
	cd.Y1 = 12.0
	cd.X2 = 24.0
	cd.Y2 = 24.0
	cd.XPivot = 15.0
	cd.Font = "Helvetica"
	cd.Fontsize = 14.0
	cd.Align = "LB"
	cd.Color = "green"
	cd.Example = "Some Example"
	cd.Presets = []string{"Some Preset"}

	test.ExpectAllExportedSet(t, cd) // to be sure that we also get all new fields

	// overvwrite type after the "expect..." check as cdType could be intentionally empty
	cd.Type = cdType

	return cd
}

func TestPresetEntry_IsNotContradictingWith(t *testing.T) {
	var err error

	cdEmpty := yaml.ContentData{}
	peEmpty := NewEntry("idEmpty", cdEmpty)

	cdAllSet := getContentDataWithDummyData(t, "type")
	peAllSet := NewEntry("idAllSet", cdAllSet)

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
		peLeft := NewEntry("idLeft", cdLeft)
		cdRight := yaml.ContentData{X2: 2.0, Font: "font"}
		peRight := NewEntry("idRight", cdRight)
		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectNoError(t, err)
	})

	t.Run("conflicting string attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Font = cdLeft.Font + "foo" // <= conflicting data
		peLeft := NewEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectError(t, err)
	})

	t.Run("conflicting float64 attribute", func(t *testing.T) {
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Fontsize = cdLeft.Fontsize + 1.0 // <= conflicting data
		peLeft := NewEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		peRight := NewEntry("idRight", cdRight)

		err = peLeft.IsNotContradictingWith(peRight)
		test.ExpectError(t, err)
	})
}
