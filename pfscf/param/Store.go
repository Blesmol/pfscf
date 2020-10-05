package param

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Store stores a list of parameter descriptions
type Store map[string]Entry

// NewStore creates a new store.
func NewStore() (s Store) {
	s = make(Store, 0)
	return s
}

// add adds an entry to the store and also sets the ID on the entry
func (s *Store) add(id string, e Entry) (err error) {
	utils.Assert(!utils.IsSet(e.ID()) || id == e.ID(), "ID must not be set here")
	if _, exists := (*s)[id]; exists {
		return fmt.Errorf("Found multiple parameter definitions with id '%v'", id)
	}

	if !utils.IsSet(e.ID()) {
		e.setID(id)
	}
	(*s)[id] = e
	return nil
}

// Get returns the Entry matching the provided id.
func (s *Store) Get(id string) (e Entry, exists bool) {
	e, exists = (*s)[id]
	return
}

// UnmarshalYAML unmarshals a Parameter Store
func (s *Store) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type storeYAML map[string]map[string]entryYAML

	sy := make(storeYAML, 0)

	err = unmarshal(&sy)
	if err != nil {
		return err
	}

	*s = NewStore()
	for groupID, group := range sy {
		for entryID, entry := range group {
			entry.e.setGroup(groupID)
			if err = s.add(entryID, entry.e); err != nil {
				return err
			}
		}
	}

	return nil
}

// InheritFrom inherits entries from another param store. An error is returned in case
// an entry exists in both stores.
func (s *Store) InheritFrom(other *Store) (err error) {
	for otherID, otherEntry := range *other {
		if err = s.add(otherID, otherEntry.deepCopy()); err != nil {
			return fmt.Errorf("Error while inheriting parent template: %v", err)
		}
	}

	return nil
}

// IsValid checks whether all entries are valid.
func (s *Store) IsValid() (err error) {
	for _, entry := range *s {
		if err = entry.isValid(); err != nil {
			return fmt.Errorf("Error while validating parameter definition '%v': %v", entry.ID(), err)
		}
	}
	return nil
}

func (s *Store) getArgNameToEntryMapping() (result map[string]Entry) {
	result = make(map[string]Entry)

	for _, paramEntry := range *s {
		for _, argName := range paramEntry.ArgStoreIDs() {
			result[argName] = paramEntry
		}
	}

	return result
}

// ValidateAndProcessArgs checks whether all arguments in the arg store have a
// corresponding parameter entry.
func (s *Store) ValidateAndProcessArgs(as *args.Store) (err error) {
	argNameToEntry := s.getArgNameToEntryMapping()

	for _, argName := range as.GetKeys() {
		paramEntry, pExists := argNameToEntry[argName]

		// check that all entries in the arg store have a corresponding parameter entry
		if !pExists {
			return fmt.Errorf("Error while validating argument '%v': No corresponding parameter registered for template", argName)
		}

		// ask each type whether the provided argument is valid, and add entries to argStore if required
		if err = paramEntry.validateAndProcessArgs(as); err != nil {
			return fmt.Errorf("Error while validating argument '%v': %v", argName, err)
		}
	}

	return nil
}

// GetExampleArguments returns an array containing all keys and example values for all parameters.
// The result can be passed to the ArgStore.
func (s *Store) GetExampleArguments() (result []string) {
	result = make([]string, 0)

	for _, entry := range *s {
		for _, argStoreID := range entry.ArgStoreIDs() {
			result = append(result, fmt.Sprintf("%v=%v", argStoreID, entry.Example()))
		}
	}

	return result
}

// GetKeysSortedByName returns the list of keys contained in this store as list sorted by rank.
func (s *Store) GetKeysSortedByName() (result []string) {
	result = make([]string, 0)

	for key := range *s {
		result = append(result, key)
	}

	sort.Strings(result)

	return result
}

// getKeysSortedByRank returns the list of keys contained in this store as list sorted by rank.
func (s *Store) getKeysSortedByRank() (result []string) {
	type pair struct {
		id   string
		rank int
	}
	sorting := make([]pair, 0)

	for key, entry := range *s {
		sorting = append(sorting, pair{key, entry.rank()})
	}

	sort.Slice(sorting, func(i, j int) bool {
		return sorting[i].rank < sorting[j].rank
	})

	result = make([]string, 0)
	for _, sortingEntry := range sorting {
		result = append(result, sortingEntry.id)
	}

	return result
}

// GetKeysForGroupSortedByRank returns the list of key IDs contained in the listed group.
func (s *Store) GetKeysForGroupSortedByRank(group string) (result []string) {
	result = make([]string, 0)

	sortedKeys := s.getKeysSortedByRank()

	for _, key := range sortedKeys {
		entry := (*s)[key]
		if group == entry.Group() {
			result = append(result, entry.ID())
		}
	}

	return result
}

// GetGroupsSortedByRank returns the list of groups sorted by rank
func (s *Store) GetGroupsSortedByRank() (result []string) {
	result = make([]string, 0)

	sortedKeys := s.getKeysSortedByRank()

	for _, key := range sortedKeys {
		entry := (*s)[key]
		if !utils.Contains(result, entry.Group()) {
			result = append(result, entry.Group())
		}
	}

	return result
}

// Describe returns a short textual description of all parameters contained in this store.
// It returns the description as a multi-line string.
func (s *Store) Describe(verbose bool) (result string) {
	var sb strings.Builder

	for _, groupName := range s.GetGroupsSortedByRank() {
		fmt.Fprintf(&sb, "%v:\n", groupName)

		for _, key := range s.GetKeysForGroupSortedByRank(groupName) {
			entry, _ := s.Get(key)
			fmt.Fprintf(&sb, entry.describe(verbose))
		}

		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}
