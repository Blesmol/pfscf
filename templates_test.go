package main

import (
	"testing"
)

func getContentEntryWithDummyData(ceType string, ceID string) (ce ContentEntry) {
	ce.Type = ceType
	ce.ID = ceID
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

func TestContentEntryIsValid_emptyType(t *testing.T) {
	ce := getContentEntryWithDummyData("", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestContentEntryIsValid_invalidType(t *testing.T) {
	ce := getContentEntryWithDummyData("textCellX", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestContentEntryIsValid_validTextCell(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell", "foo")
	isValid, err := ce.IsValid()

	expectEqual(t, isValid, true)
	expectNoError(t, err)
}

func TestContentEntryIsValid_textCellWithZeroedValues(t *testing.T) {
	ce := getContentEntryWithDummyData("textCell", "foo")
	ce.Font = ""

	isValid, err := ce.IsValid()

	expectEqual(t, isValid, false)
	expectError(t, err)
}

func TestApplyDefaults(t *testing.T) {
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
	expectNotSet(t, ce.ID)
	expectNotSet(t, ce.Desc)
	expectEqual(t, ce.X1, 5.0)
	expectNotSet(t, ce.Y1)
	expectNotSet(t, ce.X2)
	expectNotSet(t, ce.Y2)
	expectEqual(t, ce.Font, "Dingbats")
	expectEqual(t, ce.Fontsize, 10.0)
	expectNotSet(t, ce.Align)
}
