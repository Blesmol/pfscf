package canvas

import (
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Store stores a set of preset entries
type Store map[string]*Entry

// NewStore creates a new store.
func NewStore() (s Store) {
	s = make(Store, 0)
	return s
}

// Add adds an entry to the store and also sets the ID on the entry
func (s *Store) Add(id string, e *Entry) {
	utils.Assert(!utils.IsSet(e.id) || e.id == id, "ID must not be set here")
	if _, exists := (*s)[id]; exists {
		utils.Assert(false, "As we only call this from a map in yaml, no duplicate should occur")
	}

	e.id = id
	(*s)[id] = e
}

// Get returns the Entry matching the provided id.
func (s Store) Get(id string) (e *Entry, exists bool) {
	e, exists = s[id]
	return
}

// UnmarshalYAML unmarshals a canvas Store
func (s *Store) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYaml Store
	sy := storeYaml(NewStore()) // avoid unmarshalling recursion

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	// create return value and copy over content from temporary object
	*s = NewStore()
	for key, value := range sy {
		if value != nil {
			s.Add(key, value)
		}
	}

	return nil
}

// InheritFrom copies over entries from another Store that do not yet
// exist in the current Store.
func (s *Store) InheritFrom(other Store) {
	// get entries from other object
	for id, otherEntry := range other {

		if localEntry, exists := s.Get(id); !exists {
			// if no local entry with the name exists, take from parent store
			s.Add(id, otherEntry.deepCopy())
		} else {
			// take entries from parent where appropriate
			localEntry.inherit(otherEntry)
		}
	}
}

// IsValid validates whether all contained entries are valid. This should only be called after
// the store was resolved.
func (s *Store) IsValid() (err error) {
	for _, entryPtr := range *s {
		if err = entryPtr.isValid(); err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves inherited values between presets
func (s *Store) Resolve() (err error) {
	for _, entryPtr := range *s {
		if err := entryPtr.resolve(s); err != nil {
			return err
		}
	}

	return nil
}

// AddCanvasesToStamp adds all included canvases to the provided stamp.
func (s *Store) AddCanvasesToStamp(stamp *stamp.Stamp) {
	for _, entry := range *s {
		stamp.AddCanvas(entry.id, *entry.X, *entry.Y, *entry.X2, *entry.Y2)
	}
}
