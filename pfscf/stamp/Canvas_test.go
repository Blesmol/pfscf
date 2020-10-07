package stamp

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func TestCanvas_relPctToPt(t *testing.T) {
	testData := []struct {
		c              canvas
		xPct, yPct     float64
		xPtExp, yPtExp float64
	}{
		{newCanvas(0.0, 0.0, 200.0, 200.0), 10.0, 10.0, 20.0, 20.0},
		{newCanvas(100.0, 100.0, 200.0, 200.0), 10.0, 10.0, 20.0, 20.0},
	}

	for _, tt := range testData {
		xPtGot, yPtGot := tt.c.pctToRelPt(tt.xPct, tt.yPct)
		test.ExpectEqual(t, xPtGot, tt.xPtExp)
		test.ExpectEqual(t, yPtGot, tt.yPtExp)
	}
}

func TestCanvas_relPtToPct(t *testing.T) {
	testData := []struct {
		c                canvas
		xPt, yPt         float64
		xPctExp, yPctExp float64
	}{
		{newCanvas(0.0, 0.0, 200.0, 200.0), 20.0, 20.0, 10.0, 10.0},
		{newCanvas(100.0, 100.0, 200.0, 200.0), 10.0, 10.0, 5.0, 5.0},
	}

	for _, tt := range testData {
		xPtGot, yPtGot := tt.c.relPtToPct(tt.xPt, tt.yPt)
		test.ExpectEqual(t, xPtGot, tt.xPctExp)
		test.ExpectEqual(t, yPtGot, tt.yPctExp)
	}
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
