package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// MapStore stores a list of parameter descriptions
//type MapStore map[string]Entry
type MapStore map[string]Entry

// NewMapStore creates a new store.
func NewMapStore() (s MapStore) {
	s = make(MapStore, 0)
	return s
}

func (store *MapStore) add(id string, entry Entry) {
	if _, exists := (*store)[id]; exists {
		utils.Assert(false, "As we only call this from a map in yaml, duplicates should not occur")
	}
	(*store)[id] = entry
}

// InheritFrom inherits entries from another param store. An error is returned in case
// an entry exists in both stores.
func (store *MapStore) InheritFrom(other MapStore) (err error) {
	for otherID, otherEntry := range other {
		if _, exists := (*store)[otherID]; exists {
			return fmt.Errorf("Duplicate parameter ID '%v' found while inheriting", otherID)
		}
		store.add(otherID, otherEntry.deepCopy())
	}

	return nil
}

// Resolve resolves preset requirements for all entries in the ContentStore
func (store *MapStore) Resolve(ps preset.Store) (err error) {
	for _, entry := range *store {
		if err := entry.resolve(ps); err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalYAML unmarshals a Content Store
func (store *MapStore) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYAML map[string]entryYAML

	sy := make(storeYAML, 0)

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	*store = NewMapStore()
	for key, ey := range sy {
		store.add(key, ey.e)
	}

	return nil
}

// IsValid validates whether all content entries are valid. This means, e.g., that
// the already contain all required values. Thus this should only be called after
// the store was resolved.
func (store *MapStore) IsValid() (err error) {
	for _, entry := range *store {
		if err = entry.isValid(); err != nil {
			return err
		}
	}
	return nil
}

func (store *MapStore) deepCopy() (copy MapStore) {
	copy = NewMapStore()
	for key, entry := range *store {
		copy.add(key, entry.deepCopy())
	}
	return copy
}
