package content

import (
	"testing"

	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func getTextCellWithDummyData(presets ...string) (tc *textCell) {
	tc = newTextCell()

	tc.Value = "Some value"
	tc.X = 12.0
	tc.Y = 12.0
	tc.X2 = 24.0
	tc.Y2 = 24.0
	tc.Font = "Helvetica"
	tc.Fontsize = 14.0
	tc.Align = "CB"
	for _, preset := range presets {
		tc.Presets = append(tc.Presets, preset)
	}

	return tc
}

func TestTextCell_IsValid(t *testing.T) {
	paramStore := param.NewStore()
	t.Run("errors", func(t *testing.T) {
		t.Run("missing value", func(t *testing.T) {
			tc := getTextCellWithDummyData()
			tc.Font = "" // "Unset" one required value

			err := tc.isValid(&paramStore)
			test.ExpectError(t, err, "Missing value", "Font")
		})

		t.Run("value out of range", func(t *testing.T) {
			tc := getTextCellWithDummyData()
			tc.Y2 = 101.0

			err := tc.isValid(&paramStore)
			test.ExpectError(t, err, "out of range", "Y2")
		})

		t.Run("equal x axis values", func(t *testing.T) {
			tc := getTextCellWithDummyData()
			tc.X2 = tc.X

			err := tc.isValid(&paramStore)
			test.ExpectError(t, err, "Coordinates for X axis are equal")
		})

		t.Run("equal y axis values", func(t *testing.T) {
			tc := getTextCellWithDummyData()
			tc.Y2 = tc.Y

			err := tc.isValid(&paramStore)
			test.ExpectError(t, err, "Coordinates for Y axis are equal")
		})
	})

	t.Run("valid", func(t *testing.T) {
		tc := getTextCellWithDummyData()
		tc.X = 0.0 // set something to "zero", which is also acceptable

		err := tc.isValid(&paramStore)
		test.ExpectNoError(t, err)
	})
}

func TestTextCell_Resolve(t *testing.T) {
	ps := getTestPresetStore(t)

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant preset", func(t *testing.T) {
			tc := getTextCellWithDummyData("non-existing preset")

			err := tc.resolve(*ps)
			test.ExpectError(t, err, "does not exist")
		})

		t.Run("conflicting presets", func(t *testing.T) {
			tc := getTextCellWithDummyData("conflict1", "conflict2")

			err := tc.resolve(*ps)
			test.ExpectError(t, err, "Contradicting values", "font", "conflict1", "conflict2")
		})
	})

	t.Run("valid", func(t *testing.T) {
		tc := getTextCellWithDummyData("sameData1", "sameData2")
		tc.Font = ""

		err := tc.resolve(*ps)
		test.ExpectNoError(t, err)

		test.ExpectIsSet(t, tc.Font)
	})
}

func TestTextCell_generateOutput(t *testing.T) {
	stamp := stamp.NewStamp(100.0, 100.0)
	testArgName := "someId"
	testArgValue := "foobar"
	as := getTestArgStore(testArgName, testArgValue)

	t.Run("valid", func(t *testing.T) {
		tc := getTextCellWithDummyData()
		tc.Value = "param:someId"

		err := tc.generateOutput(stamp, as)
		test.ExpectNoError(t, err)
	})
}

func TestTextCell_deepCopy(t *testing.T) {
	e1 := newTextCell()
	e1.Value = "t1"
	e1.Presets = append(e1.Presets, "t1")

	e2 := e1.deepCopy().(*textCell)
	e2.Value = "t2"
	e2.Presets[0] = "t2"

	test.ExpectNotEqual(t, e1.Value, e2.Value)
	test.ExpectNotEqual(t, e1.Presets[0], e2.Presets[0])
}
