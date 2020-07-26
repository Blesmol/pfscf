package main

import "testing"

func init() {
	SetIsTestEnvironment(true)
}

func getContentDataWithDummyData(t *testing.T, cdType string) (cd ContentData) {
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
	cd.Example = "Some Example"
	cd.Presets = []string{"Some Preset"}

	expectAllExportedSet(t, cd) // to be sure that we also get all new fields

	// overvwrite type after the "expect..." check as cdType could be intentionally empty
	cd.Type = cdType

	return cd
}

func getTestPresetStore(t *testing.T) (ps PresetStore) {
	ps = NewPresetStore(0)
	var (
		data ContentData
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

			expectError(t, err)
		})

		t.Run("unknown type", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "foo")
			_, err := NewContentEntry("x", data)

			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("TextCell", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			ce, err := NewContentEntry("textCellTest", data)

			expectNoError(t, err)
			expectEqual(t, ce.Type(), "textCell")
			expectEqual(t, ce.ID(), "textCellTest")
		})

		t.Run("SocietyID", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			ce, err := NewContentEntry("societyIdTest", data)

			expectNoError(t, err)
			expectEqual(t, ce.Type(), "societyId")
			expectEqual(t, ce.ID(), "societyIdTest")
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

		expectNoError(t, err)
		expectEqual(t, tc.id, "foo")
		expectEqual(t, tc.description, data.Desc)
		expectEqual(t, tc.exampleValue, data.Example)
		expectEqual(t, len(tc.presets), len(data.Presets))
		expectEqual(t, tc.X1, data.X1)
		expectEqual(t, tc.Y1, data.Y1)
		expectEqual(t, tc.X2, data.X2)
		expectEqual(t, tc.Y2, data.Y2)
		expectEqual(t, tc.Font, data.Font)
		expectEqual(t, tc.Fontsize, data.Fontsize)
		expectEqual(t, tc.Align, data.Align)
	})
}

func TestContentTextCell_BasicGetters(t *testing.T) {
	data := getContentDataWithDummyData(t, "textCell")
	tc, err := NewContentTextCell("foo", data)
	expectNoError(t, err)

	t.Run("ID", func(t *testing.T) {
		expectEqual(t, tc.ID(), "foo")
	})

	t.Run("Type", func(t *testing.T) {
		expectEqual(t, tc.Type(), "textCell")
	})

	t.Run("ExampleValue", func(t *testing.T) {
		expectEqual(t, tc.ExampleValue(), "Some Example")
	})

	t.Run("UsageExample", func(t *testing.T) {
		expectEqual(t, tc.UsageExample(), "foo=\"Some Example\"")
	})
}

func TestContentTextCell_IsValid(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.Font = "" // "Unset" one required value
		tc, err := NewContentTextCell("foo", data)
		expectNoError(t, err)

		err = tc.IsValid()
		expectError(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell("foo", data)
		expectNoError(t, err)

		err = tc.IsValid()
		expectNoError(t, err)
	})
}

func TestContentTextCell_Describe(t *testing.T) {
	t.Run("with description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell("someId", data)
		expectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := tc.Describe(false)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "Some Description")
			expectStringContainsNot(t, desc, "textCell")
			expectStringContainsNot(t, desc, "Some Example")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := tc.Describe(true)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "Some Description")
			expectStringContains(t, desc, "textCell")
			expectStringContains(t, desc, "Some Example")
		})
	})

	t.Run("without description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.Desc = ""
		data.Example = ""
		tc, err := NewContentTextCell("someId", data)
		expectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := tc.Describe(false)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "No description available")
			expectStringContainsNot(t, desc, "textCell")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := tc.Describe(true)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "No description available")
			expectStringContains(t, desc, "textCell")
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
			expectNoError(t, err)

			_, err = tc.Resolve(ps)
			expectError(t, err)
		})

		t.Run("conflicting presets", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			data.Presets = []string{"conflict1", "conflict2"}
			tc, err := NewContentTextCell("someId", data)
			expectNoError(t, err)

			_, err = tc.Resolve(ps)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		data.Presets = []string{"sameData1", "sameData2"}
		data.Font = "" // set an empty value to be set by presets
		tc, err := NewContentTextCell("someId", data)
		expectNoError(t, err)

		ceResolved, err := tc.Resolve(ps)
		expectNoError(t, err)

		tcResolved, castWorked := ceResolved.(ContentTextCell)
		expectTrue(t, castWorked)

		expectIsSet(t, tcResolved.Font)
	})
}

func TestContentTextCell_GenerateOutput(t *testing.T) {
	stamp := NewStamp(100.0, 100.0)
	value := "foobar"

	t.Run("errors", func(t *testing.T) {
		t.Run("invalid content object", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			data.X1 = 0 // unset value, making this textCell invalid
			tc, err := NewContentTextCell("someId", data)
			expectNoError(t, err)

			err = tc.GenerateOutput(stamp, &value)
			expectError(t, err)
		})

		t.Run("missing value", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "textCell")
			tc, err := NewContentTextCell("someId", data)
			expectNoError(t, err)

			err = tc.GenerateOutput(stamp, nil)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "textCell")
		tc, err := NewContentTextCell("someId", data)
		expectNoError(t, err)

		err = tc.GenerateOutput(stamp, &value)
		expectNoError(t, err)
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

		expectNoError(t, err)
		expectEqual(t, si.id, "foo")
		expectEqual(t, si.description, data.Desc)
		expectEqual(t, si.exampleValue, data.Example)
		expectEqual(t, len(si.presets), len(data.Presets))
		expectEqual(t, si.X1, data.X1)
		expectEqual(t, si.Y1, data.Y1)
		expectEqual(t, si.X2, data.X2)
		expectEqual(t, si.Y2, data.Y2)
		expectEqual(t, si.XPivot, data.XPivot)
		expectEqual(t, si.Font, data.Font)
		expectEqual(t, si.Fontsize, data.Fontsize)
	})
}

func TestContentSocietyID_BasicGetters(t *testing.T) {
	data := getContentDataWithDummyData(t, "societyId")
	si, err := NewContentSocietyID("foo", data)
	expectNoError(t, err)

	t.Run("ID", func(t *testing.T) {
		expectEqual(t, si.ID(), "foo")
	})

	t.Run("Type", func(t *testing.T) {
		expectEqual(t, si.Type(), "societyId")
	})

	t.Run("ExampleValue", func(t *testing.T) {
		expectEqual(t, si.ExampleValue(), "Some Example")
	})

	t.Run("UsageExample", func(t *testing.T) {
		expectEqual(t, si.UsageExample(), "foo=\"Some Example\"")
	})
}

func TestContentSocietyID_IsValid(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("missing value", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyID")
			data.Font = "" // "Unset" one required value
			si, err := NewContentSocietyID("foo", data)
			expectNoError(t, err)

			err = si.IsValid()
			expectError(t, err)
		})

		t.Run("range violation", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyID")
			data.X1 = 10.0
			data.X2 = 20.0
			for _, testPivot := range []float64{-1.0, 10.0, 20.0, 30.0} {
				t.Logf("Testing pivot=%v", testPivot)
				data.XPivot = testPivot

				si, err := NewContentSocietyID("foo", data)
				expectNoError(t, err)

				err = si.IsValid()
				expectError(t, err)
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		si, err := NewContentSocietyID("foo", data)
		expectNoError(t, err)

		err = si.IsValid()
		expectNoError(t, err)
	})
}

func TestContentSocietyID_Describe(t *testing.T) {
	t.Run("with description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		si, err := NewContentSocietyID("someId", data)
		expectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := si.Describe(false)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "Some Description")
			expectStringContainsNot(t, desc, "societyId")
			expectStringContainsNot(t, desc, "Some Example")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := si.Describe(true)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "Some Description")
			expectStringContains(t, desc, "societyId")
			expectStringContains(t, desc, "Some Example")
		})
	})

	t.Run("without description and example", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		data.Desc = ""
		data.Example = ""
		si, err := NewContentSocietyID("someId", data)
		expectNoError(t, err)

		t.Run("non-verbose", func(t *testing.T) {
			desc := si.Describe(false)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "No description available")
			expectStringContainsNot(t, desc, "societyId")
		})

		t.Run("verbose", func(t *testing.T) {
			desc := si.Describe(true)
			expectStringContains(t, desc, "someId")
			expectStringContains(t, desc, "No description available")
			expectStringContains(t, desc, "societyId")
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
			expectNoError(t, err)

			_, err = si.Resolve(ps)
			expectError(t, err)
		})

		t.Run("conflicting presets", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			data.Presets = []string{"conflict1", "conflict2"}
			si, err := NewContentSocietyID("someId", data)
			expectNoError(t, err)

			_, err = si.Resolve(ps)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "societyId")
		data.Presets = []string{"sameData1", "sameData2"}
		data.Font = "" // set an empty value to be set by presets
		si, err := NewContentSocietyID("someId", data)
		expectNoError(t, err)

		ceResolved, err := si.Resolve(ps)
		expectNoError(t, err)

		siResolved, castWorked := ceResolved.(ContentSocietyID)
		expectTrue(t, castWorked)

		expectIsSet(t, siResolved.Font)
	})
}

func TestContentSocietyID_GenerateOutput(t *testing.T) {
	stamp := NewStamp(100.0, 100.0)
	validValue := "12345-678"

	t.Run("errors", func(t *testing.T) {
		t.Run("invalid content object", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			data.X1 = 0 // unset value, making this textCell invalid
			si, err := NewContentSocietyID("someId", data)
			expectNoError(t, err)

			err = si.GenerateOutput(stamp, &validValue)
			expectError(t, err)
		})

		t.Run("missing value", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			si, err := NewContentSocietyID("someId", data)
			expectNoError(t, err)

			err = si.GenerateOutput(stamp, nil)
			expectError(t, err)
		})

		t.Run("value with invalid format", func(t *testing.T) {
			data := getContentDataWithDummyData(t, "societyId")
			si, err := NewContentSocietyID("someId", data)
			expectNoError(t, err)

			for _, invalidSocietyID := range []string{"", "foo", "a123-456", "123-456b", "1"} {
				err = si.GenerateOutput(stamp, &invalidSocietyID)
				expectError(t, err)
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		data := getContentDataWithDummyData(t, "SocietyId")
		si, err := NewContentSocietyID("someId", data)
		expectNoError(t, err)

		for _, societyID := range []string{"-", "1-", "-2", "123-456"} {
			err = si.GenerateOutput(stamp, &societyID)
			expectNoError(t, err)
		}
	})
}

// ---------------------------------------------------------------------------------

func TestAddMissingValues(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("exported and unexported fields", func(t *testing.T) {
			type testStruct struct{ A, b, C, d, E float64 }
			source := testStruct{A: 1.0, b: 2.0, C: 3.0, d: 4.0}
			target := testStruct{A: 10.0, b: 11.0}
			AddMissingValues(&target, source)

			expectEqual(t, target.A, 10.0)
			expectEqual(t, target.b, 11.0)
			expectEqual(t, target.C, 3.0)
			expectNotSet(t, target.d)
			expectNotSet(t, target.E)
		})

		t.Run("supported datatypes", func(t *testing.T) {
			type testStruct struct {
				A, B, C float64
				D, E, F string
			}
			source := testStruct{A: 1.0, B: 2.0 /* C left empty */, D: "4.0", E: "5.0" /* F left empty*/}
			target := testStruct{A: 10.0, D: "14.0"}
			AddMissingValues(&target, source)

			expectEqual(t, target.A, 10.0)
			expectEqual(t, target.B, 2.0)
			expectNotSet(t, target.C)
			expectEqual(t, target.D, "14.0")
			expectEqual(t, target.E, "5.0")
			expectNotSet(t, target.F)
		})

		t.Run("different exported fields", func(t *testing.T) {
			source := struct {
				Common, OnlySource float64
			}{
				Common: 1.0, OnlySource: 2.0,
			}

			target := struct {
				Common, OnlyTarget float64
			}{}

			AddMissingValues(&target, source)

			expectEqual(t, target.Common, 1.0)
			expectNotSet(t, target.OnlyTarget)
		})

		t.Run("ignore fields", func(t *testing.T) {
			type testStruct struct {
				A, B, C, D float64
			}
			source := testStruct{A: 1.0, B: 2.0, C: 3.0, D: 4.0}
			target := testStruct{}

			AddMissingValues(&target, source, "B", "C", "a", "De")

			expectEqual(t, target.A, 1.0)
			expectNotSet(t, target.B)
			expectNotSet(t, target.C)
			expectEqual(t, target.D, 4.0)
		})
	})
}

func TestCheckThatAllExportedFieldsAreSet(t *testing.T) {
	type testStruct struct {
		A, b, C, d float64
	}

	t.Run("errors", func(t *testing.T) {
		testVal := testStruct{A: 1.0, b: 2.0}
		err := CheckThatAllExportedFieldsAreSet(testVal)
		expectError(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		testVal := testStruct{A: 1.0, C: 3.0, d: 4.0}
		err := CheckThatAllExportedFieldsAreSet(testVal)
		expectNoError(t, err)
	})
}
