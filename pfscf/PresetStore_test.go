package main

import "testing"

func TestPresetStore_PresetsAreNotContradicting(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("contradicting entries", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.presets.PresetsAreNotContradicting("p1", "p2")
			expectError(t, err)
		})

		t.Run("non-existent presets", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml") // only contains p1
			err := ct.presets.PresetsAreNotContradicting("p1", "p2")
			expectError(t, err)

			err = ct.presets.PresetsAreNotContradicting("p2")
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		t.Run("empty list", func(t *testing.T) {
			p1, exists := ct.GetPreset("p1")
			expectTrue(t, exists)
			expectEqual(t, len(p1.Presets), 0) // empty list
			err := ct.presets.PresetsAreNotContradicting(p1.Presets...)
			expectNoError(t, err)
		})

		t.Run("list with one element", func(t *testing.T) {
			p2, exists := ct.GetPreset("p2")
			expectTrue(t, exists)
			expectEqual(t, len(p2.Presets), 1)
			err := ct.presets.PresetsAreNotContradicting(p2.Presets...)
			expectNoError(t, err)
		})

		t.Run("list with more elements", func(t *testing.T) {
			p4, exists := ct.GetPreset("p4")
			expectTrue(t, exists)
			expectTrue(t, len(p4.Presets) > 1)
			err := ct.presets.PresetsAreNotContradicting(p4.Presets...)
			expectNoError(t, err)
		})

		t.Run("list with duplicate elements", func(t *testing.T) {
			p5, exists := ct.GetPreset("p5")
			expectTrue(t, exists)
			expectTrue(t, len(p5.Presets) > 1)
			err := ct.presets.PresetsAreNotContradicting(p5.Presets...)
			expectNoError(t, err)
		})
	})
}

func TestPresetStore_GetIDs(t *testing.T) {
	ct := getCTfromYamlFile(t, "valid.yml")

	idList := ct.presets.GetIDs()

	expectEqual(t, len(idList), 3)

	// check that all elements returned by list actually exist
	for _, entry := range idList {
		_, exists := ct.GetPreset(entry)
		expectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	expectEqual(t, idList[0], "p0")
	expectEqual(t, idList[1], "p1")
	expectEqual(t, idList[2], "p2")
}

func TestPresetStore_Resolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.presets.Resolve()
		expectNoError(t, err)

		p2, _ := ct.GetPreset("p2")
		expectEqual(t, p2.X1, 1.0)
		expectEqual(t, p2.X2, 23.0)
		expectEqual(t, p2.Y1, 2.0)
		expectNotSet(t, p2.Y2)
		expectEqual(t, p2.XPivot, 1.0)

		p3, _ := ct.GetPreset("p3")
		expectEqual(t, p3.X1, 1.0)
		expectEqual(t, p3.X2, 23.0)
		expectNotSet(t, p3.Y1)
		expectEqual(t, p3.Y2, 3.0)
		expectEqual(t, p3.XPivot, 1.0)

		p4, _ := ct.GetPreset("p4")
		expectEqual(t, p4.X1, 1.0)
		expectEqual(t, p4.X2, 23.0)
		expectEqual(t, p4.Y1, 2.0)
		expectEqual(t, p4.Y2, 3.0)
		expectEqual(t, p4.XPivot, 4.0)

		p5, _ := ct.GetPreset("p5")
		expectEqual(t, p5.X1, 1.0)
		expectEqual(t, p5.X2, 23.0)
		expectEqual(t, p5.Y1, 2.0)
		expectNotSet(t, p5.Y2)
		expectEqual(t, p5.XPivot, 1.0)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("cyclic presets", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveCyclic.yml")
			err := ct.presets.Resolve()
			expectError(t, err)
		})

		t.Run("self-dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveSelf.yml")
			err := ct.presets.Resolve()
			expectError(t, err)
		})

		t.Run("non-existing dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml")
			err := ct.presets.Resolve()
			expectError(t, err)
		})

		t.Run("contradicting values", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.presets.Resolve()
			expectError(t, err)
		})
	})
}
