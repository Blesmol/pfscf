package main

/*
func TestContentStore_Resolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ct := getCTfromYamlFile(t, "resolveValid.yml")

		err := ct.content.Resolve(ct.presets)
		expectNoError(t, err)

		ce1, _ := ct.GetContent("c1")
		expectEqual(t, ce1.X1(), 1.0)
		expectNotSet(t, ce1.X2())

		ce2, _ := ct.GetContent("c2")
		expectEqual(t, ce2.X1(), 2.0)
		expectEqual(t, ce2.X2(), 1.0)

		ce3, _ := ct.GetContent("c3")
		expectEqual(t, ce3.X2(), 23.0)
		expectEqual(t, ce3.Y1(), 2.0)
		expectEqual(t, ce3.Y2(), 3.0)
		expectEqual(t, ce3.XPivot(), 4.0)

		ce4, _ := ct.GetContent("c4")
		expectEqual(t, ce4.X1(), 1.0)
		expectEqual(t, ce4.X2(), 1.0)
		expectEqual(t, ce4.Y1(), 2.0)
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

*/
