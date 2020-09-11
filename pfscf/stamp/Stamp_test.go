package stamp

import (
	"os"
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

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

func TestStamp_DetermineFontSize(t *testing.T) {
	s := NewStamp(100.0, 100.0)

	var result float64

	for _, tt := range []struct {
		width, height, fontsize float64
		text                    string
		expectedFontsize        float64
	}{
		{1.0, 16.0, 14.0, "fooooooooooooooooooooooo", minFontSize},
		{100.0, 16.0, 14.0, "fooooooooooooooooooooooo", 7.5},
		{100.0, 16.0, 14.0, "foo", 14.0},
		{100.0, 10.0, 14.0, "foo", 10.0},
		{100.0,  2.0, 14.0, "foo", minFontSize},
	} {
		result = s.DeriveFontsize(tt.width, tt.height, "Arial", tt.fontsize, tt.text)
		test.ExpectEqual(t, result, tt.expectedFontsize)
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
