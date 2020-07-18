package main

import (
	"path/filepath"
)

const (
	templateDir = "templates"
)

var (
	templateTestDir string
)

// GetTemplatesDir returns the path below which the template files are stored.
// In case a test environment is recognized, a different directory with testdata is returned.
func GetTemplatesDir() (dir string) {
	if IsTestEnvironment() {
		return getTestingTemplatesDir()
	}
	return getProductiveTemplatesDir()
}

// getProductiveTemplatesDir returns the path below which the
// productive template files are stored.
func getProductiveTemplatesDir() (dir string) {
	return filepath.Join(GetExecutableDir(), templateDir)
}

// getTestingTemplatesDir returns the path below which the template files
// for tests are stored.
func getTestingTemplatesDir() (dir string) {
	return templateTestDir
}

// SetTestingTemplatesDir sets the global template dir in a testing environment
func SetTestingTemplatesDir(dir string) {
	Assert(IsTestEnvironment(), "Should only be called during tests")
	templateTestDir = dir
}
