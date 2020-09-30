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
	Choice  string
	Content MapStore
}

func newChoice() *choice {
	var ce choice
	ce.Content = NewMapStore()
	return &ce
}

func (ce *choice) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	// TODO arg paramStore to isValid to be able to validate against parameters
	err = utils.CheckFieldsAreSet(ce, "Choice")
	if err != nil {
		return contentValErr(ce, err)
	}
	return ce.Content.IsValid(paramStore, canvasStore)
}

// resolve the presets for this content object.
func (ce *choice) resolve(ps preset.Store) (err error) {
	return ce.Content.Resolve(ps)
}

// generateOutput generates the output for this content object.
func (ce *choice) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	selectedChoices := getValueList(ce.Choice, as)
	if len(selectedChoices) == 0 {
		return nil // nothing to do here...
	}

	for _, choice := range selectedChoices {
		for contentID, entry := range ce.Content {
			if contentID == choice {
				if err = entry.generateOutput(s, as); err != nil {
					return err
				}
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
