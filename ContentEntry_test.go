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

	expectAllExportedSet(t, cd) // to be sure that we also get all new fields

	return cd
}

func Test_NewContentEntry(t *testing.T) {
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

func Test_ContentEntryIsValid_invalidType(t *testing.T) {
	cd := getContentDataWithDummyData(t, "textCellX")
	ce := NewContentEntry("id", cd)
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func Test_ContentEntryIsValid_validTextCell(t *testing.T) {
	cd := getContentDataWithDummyData(t, "textCell")
	ce := NewContentEntry("id", cd)
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, true)
	expectNoError(t, err)
}

func Test_ContentEntryIsValid_textCellWithZeroedValues(t *testing.T) {
	cd := getContentDataWithDummyData(t, "textCell")
	cd.Font = ""
	ce := NewContentEntry("id", cd)

	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func Test_EntriesAreNotContradicting(t *testing.T) {
	var err error

	cdEmpty := ContentData{}
	cdAllSet := getContentDataWithDummyData(t, "type")

	ceEmpty := NewContentEntry("idEmpty", cdEmpty)
	ceAllSet := NewContentEntry("idAllSet", cdAllSet)

	// a given CE with values should not contradict itself
	err = EntriesAreNotContradicting(&ceAllSet, &ceAllSet)
	expectNoError(t, err)

	// a given CE with no values should contradict nothing
	err = EntriesAreNotContradicting(&ceEmpty, &ceEmpty)
	expectNoError(t, err)
	err = EntriesAreNotContradicting(&ceAllSet, &ceEmpty)
	expectNoError(t, err)
	err = EntriesAreNotContradicting(&ceEmpty, &ceAllSet)
	expectNoError(t, err)

	// Have to partly-set objects with non-overlapping content
	{
		cdLeft := ContentData{X1: 1.0, Desc: "desc"}
		ceLeft := NewContentEntry("idLeft", cdLeft)
		cdRight := ContentData{X2: 2.0, Font: "font"}
		ceRight := NewContentEntry("idRight", cdRight)
		err = EntriesAreNotContradicting(&ceLeft, &ceRight)
		expectNoError(t, err)
	}

	// produce some conflict in a string attribute
	{
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Font = cdLeft.Font + "foo" // <= conflicting data
		ceLeft := NewContentEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		ceRight := NewContentEntry("idRight", cdRight)

		err = EntriesAreNotContradicting(&ceLeft, &ceRight)
		expectError(t, err)
	}

	// produce some conflict in a float64 attribute
	{
		cdLeft := getContentDataWithDummyData(t, "type")
		cdLeft.Fontsize = cdLeft.Fontsize + 1.0 // <= conflicting data
		ceLeft := NewContentEntry("idLeft", cdLeft)
		cdRight := getContentDataWithDummyData(t, "type")
		ceRight := NewContentEntry("idRight", cdRight)

		err = EntriesAreNotContradicting(&ceLeft, &ceRight)
		expectError(t, err)
	}
}

func Test_ContentEntryAddMissingValuesFromOther(t *testing.T) {
	//var err error

	cdEmpty := ContentData{}
	cdAllSet := getContentDataWithDummyData(t, "type")

	// set empty from full
	{
		ceSrc := NewContentEntry("idAllSet", cdAllSet)
		ceDst := NewContentEntry("idEmpty", cdEmpty)

		ceDst.AddMissingValuesFrom(&ceSrc)
		expectAllExportedSet(t, ceDst)
	}

	// do not overwrite existing data
	{
		ceSrc := NewContentEntry("src", ContentData{Desc: "srcDesc", Font: "srcFont", X1: 1.0, Y1: 2.0})
		ceDst := NewContentEntry("dst", ContentData{Desc: "dstDesc", X1: 3.0, X2: 4.0})
		ceDst.AddMissingValuesFrom(&ceSrc)

		expectEqual(t, ceDst.Description(), "dstDesc")
		expectEqual(t, ceDst.Font(), "srcFont")
		expectEqual(t, ceDst.X1(), 3.0)
		expectEqual(t, ceDst.Y1(), 2.0)
		expectEqual(t, ceDst.X2(), 4.0)
	}

}
