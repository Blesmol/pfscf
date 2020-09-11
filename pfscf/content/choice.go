package content

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
)

const (
	typeChoice = "choice"
)

type choice struct {
	Choice  string
	Content MapStore
}

func newChoice() *choice {
	var ce choice
	ce.Content = NewMapStore()
	return &ce
}

func (ce *choice) isValid() (err error) {
	// TODO arg paramStore to isValid to be able to validate against parameters
	err = checkFieldsAreSet(ce, "Choice")
	if err != nil {
		return contentValErr(ce, err)
	}
	return ce.Content.IsValid()
}

// resolve the presets for this content object.
func (ce *choice) resolve(ps preset.Store) (err error) {
	return ce.Content.Resolve(ps)
}

// generateOutput generates the output for this content object.
func (ce *choice) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	err = ce.isValid()
	if err != nil {
		return err
	}

	selectedChoice := getValue(ce.Choice, as)
	if selectedChoice == nil {
		return nil // nothing to do here...
	}

	// TODO extend to support multiple choices
	for id, entry := range ce.Content {
		if id == *selectedChoice {
			if err = entry.generateOutput(s, as); err != nil {
				return err
			}
		}
	}

	return nil
}

// deepCopy creates a deep copy of this entry.
func (ce *choice) deepCopy() Entry {

	copy := choice{
		Choice:  ce.Choice,
		Content: ce.Content.deepCopy(),
	}

	return &copy
}
