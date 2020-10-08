package content

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/canvas"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeChoice = "choice"
)

type choice struct {
	Choices string
	Content map[string]ListStore
}

func newChoice() *choice {
	var e choice
	e.Content = make(map[string]ListStore)
	return &e
}

func (e *choice) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	// TODO arg paramStore to isValid to be able to validate against parameters
	err = utils.CheckFieldsAreSet(e, "Choices")
	if err != nil {
		return contentValErr(e, err)
	}
	for _, subStore := range e.Content {
		if err = subStore.IsValid(paramStore, canvasStore); err != nil {
			return err
		}
	}
	return nil
}

// resolve the presets for this content object.
func (e *choice) resolve(ps preset.Store) (err error) {
	for _, subStore := range e.Content {
		if err = subStore.Resolve(ps); err != nil {
			return err
		}
	}
	return nil
}

// generateOutput generates the output for this content object.
func (e *choice) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	selectedChoices := e.getChoicesFromArgStore(as)
	if len(selectedChoices) == 0 {
		return nil // nothing to do here...
	}

	for _, choice := range selectedChoices {
		for contentID, contentStore := range e.Content {
			if contentID == choice {
				if err = contentStore.GenerateOutput(s, as); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getChoicesFromArgStore returns a list of values that should be used for the current content.
// Entries are expected to be separated by comma.
func (e *choice) getChoicesFromArgStore(as *args.Store) []string {
	val := getValue(e.Choices, as)

	if val == nil {
		return make([]string, 0)
	}

	return utils.SplitAndTrim(*val, ",")
}

// deepCopy creates a deep copy of this entry.
func (e *choice) deepCopy() Entry {

	copy := choice{
		Choices: e.Choices,
		Content: make(map[string]ListStore),
	}

	for key, entry := range e.Content {
		copy.Content[key] = entry.deepCopy()
	}

	return &copy
}
