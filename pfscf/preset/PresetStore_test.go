package preset

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func getValidTestStore() (ps Store) {
	ps = NewStore()

	ps.Add(Entry{id: "p1", presets: []string{}, X1: 1.0, X2: 1.0, XPivot: 1.0}) // empty list
	ps.Add(Entry{id: "p2", presets: []string{"p1"}, X2: 23.0, Y1: 2.0})         // list with one element
	ps.Add(Entry{id: "p3", presets: []string{"p1"}, X2: 23.0, Y2: 3.0})         // not contradicting to p2
	ps.Add(Entry{id: "p4", presets: []string{"p2", "p3"}, XPivot: 4.0})         // list with multiple elements
	ps.Add(Entry{id: "p5", presets: []string{"p2", "p2", "p2"}})                // list with duplicate elements

	return ps
}

func TestPresetStore_PresetsAreNotContradicting(t *testing.T) {

	t.Run("errors", func(t *testing.T) {
		testPS := NewStore()

		// contradicting
		testPS.Add(Entry{id: "contradicting1", X1: 1.0})
		testPS.Add(Entry{id: "contradicting2", X1: 2.0})

		t.Run("contradicting entries", func(t *testing.T) {
			err := testPS.PresetsAreNotContradicting("contradicting1", "contradicting2")
			test.ExpectError(t, err, "Contradicting data", "X1")
		})

		t.Run("non-existent presets", func(t *testing.T) {
			for _, presets := range [][]string{
				{"nonExistant"},
				{"nonExistant", "contradicting1"},
				{"contradicting1", "nonExistant"},
			} {
				err := testPS.PresetsAreNotContradicting(presets...)
				test.ExpectError(t, err, "nonExistant", "does not exist")
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		testPS := getValidTestStore()

		for _, presetID := range []string{
			"p1", // empty list
			"p2", // list with one element
			"p4", // list with multiple elements
			"p5", // list with duplicate elements
		} {
			t.Logf("Checking preset '%v'", presetID)
			pe, exists := testPS.Get(presetID)
			test.ExpectTrue(t, exists)
			err := testPS.PresetsAreNotContradicting(pe.presets...)
			test.ExpectNoError(t, err)
		}
	})
}

func TestPresetStore_GetIDs(t *testing.T) {
	testPS := NewStore()
	testPS.Add(Entry{id: "p1"})
	testPS.Add(Entry{id: "p2"})
	testPS.Add(Entry{id: "p0"})

	idList := testPS.GetIDs()

	test.ExpectEqual(t, len(idList), 3)

	// check that all elements returned by list actually exist
	for _, entry := range idList {
		_, exists := testPS.Get(entry)
		test.ExpectTrue(t, exists)
	}

	// check that elements are in expected order (as the result should be sorted)
	test.ExpectEqual(t, idList[0], "p0")
	test.ExpectEqual(t, idList[1], "p1")
	test.ExpectEqual(t, idList[2], "p2")
}

func TestPresetStore_Resolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		testPS := getValidTestStore()

		err := testPS.Resolve()
		test.ExpectNoError(t, err)

		p2, _ := testPS.Get("p2")
		test.ExpectEqual(t, p2.X1, 1.0)
		test.ExpectEqual(t, p2.X2, 23.0)
		test.ExpectEqual(t, p2.Y1, 2.0)
		test.ExpectNotSet(t, p2.Y2)
		test.ExpectEqual(t, p2.XPivot, 1.0)

		p3, _ := testPS.Get("p3")
		test.ExpectEqual(t, p3.X1, 1.0)
		test.ExpectEqual(t, p3.X2, 23.0)
		test.ExpectNotSet(t, p3.Y1)
		test.ExpectEqual(t, p3.Y2, 3.0)
		test.ExpectEqual(t, p3.XPivot, 1.0)

		p4, _ := testPS.Get("p4")
		test.ExpectEqual(t, p4.X1, 1.0)
		test.ExpectEqual(t, p4.X2, 23.0)
		test.ExpectEqual(t, p4.Y1, 2.0)
		test.ExpectEqual(t, p4.Y2, 3.0)
		test.ExpectEqual(t, p4.XPivot, 4.0)

		p5, _ := testPS.Get("p5")
		test.ExpectEqual(t, p5.X1, 1.0)
		test.ExpectEqual(t, p5.X2, 23.0)
		test.ExpectEqual(t, p5.Y1, 2.0)
		test.ExpectNotSet(t, p5.Y2)
		test.ExpectEqual(t, p5.XPivot, 1.0)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("cyclic presets", func(t *testing.T) {
			testPS := NewStore()
			testPS.Add(Entry{id: "p1", presets: []string{"p2"}})
			testPS.Add(Entry{id: "p2", presets: []string{"p3"}})
			testPS.Add(Entry{id: "p3", presets: []string{"p1"}})
			testPS.Add(Entry{id: "p4", presets: []string{"p1"}})

			err := testPS.Resolve()
			test.ExpectError(t, err, "Cyclic dependency", "p1", "p2", "p3")
		})

		t.Run("self-dependency", func(t *testing.T) {
			testPS := NewStore()
			testPS.Add(Entry{id: "p1", presets: []string{"p1"}})

			err := testPS.Resolve()
			test.ExpectError(t, err, "Cyclic dependency", "p1")
		})

		t.Run("non-existing dependency", func(t *testing.T) {
			testPS := NewStore()
			testPS.Add(Entry{id: "p1", presets: []string{"p2"}})

			err := testPS.Resolve()
			test.ExpectError(t, err, "preset", "p2", "cannot be found")
		})

		t.Run("contradicting values", func(t *testing.T) {
			testPS := NewStore()
			testPS.Add(Entry{id: "p1", X1: 1.0})
			testPS.Add(Entry{id: "p2", X1: 2.0})
			testPS.Add(Entry{id: "p3", presets: []string{"p1", "p2"}})

			err := testPS.Resolve()
			test.ExpectError(t, err, "Error resolving preset", "p3", "Contradicting data for field", "X1", "p1", "p2")
		})
	})
}
