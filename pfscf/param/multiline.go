package param

import (
	"fmt"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeMultiline = "multiline"
)

type multilineEntry struct {
	commonFields

	TheExample     string `yaml:"example"`
	TheDescription string `yaml:"description"`
	NumLines       int    `yaml:"lines"`
}

func (e *multilineEntry) Type() string {
	return typeMultiline
}

func (e *multilineEntry) ArgStoreIDs() (result []string) {
	// TODO: add prefixes "_complete", "_single"?
	result = make([]string, 0)
	//result = append(result, e.id)
	for idx := 1; idx <= e.NumLines; idx++ {
		result = append(result, fmt.Sprintf("%v[%v]", e.id, idx))
	}
	return result
}

func (e *multilineEntry) Example() string {
	return e.TheExample
}

func (e *multilineEntry) Description() string {
	return e.TheDescription
}

func (e *multilineEntry) AcceptedValues() []string {
	return []string{"Any text, either as long line (will have auto-break) or split into separate lines"}
}

func (e *multilineEntry) deepCopy() Entry {
	copy := *e

	return &copy
}

func (e *multilineEntry) isValid() (err error) {
	if !utils.IsSet(e.TheExample) {
		return fmt.Errorf("Missing example")
	}
	if !utils.IsSet(e.TheDescription) {
		return fmt.Errorf("Missing description")
	}
	if !utils.IsSet(e.NumLines) {
		return fmt.Errorf("Missing number of lines")
	}
	return nil
}

func (e *multilineEntry) validateAndProcessArgs(as *args.Store) error {
	// TODO check that either single line is provided or split lines
	// TODO check that only valid indices are used

	//argIDs := e.ArgStoreIDs()

	/*
		argValue, exists := as.Get(e.ID())
		utils.Assert(exists, "Existence of entry should have been validated by caller")

		splitArgs := utils.SplitAndTrim(argValue, ",")

		for _, splitArg := range splitArgs {
			if !utils.Contains(e.TheChoices, splitArg) {
				return fmt.Errorf("Invalid choice '%v' was provided. Valid choices are: %v", splitArg, e.TheChoices)
			}
		}
	*/

	return nil
}

func (e *multilineEntry) describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v\n", e.id, e.Description())
	} else {
		fmt.Fprintf(&sb, "- %v\n", e.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", e.Description())
		fmt.Fprintf(&sb, "\tType: %v\n", e.Type())
		fmt.Fprintf(&sb, "\tLines: %v\n", e.NumLines)
		fmt.Fprintf(&sb, "\tExample: %v\n", genericContentUsageExample(e.id, e.Example()))
	}

	return sb.String()
}
