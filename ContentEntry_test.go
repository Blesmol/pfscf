package main

import "testing"

func init() {
	SetIsTestEnvironment()
}

func getContentEntryWithDummyData(ceType string) (ce ContentEntry) {
	ce.Type = ceType
	ce.Desc = "Some Description"
	ce.X1 = 12.0
	ce.Y1 = 12.0
	ce.X2 = 24.0
	ce.Y2 = 24.0
	ce.Font = "Helvetica"
	ce.Fontsize = 14.0
	ce.Align = "LB"
	return ce
}

func Test_ApplyDefaults(t *testing.T) {
	var ce ContentEntry
	ce.Type = "Foo"
	ce.Fontsize = 10.0
	ce.Font = ""

	var defaults ContentEntry
	defaults.Type = "Bar"
	defaults.X1 = 5.0
	defaults.Y1 = 0.0
	defaults.Font = "Dingbats"
	defaults.Fontsize = 20.0

	ce.applyDefaults(defaults)

	expectEqual(t, ce.Type, "Foo")
	expectNotSet(t, ce.Desc)
	expectEqual(t, ce.X1, 5.0)
	expectNotSet(t, ce.Y1)
	expectNotSet(t, ce.X2)
	expectNotSet(t, ce.Y2)
	expectEqual(t, ce.Font, "Dingbats")
	expectEqual(t, ce.Fontsize, 10.0)
	expectNotSet(t, ce.Align)
}

func Test_ContentEntryIsValid_emptyType(t *testing.T) {
	ce := getContentEntryWithDummyData("")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func Test_ContentEntryIsValid_invalidType(t *testing.T) {
	ce := getContentEntryWithDummyData("textCellX")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func Test_ContentEntryIsValid_validTextCell(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, true)
	expectNoError(t, err)
}

func Test_ContentEntryIsValid_textCellWithZeroedValues(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell")
	ce.Font = ""

	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}
