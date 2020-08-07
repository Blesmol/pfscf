package main

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func TestPresetStore_PresetsAreNotContradicting(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("contradicting entries", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.presets.PresetsAreNotContradicting("p1", "p2")
			test.ExpectError(t, err)
		})

		t.Run("non-existent presets", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml") // only contains p1
			err := ct.presets.PresetsAreNotContradicting("p1", "p2")
			test.ExpectError(t, err)

			err = ct.presets.PresetsAreNotContradicting("p2")
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		t.Run("empty list", func(t *testing.T) {
			p1, exists := ct.presets.Get("p1")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, len(p1.presets), 0) // empty list
			err := ct.presets.PresetsAreNotContradicting(p1.presets...)
			test.ExpectNoError(t, err)
		})

		t.Run("list with one element", func(t *testing.T) {
			p2, exists := ct.presets.Get("p2")
			test.ExpectTrue(t, exists)
			test.ExpectEqual(t, len(p2.presets), 1)
			err := ct.presets.PresetsAreNotContradicting(p2.presets...)
			test.ExpectNoError(t, err)
		})

		t.Run("list with more elements", func(t *testing.T) {
			p4, exists := ct.presets.Get("p4")
			test.ExpectTrue(t, exists)
			test.ExpectTrue(t, len(p4.presets) > 1)
			err := ct.presets.PresetsAreNotContradicting(p4.presets...)
			test.ExpectNoError(t, err)
		})

		t.Run("list with duplicate elements", func(t *testing.T) {
			p5, exists := ct.presets.Get("p5")
			test.ExpectTrue(t, exists)
			test.ExpectTrue(t, len(p5.presets) > 1)
			err := ct.presets.PresetsAreNotContradicting(p5.presets...)
			test.ExpectNoError(t, err)
		})
	})
}

func TestPresetStore_GetIDs(t *testing.T) {
	ct := getCTfromYamlFile(t, "valid.yml")

	idList := ct.presets.GetIDs()

	test.ExpectEqual(t, len(idList), 3)

	// check that all elements returned by list actually exist
	for _, entry := range idList {
		_, exists := ct.presets.Get(entry)
		test.ExpectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	test.ExpectEqual(t, idList[0], "p0")
	test.ExpectEqual(t, idList[1], "p1")
	test.ExpectEqual(t, idList[2], "p2")
}

func TestPresetStore_Resolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.presets.Resolve()
		test.ExpectNoError(t, err)

		p2, _ := ct.presets.Get("p2")
		test.ExpectEqual(t, p2.X1, 1.0)
		test.ExpectEqual(t, p2.X2, 23.0)
		test.ExpectEqual(t, p2.Y1, 2.0)
		test.ExpectNotSet(t, p2.Y2)
		test.ExpectEqual(t, p2.XPivot, 1.0)

		p3, _ := ct.presets.Get("p3")
		test.ExpectEqual(t, p3.X1, 1.0)
		test.ExpectEqual(t, p3.X2, 23.0)
		test.ExpectNotSet(t, p3.Y1)
		test.ExpectEqual(t, p3.Y2, 3.0)
		test.ExpectEqual(t, p3.XPivot, 1.0)

		p4, _ := ct.presets.Get("p4")
		test.ExpectEqual(t, p4.X1, 1.0)
		test.ExpectEqual(t, p4.X2, 23.0)
		test.ExpectEqual(t, p4.Y1, 2.0)
		test.ExpectEqual(t, p4.Y2, 3.0)
		test.ExpectEqual(t, p4.XPivot, 4.0)

		p5, _ := ct.presets.Get("p5")
		test.ExpectEqual(t, p5.X1, 1.0)
		test.ExpectEqual(t, p5.X2, 23.0)
		test.ExpectEqual(t, p5.Y1, 2.0)
		test.ExpectNotSet(t, p5.Y2)
		test.ExpectEqual(t, p5.XPivot, 1.0)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("cyclic presets", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveCyclic.yml")
			err := ct.presets.Resolve()
			test.ExpectError(t, err)
		})

		t.Run("self-dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveSelf.yml")
			err := ct.presets.Resolve()
			test.ExpectError(t, err)
		})

		t.Run("non-existing dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml")
			err := ct.presets.Resolve()
			test.ExpectError(t, err)
		})

		t.Run("contradicting values", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.presets.Resolve()
			test.ExpectError(t, err)
		})
	})
}
