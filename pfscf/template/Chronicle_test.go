package template

import (
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	"github.com/Blesmol/pfscf/pfscf/utils"

	"gopkg.in/yaml.v2"
)

var (
	chronicleTemplateTestDir string
)

func init() {
	utils.SetIsTestEnvironment(true)
	chronicleTemplateTestDir = filepath.Join(utils.GetExecutableDir(), "testdata")
}

func TestChronicleTemplate_inheritFrom(t *testing.T) {
	t.Run("template with empty sections inherits", func(t *testing.T) {
		parentYaml := `
id: parent
description: some description

parameters:
  group1:
    foo:
      type: text
      description: d
      example: foo

canvas:
  page:
    x: 0.0
    y: 0.0
    x2: 100.0
    y2: 100.0

presets:
  bar:
    x2: 1.0
    y2: 1.0

content:
- type: rectangle
  presets: [bar]
  color: green
`

		childYaml := `
id: child
description: some description
inherit: parent

preset:

parameters:

content:

canvas:
`

		parentTemplate := NewChronicleTemplate("parent.yml")
		err := yaml.Unmarshal([]byte(parentYaml), &parentTemplate)
		test.ExpectNoError(t, err)

		childTemplate := NewChronicleTemplate("child.yml")
		err = yaml.Unmarshal([]byte(childYaml), &childTemplate)
		childTemplate.ensureStoresAreInitialized() // this feeld so dirty...
		test.ExpectNoError(t, err)

		err = childTemplate.inheritFrom(&parentTemplate)
		test.ExpectNoError(t, err)
	})
}

func TestParseAspectRatio(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		testData := []struct{ input, errString string }{
			{":", "does not follow pattern"},
			{"1:", "does not follow pattern"},
			{":1", "does not follow pattern"},
			{"1:asdsa", "does not follow pattern"},
		}

		for _, tt := range testData {
			t.Logf("Testing input '%v'", tt.input)
			_, _, err := parseAspectRatio(tt.input)
			test.ExpectError(t, err, tt.errString)
		}
	})

	t.Run("valid", func(t *testing.T) {
		testData := []struct {
			input      string
			xExp, yExp float64
		}{
			{"1:2", 1.0, 2.0},
			{"1.:2.", 1.0, 2.0},
			{"1.23:2.34", 1.23, 2.34},
			{"  1.23  :   2.34   ", 1.23, 2.34},
		}

		for _, tt := range testData {
			t.Logf("Testing input '%v'", tt.input)
			x, y, err := parseAspectRatio(tt.input)
			test.ExpectNoError(t, err)
			test.ExpectEqual(t, x, tt.xExp)
			test.ExpectEqual(t, y, tt.yExp)
		}
	})
}
