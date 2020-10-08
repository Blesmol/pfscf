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
	typeTrigger = "trigger"
)

type trigger struct {
	Trigger string
	Content ListStore
}

func newTrigger() *trigger {
	var e trigger
	e.Content = NewListStore()
	return &e
}

func (e *trigger) isValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	err = utils.CheckFieldsAreSet(e, "Trigger")
	if err != nil {
		return contentValErr(e, err)
	}
	return e.Content.IsValid(paramStore, canvasStore)
}

// resolve the presets for this content object.
func (e *trigger) resolve(ps preset.Store) (err error) {
	return e.Content.Resolve(ps)
}

// generateOutput generates the output for this object.
func (e *trigger) generateOutput(s *stamp.Stamp, as *args.Store) (err error) {
	// will be triggered by any non-nil value
	value := getValue(e.Trigger, as)
	if value == nil {
		return nil // nothing to do here...
	}

	return e.Content.GenerateOutput(s, as)
}

// deepCopy creates a deep copy of this entry.
func (e *trigger) deepCopy() Entry {
	copy := trigger{
		Trigger: e.Trigger,
		Content: e.Content.deepCopy(),
	}

	return &copy
}
