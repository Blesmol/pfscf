package content

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func TestRectangle_deepCopy(t *testing.T) {
	e1 := newRectangle()
	e1.X = 1.0
	e1.Presets = append(e1.Presets, "t1")

	e2 := e1.deepCopy().(*rectangle)
	e2.X = 2.0
	e2.Presets[0] = "t2"

	test.ExpectNotEqual(t, e1.X, e2.X)
	test.ExpectNotEqual(t, e1.Presets[0], e2.Presets[0])
}
