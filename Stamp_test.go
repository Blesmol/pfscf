package main

import (
	"os"
	"path/filepath"
	"testing"
)

func getTextCellWithDummyData() (cd ContentData) {
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

func TestNewStamp(t *testing.T) {
	s := NewStamp(1.0, 2.0)
	expectNotNil(t, s)

	expectEqual(t, s.dimX, 1.0)
	expectEqual(t, s.dimY, 2.0)
	expectEqual(t, s.cellBorder, "0")
}

func TestSetCellBorder(t *testing.T) {
	s := NewStamp(1.0, 1.0)
	expectNotNil(t, s)

	expectEqual(t, s.cellBorder, "0") // default is that no cell border should be drawn
	s.SetCellBorder(true)
	expectEqual(t, s.cellBorder, "1")
	s.SetCellBorder(false)
	expectEqual(t, s.cellBorder, "0")
}

func TestGetXYWH(t *testing.T) {
	t.Run("x1/y1 smaller than x2/y2", func(t *testing.T) {
		x, y, w, h := getXYWH(0.0, 1.0, 100.0, 101.0)
		expectEqual(t, x, 0.0)
		expectEqual(t, y, 1.0)
		expectEqual(t, w, 100.0)
		expectEqual(t, h, 100.0)
	})

	t.Run("x1/y1 greater than x2/y2", func(t *testing.T) {
		x, y, w, h := getXYWH(100.0, 101.0, 0.0, 1.0)
		expectEqual(t, x, 0.0)
		expectEqual(t, y, 1.0)
		expectEqual(t, w, 100.0)
		expectEqual(t, h, 100.0)
	})

	t.Run("x1/y1 equal to x2/y2", func(t *testing.T) {
		x, y, w, h := getXYWH(1.0, 1.0, 1.0, 1.0)
		expectEqual(t, x, 1.0)
		expectEqual(t, y, 1.0)
		expectEqual(t, w, 0.0)
		expectEqual(t, h, 0.0)
	})
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
			err := s.AddContent(ce, &text)
			expectError(t, err)
		})

		t.Run("invalid type", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			cd.Type = "unknownType"
			ce := NewContentEntry("myId", cd)
			err := s.AddContent(ce, &text)
			expectError(t, err)
		})

		// don't check concrete invalid contents, e.g. invalid textCell.
		// That is done in the specialied test functions below
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("textCell", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := s.AddContent(ce, &text)
			expectNoError(t, err)
		})
	})
}

func TestAddTextCell(t *testing.T) {
	s := NewStamp(400.0, 400.0)
	expectNotNil(t, s)
	text := "foo"

	t.Run("error", func(t *testing.T) {
		t.Run("no input value", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := s.addTextCell(ce, nil)
			expectError(t, err)
		})

		t.Run("invalid content", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			cd.Font = ""
			ce := NewContentEntry("myId", cd)
			err := s.addTextCell(ce, &text)
			expectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("textCell", func(t *testing.T) {
			cd := getTextCellWithDummyData()
			ce := NewContentEntry("myId", cd)
			err := s.addTextCell(ce, &text)
			expectNoError(t, err)
		})
	})
}

func TestWriteToFile(t *testing.T) {

	t.Run("error", func(t *testing.T) {
		t.Run("missing filename", func(t *testing.T) {
			s := NewStamp(400.0, 400.0)
			expectNotNil(t, s)
			err := s.WriteToFile("")
			expectError(t, err)
		})

		// TODO invalid filename?
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("fiii", func(t *testing.T) {
			s := NewStamp(400.0, 400.0)
			expectNotNil(t, s)
			workDir := GetTempDir()
			defer os.RemoveAll(workDir)
			err := s.WriteToFile(filepath.Join(workDir, "stamp.pdf"))
			expectNoError(t, err)
		})
	})

}

func TestCreateMeasurementCoordinates(t *testing.T) {
	t.Run("with minor gap", func(t *testing.T) {
		s := NewStamp(395.0, 395.0)
		s.CreateMeasurementCoordinates(100.0, 25.0)
	})
	t.Run("without minor gap", func(t *testing.T) {
		s := NewStamp(395.0, 395.0)
		s.CreateMeasurementCoordinates(100.0, 0)
	})
}
