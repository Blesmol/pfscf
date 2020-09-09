package preset

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Store stores a set of preset entries
type Store map[string]*Entry

// NewStore creates a new store.
func NewStore() (s Store) {
	s = make(Store, 0)
	return s
}

// add adds an entry to the store and also sets the ID on the entry
func (s *Store) add(id string, e *Entry) {
	utils.Assert(!utils.IsSet(e.id) || e.id == id, "ID must not be set here")
	if _, exists := (*s)[id]; exists {
		utils.Assert(false, "As we only call this from a map in yaml, no duplicate should occur")
	}

	e.id = id
	(*s)[id] = e
}

// Get returns the Entry matching the provided id.
func (s Store) Get(id string) (e Entry, exists bool) {
	ePtr, exists := s[id]
	if exists {
		e = *ePtr
	}
	return e, exists
}

// UnmarshalYAML unmarshals a Parameter Store
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
			s.add(key, value)
		}
	}

	return nil
}

// InheritFrom copies over entries from another Store that do not yet
// exist in the current Store.
func (s *Store) InheritFrom(other Store) {
	// get presets from other object and intentionally ignore duplicates
	for id, otherEntry := range other {
		if _, exists := s.Get(id); !exists {
			s.add(id, otherEntry.deepCopy())
		}
	}
}

// PresetsAreNotContradicting takes an arbitrary number of preset IDs and
// checks each combination of them on whether they are contradicting or not.
func (s *Store) PresetsAreNotContradicting(IDs ...string) (err error) {
	// ensure that all provided IDs exist. Even before checking the number of arguments
	for _, id := range IDs {
		_, exists := s.Get(id)
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

	firstEntry, _ := s.Get(firstID)

	// check first versus other elements
	for _, otherID := range remainingIDs {
		otherEntry, _ := s.Get(otherID)
		err = firstEntry.doesNotContradict(otherEntry)
		if err != nil {
			return err
		}
	}

	// check for contradictions in remaining elements
	err = s.PresetsAreNotContradicting(remainingIDs...)
	if err != nil {
		return err
	}

	return nil
}

// Resolve resolves inherited values between presets
func (s *Store) Resolve() (err error) {
	resolved := make(map[string]bool)
	for _, entryPtr := range *s {
		if err := s.resolveInternal(entryPtr, &resolved); err != nil {
			return err
		}
	}

	return nil
}

// resolveInternal recursively resolves all presets
func (s *Store) resolveInternal(e *Entry, resolved *map[string]bool, resolveChain ...string) (err error) {
	// check if already resolved
	if _, exists := (*resolved)[e.id]; exists {
		return nil
	}

	// check that we do not have any cyclic dependencies
	for idx, otherID := range resolveChain {
		if e.id == otherID {
			outputChain := append(resolveChain[idx:], otherID) // reduce to relevant part, include conflicting ID again
			return fmt.Errorf("Error resolving preset '%v': Cyclic dependency, chain is %v", e.id, outputChain)
		}
	}

	// ensure that all required presets exist and are already resolved before continuing
	for _, requiredPresetID := range e.presets {
		requiredPreset, exists := s.Get(requiredPresetID)
		if !exists {
			return fmt.Errorf("Error resolving preset '%v': Consumed preset '%v' cannot be found", e.id, requiredPresetID)
		}

		tempResolveChain := append(resolveChain, e.id) // prepare resolveChain for recursive call
		if err = s.resolveInternal(&requiredPreset, resolved, tempResolveChain...); err != nil {
			return err
		}
	}

	// check that required presets are not contradicting each other
	if err = s.PresetsAreNotContradicting(e.presets...); err != nil {
		return fmt.Errorf("Error resolving preset '%v': %v", e.id, err)
	}

	// now finally include values from presets into current entry
	for _, requiredPresetID := range e.presets {
		requiredPreset, _ := s.Get(requiredPresetID)
		e.inheritFrom(requiredPreset)
	}
	(*resolved)[e.id] = true

	return nil
}
