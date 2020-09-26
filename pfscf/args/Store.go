package args

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Store holds a mapping between argument keys and values
type Store struct {
	store  map[string]string
	parent *Store
}

// StoreInit can take parameters for initialisation of an
// ArgStore object using NewArgStore()
type StoreInit struct {
	InitCapacity int
	Parent       *Store
	Args         []string
}

const (
	// TODOs:
	// - Restrict input set for key: [a-zA-Z0-9-_]
	// - key should have shortest match, to allow value to contain "="
	argPattern = "(?P<key>.*)=(?P<value>.*)"
)

var (
	argRegex *regexp.Regexp // not used at the moment
)

func init() {
	argRegex = regexp.MustCompile(argPattern)
}

// NewStore creates a new ArgStore
// TODO test this
func NewStore(init StoreInit) (s *Store, err error) {
	var initialCapacity int
	if utils.IsSet(init.InitCapacity) {
		initialCapacity = init.InitCapacity
	} else {
		initialCapacity = len(init.Args)
	}

	localStore := Store{
		store:  make(map[string]string, initialCapacity),
		parent: init.Parent,
	}

	for _, arg := range init.Args {
		key, value, err := splitArgument(arg)
		if err != nil {
			return nil, err
		}

		if !localStore.hasKey(key) {
			localStore.Set(key, value)
		} else {
			return nil, fmt.Errorf("Duplicate key '%v' found", key)
		}
	}

	return &localStore, nil
}

// splitArgument takes an arugment string and tries to split it up in key and value parts.
// TODO test this
func splitArgument(arg string) (key, value string, err error) {
	splitIdx := strings.Index(arg, "=")

	switch splitIdx {
	case -1:
		return "", "", fmt.Errorf("No '=' separator found in argument '%v'", arg)
	case 0:
		return "", "", fmt.Errorf("No key found in argument '%v'", arg)
	case len(arg) - 1:
		return "", "", fmt.Errorf("No value found in argument '%v'", arg)
	}

	key = arg[:splitIdx]
	value = arg[(splitIdx + 1):]

	return key, value, nil
}

// hasKey returns whether the ArgStore contains an entry with the given key.
func (s *Store) hasKey(key string) bool {
	_, keyExists := s.store[key]
	if !keyExists && s.hasParent() {
		keyExists = s.parent.hasKey(key)
	}
	return keyExists
}

// Set adds a new value to the ArgStore using the given key.
func (s *Store) Set(key string, value string) {
	s.store[key] = value
}

// Get looks up the given key in the ArgStore and returns the value plus a flag
// indicating whether there is a value for the given key.
func (s *Store) Get(key string) (value string, keyExists bool) {
	value, keyExists = s.store[key]
	if keyExists {
		return value, true
	}
	if s.parent != nil {
		return s.parent.Get(key)
	}
	return value, false
}

// hasParent returns whether the ArgStore already has a parent object sei
func (s *Store) hasParent() bool {
	return s.parent != nil
}

// SetParent sets the parent ArgStore for the given ArgStore. The returned bool flag
// indicates whether an existing parent was overwritten.
func (s *Store) SetParent(parent *Store) bool {
	hasParent := s.hasParent()
	s.parent = parent
	return hasParent
}

// numEntries returns the number of entries currently stored in the ArgStore
func (s *Store) numEntries() int {
	return len(s.GetKeys())
}

// GetKeys returns a sorted list of all contained keys
func (s *Store) GetKeys() (keyList []string) {
	if s.parent != nil {
		keyList = s.parent.GetKeys()
	} else {
		keyList = make([]string, 0)
	}

	for key := range s.store {
		if !utils.Contains(keyList, key) {
			keyList = append(keyList, key)
		}
	}

	sort.Strings(keyList)

	return keyList
}

// GetArgStoresFromCsvFile reads a csv file and returns a list of ArgStores that
// contain the required arguments to fill out a chronicle.
func GetArgStoresFromCsvFile(filename string) (argStores []*Store, err error) {
	records, err := csv.ReadCsvFile(filename)
	if err != nil {
		return nil, err
	}

	argStores = make([]*Store, 0)

	if len(records) == 0 {
		return argStores, nil
	}

	numPlayers := len(records[0]) - 1

	for idx := 1; idx <= numPlayers; idx++ {
		s, err := NewStore(StoreInit{InitCapacity: len(records)})
		if err != nil {

		}

		for _, record := range records {
			key := record[0]
			value := record[idx]
			if s.hasKey(key) {
				return nil, fmt.Errorf("File '%v' contains multiple lines for content ID '%v'", filename, key)
			}

			// only store if there is an actual value
			if utils.IsSet(value) {
				if !utils.IsSet(key) {
					return nil, fmt.Errorf("CSV Line has content value '%v', but is missing content ID in first column", value)
				}
				s.Set(key, value)
			}
		}

		// only add if we have at least one entry here
		if s.numEntries() >= 1 {
			argStores = append(argStores, s)
		}
	}

	return argStores, nil
}
