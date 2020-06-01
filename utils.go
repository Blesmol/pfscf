package main

import (
	"fmt"
	"io/ioutil"
)

// Assert duplicates C assert()
func Assert(cond bool, err error) {
	if cond == false {
		fmt.Printf("Error is %v\n", err)
		panic(err)
	}
}

// AssertNoError is a cheap function to check for err==nil and
// panics if condition is not met
func AssertNoError(err error) {
	Assert(err == nil, err)
}

// GetTempDir returns the location of a temporary directory in which
// files can be stored. The caller needs to ensure that the
// directory is deleted afterwards.
func GetTempDir() (name string) {
	// TODO Wait for watermarking issue to be fixed on side of pdfcpu
	// https://github.com/pdfcpu/pdfcpu/issues/195
	// Watermarking with pdfcpu currently does not work on Windows
	// when absolute paths are used.
	// So temporarily create the working dir as subdir of the local directory
	name, err := ioutil.TempDir(".", "pfsct-")
	AssertNoError(err)
	return name
}
