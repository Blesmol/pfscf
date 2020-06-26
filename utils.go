package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

var isTestEnvironment = false

// Assert will throw a panic if condition is false.
// The additional parameter is provided to panic() as argument.
func Assert(cond bool, i interface{}) {
	if cond == false {
		panic(i)
	}
}

// AssertNoError is a cheap function to check for err==nil and
// panics if condition is not met
func AssertNoError(err error) {
	Assert(err == nil, err)
}

// ExitOnError exits the program with rc=1 if err!=nil. The following
// Message and all additional arguments will be passed to fmt.Printf()
func ExitOnError(err error, errMsg string, v ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg+": ", v...)
		fmt.Fprintln(os.Stderr, err)
		// TODO add flag somewhere that will exit with panic() instead
		os.Exit(1)
	}
}

// InformOnError provides an error message on stdErr if an error
// occurs, but continues with the program
func InformOnError(err error, errMsg string, v ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg+":\n", v...)
		fmt.Fprintln(os.Stderr, err)
	}
}

// GetTempDir returns the location of a temporary directory in which
// files can be stored. The caller needs to ensure that the
// directory is deleted afterwards.
func GetTempDir() (dirName string) {
	// TODO Wait for watermarking issue to be fixed on side of pdfcpu
	// https://github.com/pdfcpu/pdfcpu/issues/195
	// Watermarking with pdfcpu currently does not work on Windows
	// when absolute paths are used.
	// So temporarily create the working dir as subdir of the local directory
	dirName, err := ioutil.TempDir(".", "pfsct-")
	AssertNoError(err)
	return dirName
}

// GetExecutableDir returns the dir in which the binary of this program is located
func GetExecutableDir() (dirName string) {
	if IsTestEnvironment() {
		// during test runs the executable will be run from some temporary
		// directory, so instead return the local directory for that case.
		return "."
	}

	dirName, err := filepath.Abs(filepath.Dir(os.Args[0]))
	AssertNoError(err)
	return dirName
}

// IsSet checks whether the provided value is different from its zero value and
// in case of a non-nil pointer it also checks whether the referenced value is
// not the zero value
func IsSet(val interface{}) (result bool) {
	x := reflect.ValueOf(val)
	if x.Kind() == reflect.Ptr {
		return !(x.IsNil() || x.Elem().IsZero())
	}
	return !x.IsZero()
}

// IsTestEnvironment should recognize whether the current run is a test run.
func IsTestEnvironment() bool {
	return isTestEnvironment
}

// SetIsTestEnvironment sets a flag that indicates that we are currently in
// a test environment.
func SetIsTestEnvironment() {
	isTestEnvironment = true
}
