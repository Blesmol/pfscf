package utils

import (
	"reflect"
	"testing"
)

func helperCopyValueIfUnset(t *testing.T) {

}

func TestCopyValueIfUnset(t *testing.T) {
	zero := 0.0
	one := 1.0
	two := 2.0

	testdata := []struct {
		PtrDst, PtrSrc *float64
		ptrExpNil      bool
		ptrExpResult   float64
	}{
		{nil, nil, true, zero},
		{nil, &zero, false, zero},
		{nil, &one, false, one},
		{&zero, nil, false, zero},
		{&zero, &one, false, zero},
		{&one, &two, false, one},
	}

	for _, tt := range testdata {
		ttPtr := &tt

		vTT := reflect.ValueOf(ttPtr).Elem()
		if !vTT.CanSet() {
			t.Log("Struct must be settable for this test")
			t.FailNow()
		}

		vSrc := vTT.FieldByName("PtrSrc")
		vDst := vTT.FieldByName("PtrDst")
		if !vDst.CanSet() {
			t.Log("Destination value must be settable for this test")
			t.FailNow()
		}

		copyValueIfUnset(vSrc, vDst)

		if tt.ptrExpNil {
			if tt.PtrDst != nil {
				t.Error("Should have been nil")
			}
		} else {
			if tt.PtrDst == nil {
				t.Error("Should not have been nil")
			} else if *tt.PtrDst != tt.ptrExpResult {
				t.Errorf("Values should have been identical: '%v' != '%v'", *tt.PtrDst, tt.ptrExpResult)
			} else if tt.PtrSrc == tt.PtrDst {
				t.Error("Pointers should have been different")
			}
		}
	}

}

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

		t.Run("supported elem datatypes", func(t *testing.T) {
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

		t.Run("supported ptr datatypes", func(t *testing.T) {
			srcElem := []float64{10.0, 11.0, 12.0}
			dstElem := []float64{0.0, 1.0}

			type testStruct struct {
				A, B, C, D *float64
			}
			source := testStruct{A: &srcElem[0], B: &srcElem[1], C: &srcElem[2], D: nil}
			target := testStruct{A: &dstElem[0], B: &dstElem[1], C: nil, D: nil}
			AddMissingValues(&target, source)

			if *target.A != dstElem[0] ||
				*target.B != dstElem[1] ||
				*target.C != srcElem[2] ||
				target.D != nil {
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
