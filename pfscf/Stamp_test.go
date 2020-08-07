package main

import (
	"os"
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

func getTextCellWithDummyData() (cd yaml.ContentData) {
	cd.Type = "textCell"
	cd.X1 = 12.0
	cd.Y1 = 12.0
	cd.X2 = 24.0
	cd.Y2 = 24.0
	cd.Font = "Helvetica"
	cd.Fontsize = 14.0
	cd.Align = "LB"

	return cd
}

func getSocietyIDWithDummyData() (cd yaml.ContentData) {
	cd.Type = "societyId"
	cd.X1 = 12.0
	cd.Y1 = 12.0
	cd.XPivot = 16.0
	cd.X2 = 24.0
	cd.Y2 = 24.0
	cd.Font = "Helvetica"
	cd.Fontsize = 14.0

	return cd
}

func TestNewStamp(t *testing.T) {
	s := NewStamp(1.0, 2.0)
	test.ExpectNotNil(t, s)

	test.ExpectEqual(t, s.dimX, 1.0)
	test.ExpectEqual(t, s.dimY, 2.0)
	test.ExpectEqual(t, s.cellBorder, "0")
}

func TestStamp_SetCellBorder(t *testing.T) {
	s := NewStamp(1.0, 1.0)
	test.ExpectNotNil(t, s)

	test.ExpectEqual(t, s.cellBorder, "0") // default is that no cell border should be drawn
	s.SetCellBorder(true)
	test.ExpectEqual(t, s.cellBorder, "1")
	s.SetCellBorder(false)
	test.ExpectEqual(t, s.cellBorder, "0")
}

func TestStamp_PctToPt(t *testing.T) {
	s := NewStamp(200.0, 200.0)

	x, y := s.pctToPt(10.0, 10.0)
	test.ExpectEqual(t, x, 20.0)
	test.ExpectEqual(t, y, 20.0)
}

func TestStamp_PtToPct(t *testing.T) {
	s := NewStamp(200.0, 200.0)

	x, y := s.ptToPct(20.0, 20.0)
	test.ExpectEqual(t, x, 10.0)
	test.ExpectEqual(t, y, 10.0)
}

func TestGetXYWH(t *testing.T) {

	for _, data := range []struct {
		desc                                   string
		x1, y1, x2, y2, xExp, yExp, wExp, hExp float64
	}{
		{"x1/y1 smaller than x2/y2", 0.0, 1.0, 100.0, 101.0, 0.0, 1.0, 100.0, 100.0},
		{"x1/y1 greater than x2/y2", 100.0, 101.0, 0.0, 1.0, 0.0, 1.0, 100.0, 100.0},
		{"x1/y1 equal to x2/y2", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0, 0.0},
	} {
		t.Logf("%v:", data.desc)
		t.Logf("  x1=%.1f, y1=%.1f, x2=%.1f, y2=%.1f", data.x1, data.y1, data.x2, data.y2)

		x, y, w, h := getXYWH(data.x1, data.y1, data.x2, data.y2)
		test.ExpectEqual(t, x, data.xExp)
		test.ExpectEqual(t, y, data.yExp)
		test.ExpectEqual(t, w, data.wExp)
		test.ExpectEqual(t, h, data.hExp)
	}
}

func TestStamp_DetermineFontSize(t *testing.T) {
	s := NewStamp(100.0, 100.0)

	var result float64

	for _, data := range []struct {
		width, fontsize  float64
		text             string
		expectedFontsize float64
	}{
		{1.0, 14.0, "fooooooooooooooooooooooo", minFontSize},
		{100.0, 14.0, "fooooooooooooooooooooooo", 7.5},
		{100.0, 14.0, "foo", 14.0},
	} {
		result = s.DeriveFontsize(data.width, "Arial", data.fontsize, data.text)
		test.ExpectEqual(t, result, data.expectedFontsize)
	}
}

func TestStamp_WriteToFile(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		t.Run("missing filename", func(t *testing.T) {
			s := NewStamp(400.0, 400.0)
			test.ExpectNotNil(t, s)
			err := s.WriteToFile("")
			test.ExpectError(t, err)
		})

		// TODO invalid filename?
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("fiii", func(t *testing.T) {
			s := NewStamp(400.0, 400.0)
			test.ExpectNotNil(t, s)
			workDir := utils.GetTempDir()
			defer os.RemoveAll(workDir)
			err := s.WriteToFile(filepath.Join(workDir, "stamp.pdf"))
			test.ExpectNoError(t, err)
		})
	})

}

func TestStamp_CreateMeasurementCoordinates(t *testing.T) {
	t.Run("with minor gap", func(t *testing.T) {
		s := NewStamp(395.0, 395.0)
		s.CreateMeasurementCoordinates(5.0, 1.0)
	})
	t.Run("without minor gap", func(t *testing.T) {
		s := NewStamp(395.0, 395.0)
		s.CreateMeasurementCoordinates(5.0, 0)
	})
}
