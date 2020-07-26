package main

import (
	"fmt"
	"sort"
)

// PresetStore stores the list of presets for a single ChronicleTemplate
type PresetStore map[string]PresetEntry

// NewPresetStore creates a new PresetStore object with the provided initial capacity
func NewPresetStore(initialCapacity int) (ps PresetStore) {
	return make(PresetStore, initialCapacity)
}

// GetIDs returns the list of IDs for the Presets currently stored in this PresetStore
func (ps PresetStore) GetIDs() (idList []string) {
	idList = make([]string, 0, len(ps))
	for id := range ps {
		idList = append(idList, id)
	}
	sort.Strings(idList)
	return idList
}

// Get returns the PresetEntry matching the provided id.
func (ps PresetStore) Get(id string) (pe PresetEntry, exists bool) {
	pe, exists = ps[id]
	return
}

// Set adds or updates the entry with the specified ID in the PresetStore to
// the provided PresetEntry
func (ps *PresetStore) Set(id string, pe PresetEntry) {
	(*ps)[id] = pe
}

// InheritFrom copies over entries from another PresetStore that do not yet
// exist in the current PresetStore
func (ps *PresetStore) InheritFrom(other PresetStore) {
	// get presets from other object and intentionally ignore duplicates
	for id, otherEntry := range other {
		if _, exists := ps.Get(id); !exists {
			ps.Set(id, otherEntry)
		}
	}
}

// PresetsAreNotContradicting takes an arbitrary number of preset IDs and
// checks each combination of them on whether they are contradicting or not.
func (ps PresetStore) PresetsAreNotContradicting(IDs ...string) (err error) {
	// ensure that all provided IDs exist. Even before checking the number of arguments
	for _, id := range IDs {
		_, exists := ps.Get(id)
		if !exists {
			return fmt.Errorf("Preset '%v' does not exist", id)
		}
	}

	// with 0 or 1 entries, no contradictions are possible
	if len(IDs) <= 1 {
		return nil
	}

	firstID := IDs[0]
	remainingIDs := IDs[1:]

	firstEntry, _ := ps.Get(firstID)

	// check first versus other elements
	for _, otherID := range remainingIDs {
		otherEntry, _ := ps.Get(otherID)
		err = firstEntry.IsNotContradictingWith(otherEntry)
		if err != nil {
			return err
		}
	}

	// check for contradictions in remaining elements
	err = ps.PresetsAreNotContradicting(remainingIDs...)
	if err != nil {
		return err
	}

	return nil
}

// Resolve resolves inherited values between presets
func (ps *PresetStore) Resolve() (err error) {
	resolved := make(map[string]bool)
	for _, pe := range *ps {
		if err := ps.resolveInternal(pe, &resolved); err != nil {
			return err
		}
	}

	return nil
}

// resolveInternal recursively resolves all presets
func (ps *PresetStore) resolveInternal(pe PresetEntry, resolved *map[string]bool, resolveChain ...string) (err error) {
	// check if already resolved
	if _, exists := (*resolved)[pe.ID]; exists {
		return nil
	}

	// check that we do not have any cyclic dependencies
	for idx, otherID := range resolveChain {
		if pe.ID == otherID {
			outputChain := append(resolveChain[idx:], otherID) // reduce to relevant part, include conflicting ID again
			return fmt.Errorf("Error resolving preset '%v': Cyclic dependency, chain is %v", pe.ID, outputChain)
		}
	}

	// ensure that all required presets exist and are already resolved before continuing
	for _, requiredPresetID := range pe.Presets {
		requiredPreset, exists := ps.Get(requiredPresetID)
		if !exists {
			return fmt.Errorf("Error resolving preset '%v': Consumed preset '%v' cannot be found", pe.ID, requiredPresetID)
		}

		tempResolveChain := append(resolveChain, pe.ID) // prepare resolveChain for recursive call
		if err = ps.resolveInternal(requiredPreset, resolved, tempResolveChain...); err != nil {
			return err
		}
	}

	// check that required presets are not contradicting each other
	if err = ps.PresetsAreNotContradicting(pe.Presets...); err != nil {
		return fmt.Errorf("Error resolving preset '%v': %v", pe.ID, err)
	}

	// now finally include values from presets into current entry
	for _, requiredPresetID := range pe.Presets {
		requiredPreset, _ := ps.Get(requiredPresetID)
		AddMissingValues(requiredPreset, &pe, "Presets", "ID")
	}

	// update entry stored in ChronicleTemplate, record that we are ready, and thats it.
	ps.Set(pe.ID, pe)
	(*resolved)[pe.ID] = true

	return nil
}
