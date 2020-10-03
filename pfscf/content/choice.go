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
	var ce choice
	ce.Content = make(map[string]ListStore)
	return &ce
}

func (ce *choice) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	// TODO arg paramStore to isValid to be able to validate against parameters
	err = utils.CheckFieldsAreSet(ce, "Choices")
	if err != nil {
		return contentValErr(ce, err)
	}
	for _, subStore := range ce.Content {
		if err = subStore.IsValid(paramStore, canvasStore); err != nil {
			return err
		}
	}
	return nil
}

// resolve the presets for this content object.
func (ce *choice) resolve(ps preset.Store) (err error) {
	for _, subStore := range ce.Content {
		if err = subStore.Resolve(ps); err != nil {
			return err
		}
	}
	return nil
}

// generateOutput generates the output for this content object.
func (ce *choice) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	selectedChoices := ce.getChoicesFromArgStore(as)
	if len(selectedChoices) == 0 {
		return nil // nothing to do here...
	}

	for _, choice := range selectedChoices {
		for contentID, contentStore := range ce.Content {
			if contentID == choice {
				if err = contentStore.GenerateOutput(s, as); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getValueList returns a list of values that should be used for the current content.
func (ce *choice) getChoicesFromArgStore(as *args.Store) []string {
	val := getValue(ce.Choices, as)

	if val == nil {
		return make([]string, 0)
	}

	return utils.SplitAndTrim(*val, ",")
}

// deepCopy creates a deep copy of this entry.
func (ce *choice) deepCopy() Entry {

	copy := choice{
		Choices: ce.Choices,
		Content: make(map[string]ListStore),
	}

	for key, entry := range ce.Content {
		copy.Content[key] = entry.deepCopy()
	}

	return &copy
}
