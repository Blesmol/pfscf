package content

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/canvas"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
)

// ListStore stores a list of parameter descriptions
type ListStore []Entry

// NewListStore creates a new store.
func NewListStore() (s ListStore) {
	s = make(ListStore, 0)
	return s
}

func (s *ListStore) add(entry Entry) {
	*s = append(*s, entry)
}

// InheritFrom copies over entries from another Store.
func (s *ListStore) InheritFrom(other ListStore) {
	for _, otherEntry := range other {
		s.add(otherEntry.deepCopy())
	}
}

// Resolve resolves preset requirements for all entries in the ContentStore
func (s *ListStore) Resolve(ps preset.Store) (err error) {
	for _, entry := range *s {
		if err := entry.resolve(ps); err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalYAML unmarshals a Content List Store
func (s *ListStore) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYAML []entryYAML

	sy := make(storeYAML, 0)

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	*s = NewListStore()
	for _, ey := range sy {
		s.add(ey.e)
	}

	return nil
}

// GenerateOutput generates the output for the current content store into the provided stamp
func (s *ListStore) GenerateOutput(stamp *stamp.Stamp, argStore *args.Store) (err error) {
	for _, entry := range *s {
		if err = entry.generateOutput(stamp, argStore); err != nil {
			return err
		}
	}
	return nil
}

// IsValid validates whether all content entries are valid. This means, e.g., that
// the already contain all required values. Thus this should only be called after
// the store was resolved.
func (s *ListStore) IsValid(paramStore *param.Store, canvasStore *canvas.Store) (err error) {
	for _, entry := range *s {
		if err = entry.isValid(paramStore, canvasStore); err != nil {
			return err
		}
	}
	return nil
}

func (s *ListStore) deepCopy() (copy ListStore) {
	copy = NewListStore()
	for _, entry := range *s {
		copy.add(entry.deepCopy())
	}
	return copy
}
