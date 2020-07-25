package main

import "testing"

func init() {
	SetIsTestEnvironment(true)
}

func getContentDataWithDummyData(t *testing.T, cdType string) (cd ContentData) {
	cd.Type = cdType
	cd.Desc = "Some Description"
	cd.X1 = 12.0
	cd.Y1 = 12.0
	cd.X2 = 24.0
	cd.Y2 = 24.0
	cd.XPivot = 36.0
	cd.Font = "Helvetica"
	cd.Fontsize = 14.0
	cd.Align = "LB"
	cd.Example = "Some Example"
	cd.Presets = []string{"Some Preset"}

	expectAllExportedSet(t, cd) // to be sure that we also get all new fields

	return cd
}

func TestNewContentEntry(t *testing.T) {
	cd := getContentDataWithDummyData(t, "myType")
	ce := NewContentEntry("myId", cd)

	expectEqual(t, ce.ID(), "myId")
	expectEqual(t, ce.Type(), "myType")
	expectEqual(t, ce.Description(), "Some Description")
	expectEqual(t, ce.X1(), 12.0)
	expectEqual(t, ce.Y1(), 12.0)
	expectEqual(t, ce.X2(), 24.0)
	expectEqual(t, ce.Y2(), 24.0)
	expectEqual(t, ce.XPivot(), 36.0)
	expectEqual(t, ce.Font(), "Helvetica")
	expectEqual(t, ce.Fontsize(), 14.0)
	expectEqual(t, ce.Align(), "LB")
	expectEqual(t, ce.Example(), "Some Example")
}

func TestCheckThatValuesArePresent(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		t.Run("missing value", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "some type")
			cd.Font = ""
			ce := NewContentEntry("id", cd)
			err := ce.CheckThatValuesArePresent("Font")
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("all values set", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "some type")
			ce := NewContentEntry("id", cd)
			err := ce.CheckThatValuesArePresent("Type", "Desc", "X1", "X2", "Y1", "Y2", "XPivot", "Font", "Fontsize", "Align", "Example")
			expectNoError(t, err)
		})

		t.Run("only check existing values", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "some type")
			cd.X2 = 0.0
			cd.Font = ""
			ce := NewContentEntry("id", cd)
			err := ce.CheckThatValuesArePresent("X1", "Y2", "Desc")
			expectNoError(t, err)
		})
	})
}

func TestContentEntry_IsValid(t *testing.T) {

	t.Run("general", func(t *testing.T) {
		t.Run("missing type", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "willBeRemovedOneLineLater")
			cd.Type = ""
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectError(t, err)
		})

		t.Run("invalid type", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "textCellX")
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectError(t, err)
		})
	})

	t.Run("textCell", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "textCell")
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectNoError(t, err)
		})

		t.Run("missing value", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "textCell")
			cd.Font = ""
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectError(t, err)
		})
	})

	t.Run("societyId", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "societyId")
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectNoError(t, err)
		})

		t.Run("missing value", func(t *testing.T) {
			cd := getContentDataWithDummyData(t, "societyId")
			cd.Font = ""
			ce := NewContentEntry("id", cd)
			err := ce.IsValid()
			expectError(t, err)
		})
	})

}

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

func TestAddMissingValues(t *testing.T) {

	cdEmpty := ContentData{}
	cdAllSet := getContentDataWithDummyData(t, "type")

	t.Run("fill empty set from full set", func(t *testing.T) {
		ceSrc := NewContentEntry("idAllSet", cdAllSet)
		ceDst := NewContentEntry("idEmpty", cdEmpty)

		AddMissingValues(ceSrc.data, &ceDst.data, "Presets", "ID")
		expectAllExportedSet(t, ceDst)
	})

	t.Run("do not overwrite existing data", func(t *testing.T) {
		ceSrc := NewContentEntry("src", ContentData{Desc: "srcDesc", Font: "srcFont", X1: 1.0, Y1: 2.0})
		ceDst := NewContentEntry("dst", ContentData{Desc: "dstDesc", X1: 3.0, X2: 4.0})

		AddMissingValues(ceSrc.data, &ceDst.data, "Presets", "ID")

		expectEqual(t, ceDst.Description(), "dstDesc")
		expectEqual(t, ceDst.Font(), "srcFont")
		expectEqual(t, ceDst.X1(), 3.0)
		expectEqual(t, ceDst.Y1(), 2.0)
		expectEqual(t, ceDst.X2(), 4.0)
	})

	// TODO add tests where some fields do not exist in source structure and others not in target structure
}

func TestAddContent(t *testing.T) {
	s := NewStamp(400.0, 400.0)
	expectNotNil(t, s)
	text := "foo"

	t.Run("error", func(t *testing.T) {
		t.Run("no type", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			cd.Type = ""
			ce := NewContentEntry("myId", cd)
			err := ce.AddContent(s, &text)
			expectError(t, err)
		})

		t.Run("invalid type", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			cd.Type = "unknownType"
			ce := NewContentEntry("myId", cd)
			err := ce.AddContent(s, &text)
			expectError(t, err)
		})

		// don't check concrete invalid contents, e.g. invalid textCell.
		// That is done in the specialied test functions below
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("textCell", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := ce.AddContent(s, &text)
			expectNoError(t, err)
		})
		t.Run("societyId", func(t *testing.T) {
			societyID := "123456-789"
			cd := getSocietyIDWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := ce.AddContent(s, &societyID)
			expectNoError(t, err)
		})
	})
}

func TestAddTextCell(t *testing.T) {
	s := NewStamp(400.0, 400.0)
	expectNotNil(t, s)
	text := "foo"

	t.Run("error", func(t *testing.T) {
		t.Run("missing input value", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := ce.addTextCell(s, nil)
			expectError(t, err)
		})

		t.Run("missing content", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			cd.Font = ""
			ce := NewContentEntry("myId", cd)
			err := ce.addTextCell(s, &text)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("textCell", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := ce.addTextCell(s, &text)
			expectNoError(t, err)
		})
	})
}

func TestAddSocietyID(t *testing.T) {
	s := NewStamp(400.0, 400.0)
	expectNotNil(t, s)
	text := "123456-789"

	t.Run("error", func(t *testing.T) {
		t.Run("missing input value", func(t *testing.T) {
			cd := getSocietyIDWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := ce.addSocietyID(s, nil)
			expectError(t, err)
		})

		t.Run("missing content content", func(t *testing.T) {
			cd := getSocietyIDWithDummyData()
			cd.Font = ""
			ce := NewContentEntry("myId", cd)
			err := ce.addSocietyID(s, &text)
			expectError(t, err)
		})

		t.Run("xpivot left-outside boundaries", func(t *testing.T) {
			cd := getSocietyIDWithDummyData()
			cd.XPivot = cd.X1 - 1.0
			ce := NewContentEntry("myId", cd)
			err := ce.addSocietyID(s, &text)
			expectError(t, err)
		})

		t.Run("xpivot right-outside boundaries", func(t *testing.T) {
			cd := getSocietyIDWithDummyData()
			cd.XPivot = cd.X2 + 1.0
			ce := NewContentEntry("myId", cd)
			err := ce.addSocietyID(s, &text)
			expectError(t, err)
		})

		t.Run("societyId with wrong format", func(t *testing.T) {
			for _, societyID := range []string{"", "foo", "a123-456", "123-456b", "1"} {
				cd := getSocietyIDWithDummyData()
				ce := NewContentEntry("myId", cd)
				t.Logf("Testing society id '%v'", societyID)
				err := ce.addSocietyID(s, &societyID)
				expectError(t, err)
			}
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("societyId", func(t *testing.T) {
			for _, societyID := range []string{"-", "1-", "-2", "123-456"} {
				cd := getSocietyIDWithDummyData()
				ce := NewContentEntry("myId", cd)
				t.Logf("Testing society id '%v'", societyID)
				err := ce.addSocietyID(s, &societyID)
				expectNoError(t, err)
			}
		})
	})
}
