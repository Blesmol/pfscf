package testutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	util "github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	printCallStackOnFailingTest = false
)

func callStack() {
	if printCallStackOnFailingTest {
		debug.PrintStack()
	}
}

// ExpectEqual expects the provided values to be equal.
func ExpectEqual(t *testing.T, got interface{}, exp interface{}) {
	t.Helper()

	if exp == got {
		return
	}

	// If got is ptr to type t and exp is element of type t, then compare elements instead
	typeGot := reflect.TypeOf(got)
	typeExp := reflect.TypeOf(exp)
	if typeGot != nil && typeExp != nil { // beware of nil types
		// ensure that elem typ of got matches the type of exp
		if typeGot.Kind() == reflect.Ptr && typeGot.Elem().Kind() == typeExp.Kind() {
			vGot := reflect.ValueOf(got)
			if !vGot.IsNil() {
				vGotElem := vGot.Elem()
				vExp := reflect.ValueOf(exp)
				if vGotElem.Interface() == vExp.Interface() {
					return
				}
			}
		}
	}

	callStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

// ExpectNotEqual expects the provided values to not be equal.
func ExpectNotEqual(t *testing.T, got interface{}, notExp interface{}) {
	t.Helper()

	typeNotExp := reflect.TypeOf(notExp)
	typeGot := reflect.TypeOf(got)

	// we always require that both types are identical.
	// Without that, testing can be a real pain
	if typeNotExp != typeGot {
		callStack()
		t.Errorf("Types do not match! Expected '%v', got '%v'", typeNotExp, typeGot)
		return
	}

	if notExp == got {
		callStack()
		t.Errorf("Expected something different than '%v' (type %v)", notExp, typeNotExp)
	}
}

// ExpectNil expects the provided argument to be nil.
func ExpectNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	t.Helper()

	if reflect.TypeOf(got) != nil && !reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

// ExpectNotNil expects the provided argument to not be nil.
func ExpectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	t.Helper()

	if reflect.TypeOf(got) == nil || reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

// ExpectError checks that the provided error does not equal nil. If additional
// string arguments were passend, then it is checked that the error message
// contains all of them.
func ExpectError(t *testing.T, err error, expContent ...string) {
	t.Helper()

	if err == nil {
		callStack()
		t.Error("Expected an error, got nil")
		return
	}

	for _, expPartialError := range expContent {
		if !strings.Contains(err.Error(), expPartialError) {
			callStack()
			t.Errorf("Expected string '%v' to be contained in error message:\n%v", expPartialError, err.Error())
			return
		}
	}
}

// ExpectNoError expects that the provided error argument is nil.
func ExpectNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		callStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

// ExpectNotSet expects that the provided argument is not set.
func ExpectNotSet(t *testing.T, got interface{}) {
	t.Helper()

	if util.IsSet(got) {
		callStack()
		t.Errorf("Expected not set, got '%v'", got)
	}
}

// ExpectIsSet expects that the provided argument is set.
func ExpectIsSet(t *testing.T, got interface{}) {
	t.Helper()

	if !util.IsSet(got) {
		callStack()
		t.Errorf("Expected value of type '%v' to be set, but was not", reflect.TypeOf(got))
	}
}

// ExpectAllExportedSet expects that all exported fields in the passed struct argument are set.
func ExpectAllExportedSet(t *testing.T, got interface{}) {
	t.Helper()

	vGot := reflect.ValueOf(got)
	switch vGot.Kind() {
	case reflect.Struct:
		for i := 0; i < vGot.NumField(); i++ {
			field := vGot.Field(i)
			if !util.IsExported(field) {
				continue // skip non-exported fields
			}
			t.Logf("Testing field '%v'", reflect.TypeOf(got).Field(i).Name)
			ExpectAllExportedSet(t, field.Interface())
		}
	case reflect.Ptr:
		if util.IsSet(got) {
			ExpectAllExportedSet(t, vGot.Elem().Interface())
		} else {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	default:
		if !util.IsSet(got) {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	}
}

// ExpectFileExists expects that the provided file name references an existing file.
func ExpectFileExists(t *testing.T, filename string) {
	t.Helper()

	info, err := os.Stat(filename)
	if err != nil {
		t.Errorf("Expected file '%v' is missing: %v", filename, err)

		// for debugging reasons, provide all filenames from the containing directory as well
		dirname := filepath.Dir(filename)
		files, errDir := ioutil.ReadDir(dirname)
		if errDir != nil {
			t.Logf("Cannot read dir '%v' to analyze issue: %v", dirname, errDir)
			return
		}
		t.Logf("Files in directory %v:", dirname)
		for _, file := range files {
			t.Logf("- %v\n", file.Name())
		}
	} else if info.IsDir() {
		t.Errorf("Expected file '%v' is a directory", filename)
	}
}

// ExpectKeyExists expects that the provided key exists in the provided map.
func ExpectKeyExists(t *testing.T, tMap interface{}, key interface{}) {
	t.Helper()

	vMap := reflect.ValueOf(tMap)
	if vMap.Kind() != reflect.Map {
		panic("Only maps should be provided here")
	}

	keyKind := reflect.ValueOf(key).Kind()

	mapVKeys := vMap.MapKeys()
	for _, vKey := range mapVKeys {
		if keyKind != vKey.Kind() {
			t.Errorf("Key kinds do not match! '%v' vs '%v'", keyKind, vKey.Kind())
			return
		}
		if key == vKey.Interface() {
			return
		}
	}

	callStack()
	t.Errorf("Key '%v' was not found in map '%v'", key, tMap)
}

// ExpectTrue expects that the provided bool argument is true.
func ExpectTrue(t *testing.T, v bool) {
	t.Helper()

	if !v {
		callStack()
		t.Errorf("Expected true, but was false")
	}
}

// ExpectFalse expects that the provided boll argument is false.
func ExpectFalse(t *testing.T, v bool) {
	t.Helper()

	if v {
		callStack()
		t.Errorf("Expected false, but was true")
	}
}

// ExpectStringContains expects that the provided string contains the provided substring.
func ExpectStringContains(t *testing.T, got string, exp string) {
	t.Helper()

	if !strings.Contains(got, exp) {
		callStack()
		t.Errorf("Expected string '%v' to contain '%v', which it does not", got, exp)
	}
}

// ExpectStringContainsNot expects that the provided string does not contain the provided substring.
func ExpectStringContainsNot(t *testing.T, got string, exp string) {
	t.Helper()

	if strings.Contains(got, exp) {
		callStack()
		t.Errorf("Expected string '%v' to NOT contain '%v', but it does", got, exp)
	}
}
