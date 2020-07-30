package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
)

const (
	printCallStackOnFailingTest = false
)

var (
	isTestEnvironment = false
)

func callStack() {
	if printCallStackOnFailingTest {
		debug.PrintStack()
	}
}

// IsTestEnvironment should indicate whether the current run is a test run.
func IsTestEnvironment() bool {
	return isTestEnvironment
}

// SetIsTestEnvironment sets a flag that indicates that we are currently in
// a test environment.
func SetIsTestEnvironment(isTestEnv bool) {
	isTestEnvironment = isTestEnv
}

func expectEqual(t *testing.T, got interface{}, exp interface{}) {
	t.Helper()

	if exp == got {
		return
	}

	callStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

func expectNotEqual(t *testing.T, got interface{}, notExp interface{}) {
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

func expectNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	t.Helper()

	if reflect.TypeOf(got) != nil && !reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	t.Helper()

	if reflect.TypeOf(got) == nil || reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

// expectError checks that the provided error does not equal nil. If additional
// string arguments were passend, then it is checked that the error message
// contains all of them.
func expectError(t *testing.T, err error, expContent ...string) {
	t.Helper()

	if err == nil {
		callStack()
		t.Error("Expected an error, got nil")
	}

	for _, expPartialError := range expContent {
		if !strings.Contains(err.Error(), expPartialError) {
			callStack()
			t.Errorf("Expected string '%v' to be contained in error message:\n%v", expPartialError, err.Error())
			return
		}
	}
}

func expectNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		callStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func expectNotSet(t *testing.T, got interface{}) {
	t.Helper()

	if IsSet(got) {
		callStack()
		t.Errorf("Expected not set, got '%v'", got)
	}
}

func expectIsSet(t *testing.T, got interface{}) {
	t.Helper()

	if !IsSet(got) {
		callStack()
		t.Errorf("Expected value of type '%v' to be set, but was not", reflect.TypeOf(got))
	}
}

func expectAllExportedSet(t *testing.T, got interface{}) {
	t.Helper()

	vGot := reflect.ValueOf(got)
	switch vGot.Kind() {
	case reflect.Struct:
		for i := 0; i < vGot.NumField(); i++ {
			field := vGot.Field(i)
			if !IsExported(field) {
				continue // skip non-exported fields
			}
			t.Logf("Testing field '%v'", reflect.TypeOf(got).Field(i).Name)
			expectAllExportedSet(t, field.Interface())
		}
	case reflect.Ptr:
		if IsSet(got) {
			expectAllExportedSet(t, vGot.Elem().Interface())
		} else {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	default:
		if !IsSet(got) {
			callStack()
			t.Errorf("Expected to be set, but was not: %v / %v", vGot.Type(), vGot.Kind())
		}
	}
}

func expectFileExists(t *testing.T, filename string) {
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

func expectKeyExists(t *testing.T, tMap interface{}, key interface{}) {
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

func expectTrue(t *testing.T, v bool) {
	t.Helper()

	if !v {
		callStack()
		t.Errorf("Expected true, but was false")
	}
}

func expectFalse(t *testing.T, v bool) {
	t.Helper()

	if v {
		callStack()
		t.Errorf("Expected false, but was true")
	}
}

func expectStringContains(t *testing.T, got string, exp string) {
	t.Helper()

	if !strings.Contains(got, exp) {
		callStack()
		t.Errorf("Expected string '%v' to contain '%v', which it does not", got, exp)
	}
}

func expectStringContainsNot(t *testing.T, got string, exp string) {
	t.Helper()

	if strings.Contains(got, exp) {
		callStack()
		t.Errorf("Expected string '%v' to NOT contain '%v', but it does", got, exp)
	}
}
