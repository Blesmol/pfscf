package main

import (
	"reflect"
	"runtime/debug"
	"testing"
)

const (
	printCallStackOnFailingTest = true
)

func callStack() {
	if printCallStackOnFailingTest {
		debug.PrintStack()
	}
}

func expectEqual(t *testing.T, got interface{}, exp interface{}) {
	if exp == got {
		return
	}
	callStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

func expectNotEqual(t *testing.T, got interface{}, notExp interface{}) {
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
	if !reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if reflect.ValueOf(got).IsNil() {
		callStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		callStack()
		t.Error("Expected an error, got nil")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		callStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func expectNotSet(t *testing.T, got interface{}) {
	if IsSet(got) {
		callStack()
		t.Errorf("Expected not set, got '%v'", got)
	}
}

func expectAllSet(t *testing.T, got interface{}) {
	vGot := reflect.ValueOf(got)

	switch vGot.Kind() {
	case reflect.Struct:
		for i := 0; i < vGot.NumField(); i++ {
			field := vGot.Field(i)
			expectAllSet(t, field.Interface())
		}
	case reflect.Ptr:
		if IsSet(got) {
			expectAllSet(t, vGot.Elem().Interface())
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

func expectKeyExists(t *testing.T, tmap interface{}, key interface{}) {

}
