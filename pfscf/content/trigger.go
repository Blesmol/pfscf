package content

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
)

const (
	typeTrigger = "trigger"
)

type trigger struct {
	Trigger string
	Content ListStore
}

func newTrigger() *trigger {
	var ce trigger
	ce.Content = NewListStore()
	return &ce
}

func (ce *trigger) isValid() (err error) {
	err = checkFieldsAreSet(ce, "Trigger")
	if err != nil {
		return contentValErr(ce, err)
	}
	return ce.Content.IsValid()
}

// resolve the presets for this content object.
func (ce *trigger) resolve(ps preset.Store) (err error) {
	return ce.Content.Resolve(ps)
}

// generateOutput generates the output for this textCell object.
func (ce *trigger) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	err = ce.isValid()
	if err != nil {
		return err
	}

	// will be triggered by any non-nil value
	value := getValue(ce.Trigger, as)
	if value == nil {
		return nil // nothing to do here...
	}

	return ce.Content.GenerateOutput(s, as)
}

// deepCopy creates a deep copy of this entry.
func (ce *trigger) deepCopy() Entry {

	copy := trigger{
		Trigger: ce.Trigger,
		Content: ce.Content.deepCopy(),
	}

	return &copy
}
