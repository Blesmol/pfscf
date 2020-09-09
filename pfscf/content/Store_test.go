package content

import (
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
)

func TestStore_InheritFrom(t *testing.T) {
	t.Run("ensure deep copies", func(t *testing.T) {
		var tc textCell
		tc.X = 1.0
		tc.Presets = []string{"foo", "bar"}

		s1 := NewStore()
		s1.add(&tc)

		s2 := NewStore()
		s2.InheritFrom(s1)

		tcStore1 := s1[0].(*textCell)
		tcStore2 := s2[0].(*textCell)

		// modifications in an entry from one store should not be reflected in the other
		tcStore1.X = 2.0
		test.ExpectNotEqual(t, tcStore1.X, tcStore2.X)

		tcStore1.Presets[0] = "foobar"
		test.ExpectNotEqual(t, tcStore1.Presets[0], tcStore2.Presets[0])
	})
}
