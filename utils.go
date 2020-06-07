package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	dirName, err := filepath.Abs(filepath.Dir(os.Args[0]))
	AssertNoError(err)
	return dirName
}
