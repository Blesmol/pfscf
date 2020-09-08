package content

import (
	"testing"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/preset"
	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"gopkg.in/yaml.v2"
)

func getTestPresetStore(t *testing.T) (ps *preset.Store) {
	presetContent := `
sameData1:
  x: 10.0
  font: Helvetica
sameData2:
  x: 10.0
  font: Helvetica
conflict1:
  font: Arial
conflict2:
  font: Helvetica
`
	var store preset.Store
	err := yaml.Unmarshal([]byte(presetContent), &store)
	test.ExpectNoError(t, err)

	return &store
}

func getTestArgStore(key, value string) (as *args.Store) {
	as, _ = args.NewStore(args.StoreInit{})
	as.Set(key, value)
	return as
}
