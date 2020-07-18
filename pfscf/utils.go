package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

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

// ExitWithMessage prints out a message on StdErr and then exits the
// program with rc=1
func ExitWithMessage(errMsg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+errMsg, v...)
	os.Exit(1)
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
	dirName, err := ioutil.TempDir("", "pfscf-")
	AssertNoError(err)
	return dirName
}

// GetExecutableDir returns the dir in which the binary of this program is located
func GetExecutableDir() (dirName string) {
	var baseDir string

	if IsTestEnvironment() {
		// during test runs the executable will be run from some temporary
		// directory, so instead return the local directory for that case.
		baseDir = "."
	} else {
		baseDir = filepath.Dir(os.Args[0])
	}

	dirName, err := filepath.Abs(baseDir)
	AssertNoError(err)
	return dirName
}

// IsFile checks wether a file exists and is not a directory.
func IsFile(filename string) (exists bool, err error) {
	info, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("Path %v is a directory", filename)
	}
	return true, nil
}

// IsDir checks wether a directory exists.
func IsDir(dirname string) (exists bool, err error) {
	info, err := os.Stat(dirname)
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("Path %v is not a directory", dirname)
	}
	return true, nil
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

// HasWhitespace checks if a string includes at least one whitespace character
func HasWhitespace(s string) bool {
	return strings.Index(s, " ") > -1
}

// QuoteStringIfRequired takes a string as input. If the string contains whitespace,
// then it is returned enclosed by quotation marks, else it is returned as-is
func QuoteStringIfRequired(s string) string {
	if HasWhitespace(s) {
		return fmt.Sprintf("\"%v\"", s)
	}
	return s
}
