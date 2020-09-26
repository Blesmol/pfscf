package param

import (
	"fmt"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeText = "text"
)

type textEntry struct {
	id      string
	theRank int
	group   string

	TheExample     string `yaml:"example"`
	TheDescription string `yaml:"description"`
}

func (e *textEntry) Type() string {
	return typeText
}

func (e *textEntry) ID() string {
	return e.id
}

func (e *textEntry) setID(id string) {
	utils.Assert(!utils.IsSet(e.id), "Should only be called once per object")
	e.id = id
}

func (e *textEntry) Example() string {
	return e.TheExample
}

func (e *textEntry) Description() string {
	return e.TheDescription
}

func (e *textEntry) AcceptedValues() []string {
	return []string{"Any text"}
}

func (e *textEntry) Group() string {
	return e.group
}

func (e *textEntry) setGroup(groupID string) {
	e.group = groupID
}

func (e *textEntry) rank() int {
	return e.theRank
}

func (e *textEntry) setRank(rank int) {
	e.theRank = rank
}

func (e *textEntry) deepCopy() Entry {
	copy := *e
	return &copy
}

func (e *textEntry) isValid() (err error) {
	if !utils.IsSet(e.TheExample) {
		return fmt.Errorf("Missing example")
	}
	if !utils.IsSet(e.TheDescription) {
		return fmt.Errorf("Missing description")
	}
	return nil
}

func (e *textEntry) validateAndProcessArgs(*args.Store) error {
	// text entries have not much to validate...
	return nil
}

func (e *textEntry) describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v\n", e.id, e.TheDescription)
	} else {
		fmt.Fprintf(&sb, "- %v\n", e.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", e.TheDescription)
		fmt.Fprintf(&sb, "\tType: %v\n", e.Type())
		fmt.Fprintf(&sb, "\tExample: %v\n", genericContentUsageExample(e.id, e.TheExample))
	}

	return sb.String()
}
