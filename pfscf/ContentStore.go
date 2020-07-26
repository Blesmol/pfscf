package main

import (
	"fmt"
	"sort"
)

// ContentStore stores the list of ContentEntries for a single ChronicleTemplate
type ContentStore map[string]ContentEntry

// NewContentStore creates a new ContentStore object with the provided initial capacity
func NewContentStore(initialCapacity int) (cs ContentStore) {
	return make(ContentStore, initialCapacity)
}

// GetIDs returns the list of IDs for the Presets currently stored in this PresetStore
func (cs ContentStore) GetIDs(includeAliases bool) (idList []string) {
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
func (cs ContentStore) Get(id string) (pe ContentEntry, exists bool) {
	pe, exists = cs[id]
	return
}

// Set adds or updates the entry with the specified ID in the ContentStore to
// the provided ContentEntry
func (cs *ContentStore) Set(id string, ce ContentEntry) {
	(*cs)[id] = ce
}

// InheritFrom copies over entries from another ContentStore. An error is thrown
// if an entry already exists in both ContentStores.
func (cs *ContentStore) InheritFrom(other ContentStore) (err error) {
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
func (cs *ContentStore) Resolve(ps PresetStore) (err error) {
	for _, ce := range *cs {

		// check that required presets are not contradicting each other
		if err = ps.PresetsAreNotContradicting(ce.Presets()...); err != nil {
			return fmt.Errorf("Error resolving content '%v': %v", ce.ID(), err)
		}

		for _, presetID := range ce.Presets() {
			preset, _ := ps.Get(presetID)
			AddMissingValues(preset, &ce.data, "Presets", "ID")
		}

		cs.Set(ce.ID(), ce)
	}

	return nil
}
