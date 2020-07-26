package main

import "testing"

func TestContentStore_Resolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.content.Resolve(ct.presets)
		expectNoError(t, err)

		c1, _ := ct.GetContent("c1")
		expectEqual(t, c1.X1(), 1.0)
		expectNotSet(t, c1.X2())

		c2, _ := ct.GetContent("c2")
		expectEqual(t, c2.X1(), 2.0)
		expectEqual(t, c2.X2(), 1.0)

		c3, _ := ct.GetContent("c3")
		expectEqual(t, c3.X2(), 23.0)
		expectEqual(t, c3.Y1(), 2.0)
		expectEqual(t, c3.Y2(), 3.0)
		expectEqual(t, c3.XPivot(), 4.0)

		c4, _ := ct.GetContent("c4")
		expectEqual(t, c4.X1(), 1.0)
		expectEqual(t, c4.X2(), 1.0)
		expectEqual(t, c4.Y1(), 2.0)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existing dependency", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveNonExisting.yml")
			err := ct.content.Resolve(ct.presets)
			expectError(t, err)
		})

		t.Run("contradicting values", func(t *testing.T) {
			ct := getCTfromYamlFile(t, "resolveContradicting.yml")
			err := ct.content.Resolve(ct.presets)
			expectError(t, err)
		})
	})
}
