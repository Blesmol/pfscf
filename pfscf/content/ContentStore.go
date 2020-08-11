package content

import (
	"fmt"
	"sort"

	"github.com/Blesmol/pfscf/pfscf/preset"
)

// Store stores the list of ContentEntries for a single ChronicleTemplate
type Store map[string]Entry

// NewContentStore creates a new ContentStore object with the provided initial capacity
func NewContentStore(initialCapacity int) (cs Store) {
	return make(Store, initialCapacity)
}

// GetIDs returns the list of IDs for the Presets currently stored in this PresetStore
func (cs Store) GetIDs(includeAliases bool) (idList []string) {
	idList = make([]string, 0, len(cs))
	for id, entry := range cs {
		if includeAliases || id == entry.ID() {
			idList = append(idList, id)
		}
	}
	sort.Strings(idList)
	return idList
}

// Get returns the ContentEntry matching the provided id.
func (cs Store) Get(id string) (ce Entry, exists bool) {
	ce, exists = cs[id]
	return
}

// Set adds or updates the entry with the specified ID in the ContentStore to
// the provided ContentEntry
func (cs *Store) Set(id string, ce Entry) {
	(*cs)[id] = ce
}

// InheritFrom copies over entries from another ContentStore. An error is thrown
// if an entry already exists in both ContentStores.
func (cs *Store) InheritFrom(other Store) (err error) {
	// get content from other object and throw error on duplicates
	for id, otherEntry := range other {
		if _, exists := cs.Get(id); exists {
			return fmt.Errorf("Inheritance error: Content ID '%v' cannot be inherited, because it already exists", id)
		}
		cs.Set(id, otherEntry)
	}

	return nil
}

// Resolve resolves preset requirements for all entries in the ContentStore
func (cs *Store) Resolve(ps preset.Store) (err error) {
	for _, ci := range *cs {
		resolvedCI, err := ci.Resolve(ps)
		if err != nil {
			return err
		}
		cs.Set(ci.ID(), resolvedCI)
	}

	return nil
}
