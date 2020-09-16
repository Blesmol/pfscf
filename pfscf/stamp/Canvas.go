package stamp

type canvas struct {
	xPt, yPt float64
	wPt, hPt float64
}

func newCanvas(xPt, yPt, wPt, hPt float64) (c canvas) {
	c = canvas{xPt, yPt, wPt, hPt}
	return
}

func (c canvas) getSubCanvas(x1Pct, y1Pct, x2Pct, y2Pct float64) (subC canvas) {
	xPt, yPt, wPt, hPt := c.pctToPt(x1Pct, y1Pct, x2Pct, y2Pct)

	subC = canvas{xPt, yPt, wPt, hPt}
	return
}

// relPctToPt converts the provided percent coordinates into absolute
// point coordinates for the current stamp object.
// A value of, e.g. 10% should be passed as 10.0, not as 0.10
func (c canvas) relPctToPt(xPct, yPct float64) (xPt, yPt float64) {
	xPt = xPct * c.wPt / 100.0
	yPt = yPct * c.hPt / 100.0
	return
}

func (c canvas) pctToPt(x1Pct, y1Pct, x2Pct, y2Pct float64) (xPt, yPt, wPt, hPt float64) {
	xPct, yPct, wPct, hPct := getXYWH(x1Pct, y1Pct, x2Pct, y2Pct)

	xPt, yPt = c.relPctToPt(xPct, yPct)
	xPt += c.xPt // coordinates should be absolute, not relative to canvas
	yPt += c.yPt

	wPt, hPt = c.relPctToPt(wPct, hPct)

	return
}

// relPtToPct converts the provided point coordinates into percent
// coordinates for the current stamp object.
// The point coordinates are expected to be relative to the current canvas,
// i.e. a pt coordinate of 0,0 would referr to the upper left corner of
// the canvas.
// A value of, e.g. 10% will be returned as 10.0, not as 0.10
func (c canvas) relPtToPct(xPt, yPt float64) (xPct, yPct float64) {
	xPct = 100.0 * xPt / c.wPt
	yPct = 100.0 * yPt / c.hPt

	return
}

// getXYWH transforms two sets of x/y coordinates into a single set of
// x/y coordinates and a pair of width/height values.
func getXYWH(x1, y1, x2, y2 float64) (x, y, w, h float64) {
	if x1 < x2 {
		x = x1
		w = x2 - x1
	} else {
		x = x2
		w = x1 - x2
	}
	if y1 < y2 {
		y = y1
		h = y2 - y1
	} else {
		y = y2
		h = y1 - y2
	}
	return
}
