package template

import (
	"path/filepath"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	chronicleTemplateTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	chronicleTemplateTestDir = filepath.Join(utils.GetExecutableDir(), "testdata")
}
