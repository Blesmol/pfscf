package cfg

import (
	"path/filepath"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	templateDir = "templates"
)

var (
	templateTestDir string

	// Global holds global config flags
	Global globalFlags
)

type globalFlags struct {
	Verbose        bool
	DrawCellBorder bool
	DrawGrid       bool
}

func init() {
	Global = globalFlags{
		Verbose:        false,
		DrawCellBorder: false,
		DrawGrid:       false,
	}
}

// GetTemplatesDir returns the path below which the template files are stored.
// In case a test environment is recognized, a different directory with testdata is returned.
func GetTemplatesDir() (dir string) {
	if utils.IsTestEnvironment() {
		return getTestingTemplatesDir()
	}
	return getProductiveTemplatesDir()
}

// getProductiveTemplatesDir returns the path below which the
// productive template files are stored.
func getProductiveTemplatesDir() (dir string) {
	return filepath.Join(utils.GetExecutableDir(), templateDir)
}

// getTestingTemplatesDir returns the path below which the template files
// for tests are stored.
func getTestingTemplatesDir() (dir string) {
	return templateTestDir
}

// SetTestingTemplatesDir sets the global template dir in a testing environment
func SetTestingTemplatesDir(dir string) {
	utils.Assert(utils.IsTestEnvironment(), "Should only be called during tests")
	templateTestDir = dir
}
