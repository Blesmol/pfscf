package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// ArgStore holds a mapping between argument keys and values
type ArgStore struct {
	store  map[string]string
	parent *ArgStore
}

// ArgStoreInit can take parameters for initialisation of an
// ArgStore object using NewArgStore()
type ArgStoreInit struct {
	initCapacity int
	parent       *ArgStore
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

// NewArgStore creates a new ArgStore
func NewArgStore(init *ArgStoreInit) *ArgStore {
	as := ArgStore{
		store:  make(map[string]string, init.initCapacity),
		parent: init.parent,
	}
	return &as
}

// HasKey returns whether the ArgStore contains an entry with the given key.
func (as *ArgStore) HasKey(key string) bool {
	_, keyExists := as.store[key]
	if !keyExists && as.HasParent() {
		keyExists = as.parent.HasKey(key)
	}
	return keyExists
}

// Set adds a new value to the ArgStore using the given key.
func (as *ArgStore) Set(key string, value string) {
	as.store[key] = value
}

// Get looks up the given key in the ArgStore and returns the value plus a flag
// indicating whether there is a value for the given key.
func (as ArgStore) Get(key string) (value string, keyExists bool) {
	value, keyExists = as.store[key]
	if keyExists {
		return value, true
	}
	if as.parent != nil {
		return as.parent.Get(key)
	}
	return value, false
}

// HasParent returns whether the ArgStore already has a parent object sei
func (as ArgStore) HasParent() bool {
	return as.parent != nil
}

// SetParent sets the parent ArgStore for the given ArgStore. The returned bool flag
// indicates whether an existing parent was overwritten.
func (as *ArgStore) SetParent(parent *ArgStore) bool {
	hasParent := as.HasParent()
	as.parent = parent
	return hasParent
}

// NumEntries returns the number of entries currently stored in the ArgStore
func (as ArgStore) NumEntries() int {
	return len(as.GetKeys())
}

// GetKeys returns a sorted list of all contained keys
func (as ArgStore) GetKeys() (keyList []string) {
	if as.parent != nil {
		keyList = as.parent.GetKeys()
	} else {
		keyList = make([]string, 0)
	}

	for key := range as.store {
		if !utils.Contains(keyList, key) {
			keyList = append(keyList, key)
		}
	}

	sort.Strings(keyList)

	return keyList
}

// ArgStoreFromArgs takes a list of provided arguments and checks
// that they are in format "<key>=<value>". It will return
// an error on duplicate keys.
func ArgStoreFromArgs(args []string) (as *ArgStore) {
	as = NewArgStore(&ArgStoreInit{initCapacity: len(args)})

	for _, arg := range args {
		splitIdx := strings.Index(arg, "=")
		utils.Assert(splitIdx != -1, "No '=' found in argument")
		utils.Assert(splitIdx != 0, "No key found in argument")
		utils.Assert(splitIdx != (len(arg)-1), "No value found in argument")

		key := arg[:splitIdx]
		value := arg[(splitIdx + 1):]

		if !as.HasKey(key) {
			as.Set(key, value)
		} else {
			panic("Duplicate key found: " + key)
		}
	}

	return as
}

// ArgStoreFromTemplateExamples returns an ArgStore that is filled with
// all the example texts contained in the provided template
func ArgStoreFromTemplateExamples(ct *ChronicleTemplate) (as *ArgStore) {
	contentIDs := ct.GetContentIDs(false)
	as = NewArgStore(&ArgStoreInit{initCapacity: len(contentIDs)})

	for _, id := range contentIDs {
		ce, _ := ct.GetContent(id)
		if utils.IsSet(ce.ExampleValue()) {
			as.Set(id, ce.ExampleValue())
		}
	}

	return as
}

// GetArgStoresFromCsvFile reads a csv file and returns a list of ArgStores that
// contain the required arguments to fill out a chronicle.
func GetArgStoresFromCsvFile(filename string) (argStores []*ArgStore, err error) {
	records, err := csv.ReadCsvFile(filename)
	if err != nil {
		return nil, err
	}

	argStores = make([]*ArgStore, 0)

	if len(records) == 0 {
		return argStores, nil
	}

	numPlayers := len(records[0]) - 1


	for idx := 1; idx <= numPlayers; idx++ {
		as := NewArgStore(&ArgStoreInit{initCapacity: len(records)})

		for _, record := range records {
			key := record[0]
			value := record[idx]
			if as.HasKey(key) {
				return nil, fmt.Errorf("File '%v' contains multiple lines for content ID '%v'", filename, key)
			}

			// only store if there is an actual value
			if utils.IsSet(value) {
				if !utils.IsSet(key) {
					return nil, fmt.Errorf("CSV Line has content value '%v', but is missing content ID in first column", value)
				}
				as.Set(key, value)
			}
		}

		// only add if we have at least one entry here
		if as.NumEntries() >= 1 {
			argStores = append(argStores, as)
		}
	}

	return argStores, nil
}
