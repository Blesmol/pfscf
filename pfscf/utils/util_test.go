package utils

import (
	"testing"
)

func TestAddMissingValues(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("exported and unexported fields", func(t *testing.T) {
			type testStruct struct{ A, b, C, d, E float64 }
			source := testStruct{A: 1.0, b: 2.0, C: 3.0, d: 4.0}
			target := testStruct{A: 10.0, b: 11.0}
			AddMissingValues(&target, source)

			if target.A != 10.0 ||
				target.b != 11.0 ||
				target.C != 3.0 ||
				target.d != 0.0 ||
				target.E != 0.0 {
				t.Errorf("Result was different than expected")
			}
		})

		t.Run("supported datatypes", func(t *testing.T) {
			type testStruct struct {
				A, B, C float64
				D, E, F string
			}
			source := testStruct{A: 1.0, B: 2.0 /* C left empty */, D: "4.0", E: "5.0" /* F left empty*/}
			target := testStruct{A: 10.0, D: "14.0"}
			AddMissingValues(&target, source)

			if target.A != 10.0 ||
				target.B != 2.0 ||
				target.C != 0.0 ||
				target.D != "14.0" ||
				target.E != "5.0" ||
				target.F != "" {
				t.Errorf("Result was different than expected")
			}
		})

		t.Run("different exported fields", func(t *testing.T) {
			source := struct {
				Common, OnlySource float64
			}{
				Common: 1.0, OnlySource: 2.0,
			}

			target := struct {
				Common, OnlyTarget float64
			}{}

			AddMissingValues(&target, source)

			if target.Common != 1.0 ||
				target.OnlyTarget != 0.0 {
				t.Errorf("Result was different than expected")
			}
		})

		t.Run("ignore fields", func(t *testing.T) {
			type testStruct struct {
				A, B, C, D float64
			}
			source := testStruct{A: 1.0, B: 2.0, C: 3.0, D: 4.0}
			target := testStruct{}

			AddMissingValues(&target, source, "B", "C", "a", "De")

			if target.A != 1.0 ||
				target.B != 0.0 ||
				target.C != 0.0 ||
				target.D != 4.0 {
				t.Errorf("Result was different than expected")
			}
		})
	})
}
