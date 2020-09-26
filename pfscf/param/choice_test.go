package param

import (
	"fmt"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"

	"gopkg.in/yaml.v2"
)

func Test_Unmarshal(t *testing.T) {
	baseInput := `
"group":
  "entry":
    type: choice
    example: foo
    description: some desc
    choices: %v
`

	for index, tt := range []struct {
		choicesString        string
		expectUnmarshalError bool
		expectValid          bool
		choices              []string
	}{
		{"", false, false, []string{}},
		{"foo", true, false, []string{}},
		{"[]", false, false, []string{}},
		{"[ foo ]", false, true, []string{"foo"}},
		{"[foo, bar]", false, true, []string{"foo", "bar"}},
	} {
		t.Logf("Testing case %v", index)

		input := fmt.Sprintf(baseInput, tt.choicesString)

		var s Store
		s = NewStore()
		err := yaml.Unmarshal([]byte(input), &s)

		if tt.expectUnmarshalError {
			test.ExpectError(t, err)
		} else {
			test.ExpectNoError(t, err)

			e, exists := s.Get("entry")
			test.ExpectTrue(t, exists)
			test.ExpectNotNil(t, e)
			if e == nil {
				t.FailNow()
			}

			test.ExpectEqual(t, e.Type(), typeChoice)
			test.ExpectEqual(t, e.ID(), "entry")
			test.ExpectEqual(t, e.Group(), "group")
			test.ExpectEqual(t, e.Description(), "some desc")
			test.ExpectEqual(t, e.Example(), "foo")

			err = e.isValid()
			if tt.expectValid {
				test.ExpectNoError(t, err)
			} else {
				test.ExpectError(t, err, "Missing choices")
			}
		}
	}
}
