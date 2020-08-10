package main

import (
	"fmt"
	"testing"

	"github.com/Blesmol/pfscf/pfscf/stamp"
	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

func init() {
	utils.SetIsTestEnvironment(true)
}

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

func getTestArgStore(key, value string) (as *ArgStore) {
	as, _ = NewArgStore(ArgStoreInit{})
	as.Set(key, value)
	return as
}

func getTestPresetStore(t *testing.T) (ps PresetStore) {
	ps = NewPresetStore(0)
	var (
		data yaml.ContentData
		pe   PresetEntry
	)

	// Add two new presets with same data
	data = getContentDataWithDummyData(t, "unusedType")
	pe = NewPresetEntry("sameData1", data)
	ps.Set(pe.id, pe)
	pe = NewPresetEntry("sameData2", data)
	ps.Set(pe.id, pe)

	// add two conflicting presets
	data = getContentDataWithDummyData(t, "unusedType")
	data.X1 = 10.0
	pe = NewPresetEntry("conflict1", data)
	ps.Set(pe.id, pe)
	data.X1 = 11.0
	pe = NewPresetEntry("conflict2", data)
	ps.Set(pe.id, pe)

	return ps
}

func TestNewContentEntry(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("empty type", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "")
			_, err := NewContentEntry("x", data)

			test.ExpectError(t, err, "No content type provided")
		})

		t.Run("unknown type", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "foo")
			_, err := NewContentEntry("x", data)

			test.ExpectError(t, err, "Unknown content type")
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("TextCell", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			ce, err := NewContentEntry("textCellTest", data)

			test.ExpectNoError(t, err)
			test.ExpectEqual(t, ce.Type(), "textCell")
			test.ExpectEqual(t, ce.ID(), "textCellTest")
		})

		t.Run("SocietyID", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			ce, err := NewContentEntry("societyIdTest", data)

			test.ExpectNoError(t, err)
			test.ExpectEqual(t, ce.Type(), "societyId")
			test.ExpectEqual(t, ce.ID(), "societyIdTest")
		})
	})
}

func TestNewContentTextCell(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		// TODO fill as soon as the ctor returns errors
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell("foo", data)

		test.ExpectNoError(t, err)
		test.ExpectEqual(t, tc.id, "foo")
		test.ExpectEqual(t, tc.description, data.Desc)
		test.ExpectEqual(t, tc.exampleValue, data.Example)
		test.ExpectEqual(t, len(tc.presets), len(data.Presets))
		test.ExpectEqual(t, tc.X1, data.X1)
		test.ExpectEqual(t, tc.Y1, data.Y1)
		test.ExpectEqual(t, tc.X2, data.X2)
		test.ExpectEqual(t, tc.Y2, data.Y2)
		test.ExpectEqual(t, tc.Font, data.Font)
		test.ExpectEqual(t, tc.Fontsize, data.Fontsize)
		test.ExpectEqual(t, tc.Align, data.Align)
	})
}

func TestContentTextCell_BasicGetters(t *testing.T) {
	data := getContentDataWithDummyData(t, "textCell")
	tc, err := NewContentTextCell("foo", data)
	test.ExpectNoError(t, err)

	t.Run("ID", func(t *testing.T) {
		test.ExpectEqual(t, tc.ID(), "foo")
	})

	t.Run("Type", func(t *testing.T) {
		test.ExpectEqual(t, tc.Type(), "textCell")
	})

	t.Run("ExampleValue", func(t *testing.T) {
		test.ExpectEqual(t, tc.ExampleValue(), "Some Example")
	})

	t.Run("UsageExample", func(t *testing.T) {
		test.ExpectEqual(t, tc.UsageExample(), "foo=\"Some Example\"")
	})
}

func TestContentTextCell_IsValid(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")

		t.Run("missing value", func(t *testing.T) {
			tc, _ := NewContentTextCell("foo", data)
			tc.Font = "" // "Unset" one required value

			err := tc.IsValid()
			test.ExpectError(t, err, "Missing value", "Font")
		})

		t.Run("value out of range", func(t *testing.T) {
			tc, _ := NewContentTextCell("foo", data)
			tc.Y2 = 101.0

			err := tc.IsValid()
			test.ExpectError(t, err, "out of range", "Y2")
		})

		t.Run("equal x axis values", func(t *testing.T) {
			tc, _ := NewContentTextCell("foo", data)
			tc.X2 = tc.X1

			err := tc.IsValid()
			test.ExpectError(t, err, "Coordinates for X axis are equal")
		})

		t.Run("equal y axis values", func(t *testing.T) {
			tc, _ := NewContentTextCell("foo", data)
			tc.Y2 = tc.Y1

			err := tc.IsValid()
			test.ExpectError(t, err, "Coordinates for Y axis are equal")
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.X1 = 0.0 // set something to "zero", which is also acceptable
		tc, err := NewContentTextCell("foo", data)
		test.ExpectNoError(t, err)

		err = tc.IsValid()
		test.ExpectNoError(t, err)
	})
}

func TestContentTextCell_Describe(t *testing.T) {
	t.Run("with description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell("someId", data)
		test.ExpectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := tc.Describe(false)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "Some Description")
			test.ExpectStringContainsNot(t, desc, "textCell")
			test.ExpectStringContainsNot(t, desc, "Some Example")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := tc.Describe(true)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "Some Description")
			test.ExpectStringContains(t, desc, "textCell")
			test.ExpectStringContains(t, desc, "Some Example")
		})
	})

	t.Run("without description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.Desc = ""
		data.Example = ""
		tc, err := NewContentTextCell("someId", data)
		test.ExpectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := tc.Describe(false)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "No description available")
			test.ExpectStringContainsNot(t, desc, "textCell")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := tc.Describe(true)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "No description available")
			test.ExpectStringContains(t, desc, "textCell")
		})
	})
}

func TestContentTextCell_Resolve(t *testing.T) {
	ps := getTestPresetStore(t)

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant preset", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			data.Presets = []string{"foo"}
			tc, err := NewContentTextCell("someId", data)
			test.ExpectNoError(t, err)

			_, err = tc.Resolve(ps)
			test.ExpectError(t, err, "does not exist")
		})

		t.Run("conflicting presets", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			data.Presets = []string{"conflict1", "conflict2"}
			tc, err := NewContentTextCell("someId", data)
			test.ExpectNoError(t, err)

			_, err = tc.Resolve(ps)
			test.ExpectError(t, err, "Contradicting data", "X1", "conflict1", "conflict2")
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.Presets = []string{"sameData1", "sameData2"}
		data.Font = "" // set an empty value to be set by presets
		tc, err := NewContentTextCell("someId", data)
		test.ExpectNoError(t, err)

		ceResolved, err := tc.Resolve(ps)
		test.ExpectNoError(t, err)

		tcResolved, castWorked := ceResolved.(ContentTextCell)
		test.ExpectTrue(t, castWorked)

		test.ExpectIsSet(t, tcResolved.Font)
	})
}

func TestContentTextCell_GenerateOutput(t *testing.T) {
	stamp := stamp.NewStamp(100.0, 100.0)
	testID := "someId"
	as := getTestArgStore(testID, "foobar")

	t.Run("errors", func(t *testing.T) {
		t.Run("invalid content object", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			data.Fontsize = 0.0 // unset value, making this textCell invalid
			tc, err := NewContentTextCell(testID, data)
			test.ExpectNoError(t, err)

			err = tc.GenerateOutput(stamp, as)
			test.ExpectError(t, err, "Missing value", "Fontsize")
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell(testID, data)
		test.ExpectNoError(t, err)

		err = tc.GenerateOutput(stamp, as)
		test.ExpectNoError(t, err)
	})
}

// ---------------------------------------------------------------------------------

func TestNewContentSocietyID(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		// TODO fill as soon as the ctor returns errors
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		si, err := NewContentSocietyID("foo", data)

		test.ExpectNoError(t, err)
		test.ExpectEqual(t, si.id, "foo")
		test.ExpectEqual(t, si.description, data.Desc)
		test.ExpectEqual(t, si.exampleValue, data.Example)
		test.ExpectEqual(t, len(si.presets), len(data.Presets))
		test.ExpectEqual(t, si.X1, data.X1)
		test.ExpectEqual(t, si.Y1, data.Y1)
		test.ExpectEqual(t, si.X2, data.X2)
		test.ExpectEqual(t, si.Y2, data.Y2)
		test.ExpectEqual(t, si.XPivot, data.XPivot)
		test.ExpectEqual(t, si.Font, data.Font)
		test.ExpectEqual(t, si.Fontsize, data.Fontsize)
	})
}

func TestContentSocietyID_BasicGetters(t *testing.T) {
	data := getContentDataWithDummyData(t, "societyId")
	si, err := NewContentSocietyID("foo", data)
	test.ExpectNoError(t, err)

	t.Run("ID", func(t *testing.T) {
		test.ExpectEqual(t, si.ID(), "foo")
	})

	t.Run("Type", func(t *testing.T) {
		test.ExpectEqual(t, si.Type(), "societyId")
	})

	t.Run("ExampleValue", func(t *testing.T) {
		test.ExpectEqual(t, si.ExampleValue(), "Some Example")
	})

	t.Run("UsageExample", func(t *testing.T) {
		test.ExpectEqual(t, si.UsageExample(), "foo=\"Some Example\"")
	})
}

func TestContentSocietyID_IsValid(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyID")

		t.Run("missing value", func(t *testing.T) {
			si, err := NewContentSocietyID("foo", data)
			si.Font = "" // "Unset" one required value

			err = si.IsValid()
			test.ExpectError(t, err, "Missing value")
		})

		t.Run("xpivot range violation", func(t *testing.T) {
			for _, testPivot := range []float64{5.0, 10.0, 20.0, 30.0} {
				t.Logf("Testing pivot=%v", testPivot)
				data.XPivot = testPivot

				si, err := NewContentSocietyID("foo", data)
				si.X1 = 10.0
				si.X2 = 20.0

				err = si.IsValid()
				test.ExpectError(t, err, "xpivot value must lie between")
			}
		})

		t.Run("value out of range", func(t *testing.T) {
			tc, _ := NewContentSocietyID("foo", data)
			tc.Y2 = 101.0

			err := tc.IsValid()
			test.ExpectError(t, err, "out of range", "Y2")
		})

		t.Run("equal x axis values", func(t *testing.T) {
			tc, _ := NewContentSocietyID("foo", data)
			tc.X2 = tc.X1
			tc.XPivot = tc.X1

			err := tc.IsValid()
			test.ExpectError(t, err, "Coordinates for X axis are equal")
		})

		t.Run("equal y axis values", func(t *testing.T) {
			tc, _ := NewContentSocietyID("foo", data)
			tc.Y2 = tc.Y1

			err := tc.IsValid()
			test.ExpectError(t, err, "Coordinates for Y axis are equal")
		})

	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		si, err := NewContentSocietyID("foo", data)
		test.ExpectNoError(t, err)
		si.X1 = 0.0 // also acceptable now

		err = si.IsValid()
		test.ExpectNoError(t, err)
	})
}

func TestContentSocietyID_Describe(t *testing.T) {
	t.Run("with description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		si, err := NewContentSocietyID("someId", data)
		test.ExpectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := si.Describe(false)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "Some Description")
			test.ExpectStringContainsNot(t, desc, "societyId")
			test.ExpectStringContainsNot(t, desc, "Some Example")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := si.Describe(true)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "Some Description")
			test.ExpectStringContains(t, desc, "societyId")
			test.ExpectStringContains(t, desc, "Some Example")
		})
	})

	t.Run("without description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		data.Desc = ""
		data.Example = ""
		si, err := NewContentSocietyID("someId", data)
		test.ExpectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := si.Describe(false)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "No description available")
			test.ExpectStringContainsNot(t, desc, "societyId")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := si.Describe(true)
			test.ExpectStringContains(t, desc, "someId")
			test.ExpectStringContains(t, desc, "No description available")
			test.ExpectStringContains(t, desc, "societyId")
		})
	})
}

func TestContentSocietyID_Resolve(t *testing.T) {
	ps := getTestPresetStore(t)

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant preset", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			data.Presets = []string{"foo"}
			si, err := NewContentSocietyID("someId", data)
			test.ExpectNoError(t, err)

			_, err = si.Resolve(ps)
			test.ExpectError(t, err, "does not exist")
		})

		t.Run("conflicting presets", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			data.Presets = []string{"conflict1", "conflict2"}
			si, err := NewContentSocietyID("someId", data)
			test.ExpectNoError(t, err)

			_, err = si.Resolve(ps)
			test.ExpectError(t, err, "Contradicting", "X1", "conflict1", "conflict2")
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		data.Presets = []string{"sameData1", "sameData2"}
		data.Font = "" // set an empty value to be set by presets
		si, err := NewContentSocietyID("someId", data)
		test.ExpectNoError(t, err)

		ceResolved, err := si.Resolve(ps)
		test.ExpectNoError(t, err)

		siResolved, castWorked := ceResolved.(ContentSocietyID)
		test.ExpectTrue(t, castWorked)

		test.ExpectIsSet(t, siResolved.Font)
	})
}

func TestContentSocietyID_GenerateOutput(t *testing.T) {
	stamp := stamp.NewStamp(100.0, 100.0)
	testID := "someId"
	as := getTestArgStore(testID, "12345-678")

	t.Run("errors", func(t *testing.T) {
		t.Run("invalid content object", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			data.Font = "" // unset value, making this textCell invalid
			si, err := NewContentSocietyID(testID, data)
			test.ExpectNoError(t, err)

			err = si.GenerateOutput(stamp, as)
			test.ExpectError(t, err, "Missing value", "Font")
		})

		t.Run("value with invalid format", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			si, err := NewContentSocietyID(testID, data)
			test.ExpectNoError(t, err)

			for _, invalidSocietyID := range []string{"", "foo", "a123-456", "123-456b", "1"} {
				asInvalid := getTestArgStore(testID, invalidSocietyID)
				err = si.GenerateOutput(stamp, asInvalid)
				test.ExpectError(t, err, "does not follow the pattern")
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "SocietyId")
		si, err := NewContentSocietyID(testID, data)
		test.ExpectNoError(t, err)

		for _, societyID := range []string{"-", "1-", "-2", "123-456"} {
			asValid := getTestArgStore(testID, societyID)
			err = si.GenerateOutput(stamp, asValid)
			test.ExpectNoError(t, err)
		}
	})
}

// ---------------------------------------------------------------------------------

func TestCheckThatAllExportedFieldsAreSet(t *testing.T) {
	type testStruct struct {
		A, b, C, d float64
	}

	t.Run("errors", func(t *testing.T) {
		testVal := testStruct{A: 1.0, b: 2.0}
		err := CheckThatAllExportedFieldsAreSet(testVal)
		test.ExpectError(t, err, "Missing value", "C")
	})

	t.Run("valid", func(t *testing.T) {
		testVal := testStruct{A: 1.0, C: 3.0, d: 4.0}
		err := CheckThatAllExportedFieldsAreSet(testVal)
		test.ExpectNoError(t, err)
	})
}

func TestCheckFieldsAreSet(t *testing.T) {
	type testStruct struct {
		A, b, C, d float64
		E, F       string
	}

	t.Run("errors", func(t *testing.T) {
		testVal := testStruct{A: 1.0, b: 2.0, E: "foo"}
		for _, testFields := range [][]string{
			{"A", "C"},
			{"C"},
			{"A", "E", "F"},
			{"F"},
		} {
			err := checkFieldsAreSet(testVal, testFields...)
			test.ExpectError(t, err, "Missing value")
		}
	})

	t.Run("valid", func(t *testing.T) {
		testVal := testStruct{A: 1.0, C: 3.0, E: "foo"}
		for _, testFields := range [][]string{
			{"A"},
			{"C"},
			{"A", "C"},
			{"E"},
			{"A", "C", "E"},
			{},
		} {
			err := checkFieldsAreSet(testVal, testFields...)
			test.ExpectNoError(t, err)
		}
	})
}

func TestCheckFieldsAreInRange(t *testing.T) {
	type testStruct struct {
		A, B, C float64
	}

	t.Run("errors", func(t *testing.T) {
		testVal := testStruct{A: -1.0, B: 101.0}
		for _, testFields := range [][]string{
			{"A"},
			{"B"},
			{"A", "B"},
		} {
			err := checkFieldsAreInRange(testVal, testFields...)
			test.ExpectError(t, err, "out of range")
		}
	})

	t.Run("valid", func(t *testing.T) {
		testVal := testStruct{A: 0.0, B: 50.0, C: 100.0}
		for _, testFields := range [][]string{
			{"A"},
			{"B"},
			{"C"},
			{"A", "B", "C"},
			{},
		} {
			err := checkFieldsAreInRange(testVal, testFields...)
			test.ExpectNoError(t, err)
		}
	})
}

func TestContentValErr(t *testing.T) {
	data := getContentDataWithDummyData(t, "textCell")
	tc, err := NewContentTextCell("testId", data)
	test.ExpectNoError(t, err)

	errIn := fmt.Errorf("Test text")
	errOut := contentValErr(tc, errIn)
	test.ExpectError(t, errOut, "Error validating content", "testId", "Test text")
}
