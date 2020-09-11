package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

// Store stores multiple ChronicleTemplates and provides means
// to retrieve them by name.
type Store map[string]*ChronicleTemplate // Store as ptrs so that it is easier to modify them do things like aliasing

// newStore creates a new Store object
func newStore() (store *Store) {
	s := make(Store, 0)
	return &s
}

// GetStore returns a template store that is already filled with all templates
// contained in the main template directory. If some error showed up during reading and
// parsing files, resolving dependencies etc, then nil is returned together with an error.
func GetStore() (ts *Store, err error) {
	return getStoreForDir(cfg.GetTemplatesDir())
}

// GetTemplateIDs returns a sorted list of keys contained in this Store
func (store *Store) GetTemplateIDs() (keyList []string) {
	keyList = make([]string, 0, len(*store))
	for key := range *store {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	return keyList
}

// Get returns the ChronicleTemplate matching the provided id.
func (store *Store) Get(id string) (ct *ChronicleTemplate, exists bool) {
	ct, exists = (*store)[id]
	return
}

// getStoreForDir takes a directory and returns a template store
// for all entries in that directory, including its subdirectories
func getStoreForDir(dir string) (store *Store, err error) {
	filenames, err := yaml.GetYamlFilenamesFromDir(dir)
	if err != nil {
		return nil, err
	}

	store = newStore()

	// read all templates from files and put into store
	for _, filename := range filenames {
		ct := NewChronicleTemplate(filename)
		err = yaml.ReadYamlFile(filename, &ct)
		if err != nil {
			return nil, err
		}
		ct.ensureStoresAreInitialized() // workaround for bug / shitty behavior in go-yaml

		// check for duplicate IDs
		if other, exists := store.Get(ct.ID); exists {
			return nil, fmt.Errorf("Found multiple templates with ID '%v':\n- %v\n- %v", ct.ID, ct.filename, other.filename)
		}

		(*store)[ct.ID] = &ct
	}

	if err = store.resolveInheritanceBetweenTemplates(); err != nil {
		return nil, err
	}

	if err = store.resolveTemplates(); err != nil {
		return nil, err
	}

	if err = store.isValid(); err != nil {
		return nil, err
	}

	return store, nil
}

// resolveTemplates resolves the relations inside each template contained in this store
func (store *Store) resolveTemplates() (err error) {
	for _, currentID := range store.GetTemplateIDs() {
		ct, _ := store.Get(currentID)
		if err := ct.resolve(); err != nil {
			return err
		}
	}

	return nil
}

// resolveInheritanceBetweenTemplates resolves the inheritance relations between different templates by copying
// over relevant entries, e.g. from the content or presets sections.
func (store *Store) resolveInheritanceBetweenTemplates() (err error) {
	resolvedIDs := make(map[string]bool, 0) // stores IDs of all entries that are already resolved
	for _, currentID := range store.GetTemplateIDs() {
		ct, _ := store.Get(currentID)
		err := store.resolveInheritanceBetweenTemplatesInternal(ct, &resolvedIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) resolveInheritanceBetweenTemplatesInternal(ct *ChronicleTemplate, resolvedIDs *map[string]bool, resolveChain ...string) (err error) {
	// check if we have already seen that entry
	if _, exists := (*resolvedIDs)[ct.ID]; exists {
		return nil
	}

	// check if we have a cyclic dependency
	for idx, inheritedID := range resolveChain {
		if inheritedID == ct.ID {
			resolveChain = append(resolveChain, ct.ID) // add entry before printing to have complete cycle in output
			return fmt.Errorf("Error resolving dependencies of template '%v'. Inheritance chain is %v", ct.ID, resolveChain[idx:])
		}
	}

	// entries without inheritance information can simply be added to the list of resolved IDs
	if ct.Inherit == "" {
		(*resolvedIDs)[ct.ID] = true
		return nil
	}

	// check if inherited ID exists and retrieve entry
	inheritedCt, exists := store.Get(ct.Inherit)
	if !exists {
		return fmt.Errorf("Template '%v' inherits from template '%v', but that template cannot be found", ct.ID, ct.Inherit)
	}

	// add current id to inheritance list and perform recursive call
	resolveChain = append(resolveChain, ct.ID)
	err = store.resolveInheritanceBetweenTemplatesInternal(inheritedCt, resolvedIDs, resolveChain...)
	if err != nil {
		return err
	}

	// now resolve chronicle inheritance
	err = ct.inheritFrom(inheritedCt)
	if err != nil {
		return err
	}

	// add to list of resolved entries
	(*resolvedIDs)[ct.ID] = true

	return nil
}

func (store *Store) isValid() (err error) {
	for _, entry := range *store {
		if err = entry.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) getTemplatesInheritingFrom(parentID string) (childIDs []string) {
	childIDs = make([]string, 0)

	for key, template := range *store {
		if (!utils.IsSet(parentID) && !utils.IsSet(template.Inherit)) ||
			(template.Inherit == parentID) {
			childIDs = append(childIDs, key)
		}
	}
	sort.Strings(childIDs)

	return childIDs
}

// ListTemplates lists the available templates. Result is returned as multi-line string.
func (store *Store) ListTemplates() (result string) {
	var sb strings.Builder

	completeList := store.listTemplatesInheritingFrom("")
	for _, line := range completeList {
		fmt.Fprintf(&sb, "%v\n", line)
	}

	return sb.String()
}

func (store *Store) listTemplatesInheritingFrom(parentID string) (result []string) {
	result = make([]string, 0)

	for _, childID := range store.getTemplatesInheritingFrom(parentID) {
		template, _ := store.Get(childID)
		result = append(result, fmt.Sprintf("- %v: %v", template.ID, template.Description))

		childrenDesc := store.listTemplatesInheritingFrom(childID)
		for _, childDesc := range childrenDesc {
			result = append(result, fmt.Sprintf("  %v", childDesc))
		}
	}

	return result
}
