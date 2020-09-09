package content

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
)

// Store stores a list of parameter descriptions
type Store []Entry

// NewStore creates a new store.
func NewStore() (s Store) {
	s = make(Store, 0)
	return s
}

func (store *Store) add(entry Entry) {
	*store = append(*store, entry)
}

// InheritFrom copies over entries from another Store.
func (store *Store) InheritFrom(other Store) {
	for _, otherEntry := range other {
		store.add(otherEntry.deepCopy())
	}
}

// Resolve resolves preset requirements for all entries in the ContentStore
func (store *Store) Resolve(ps preset.Store) (err error) {
	for _, entry := range *store {
		if err := entry.resolve(ps); err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalYAML unmarshals a Content Store
func (store *Store) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYAML []entryYAML

	sy := make(storeYAML, 0)

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	*store = NewStore()
	for _, ey := range sy {
		store.add(ey.e)
	}

	return nil
}

// GenerateOutput generates the output for the current content store into the provided stamp
func (store *Store) GenerateOutput(stamp *stamp.Stamp, argStore *args.Store) (err error) {
	for _, entry := range *store {
		if err = entry.generateOutput(stamp, argStore); err != nil {
			return err
		}
	}
	return nil
}

// IsValid validates whether all content entries are valid. This means, e.g., that
// the already contain all required values. Thus this should only be called after
// the store was resolved.
func (store *Store) IsValid() (err error) {
	for _, entry := range *store {
		if err = entry.isValid(); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) deepCopy() (copy Store) {
	copy = NewStore()
	for _, entry := range *store {
		copy.add(entry.deepCopy())
	}
	return copy
}
