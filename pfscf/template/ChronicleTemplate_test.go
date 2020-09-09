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
  foo:
    type: text
    description: d
    example: foo

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
