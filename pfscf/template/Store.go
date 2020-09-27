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
type Store map[string]*Chronicle // Store as ptrs so that it is easier to modify them do things like aliasing

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
func (s *Store) GetTemplateIDs() (keyList []string) {
	keyList = make([]string, 0, len(*s))
	for key := range *s {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	return keyList
}

// Get returns the ChronicleTemplate matching the provided id.
func (s *Store) Get(id string) (ct *Chronicle, exists bool) {
	ct, exists = (*s)[id]
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
func (s *Store) resolveTemplates() (err error) {
	for _, currentID := range s.GetTemplateIDs() {
		ct, _ := s.Get(currentID)
		if err := ct.resolve(); err != nil {
			return err
		}
	}

	return nil
}

// resolveInheritanceBetweenTemplates resolves the inheritance relations between different templates by copying
// over relevant entries, e.g. from the content or presets sections.
func (s *Store) resolveInheritanceBetweenTemplates() (err error) {
	resolvedIDs := make(map[string]bool, 0) // stores IDs of all entries that are already resolved
	for _, currentID := range s.GetTemplateIDs() {
		ct, _ := s.Get(currentID)
		err := s.resolveInheritanceBetweenTemplatesInternal(ct, &resolvedIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) resolveInheritanceBetweenTemplatesInternal(ct *Chronicle, resolvedIDs *map[string]bool, resolveChain ...string) (err error) {
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
	inheritedCt, exists := s.Get(ct.Inherit)
	if !exists {
		return fmt.Errorf("Template '%v' inherits from template '%v', but that template cannot be found", ct.ID, ct.Inherit)
	}

	// add current id to inheritance list and perform recursive call
	resolveChain = append(resolveChain, ct.ID)
	err = s.resolveInheritanceBetweenTemplatesInternal(inheritedCt, resolvedIDs, resolveChain...)
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

func (s *Store) isValid() (err error) {
	// get deterministic template order for validation. The order itself is not relevant
	// for the validation. But if a parent template has invalid entries, then the error
	// message should referr to that template, not to some other template that inherits it.
	// Order is:
	// 1. Parent templates come before their child templates
	// 2. Alphabetical if there is no inheritance relation

	sortedList := newHierarchieStore(s, "").flatten()

	for _, ct := range sortedList {
		if err = ct.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) getTemplatesInheritingFrom(parentID string) (childIDs []string) {
	childIDs = make([]string, 0)

	for key, template := range *s {
		if (!utils.IsSet(parentID) && !utils.IsSet(template.Inherit)) ||
			(template.Inherit == parentID) {
			childIDs = append(childIDs, key)
		}
	}
	sort.Strings(childIDs)

	return childIDs
}

// ListTemplates lists the available templates. Result is returned as multi-line string.
func (s *Store) ListTemplates() (result string) {
	var sb strings.Builder

	completeList := s.listTemplatesInheritingFrom("")
	for _, line := range completeList {
		fmt.Fprintf(&sb, "%v\n", line)
	}

	return sb.String()
}

func (s *Store) listTemplatesInheritingFrom(parentID string) (result []string) {
	result = make([]string, 0)

	for _, childID := range s.getTemplatesInheritingFrom(parentID) {
		template, _ := s.Get(childID)
		result = append(result, fmt.Sprintf("- %v: %v", template.ID, template.Description))

		childrenDesc := s.listTemplatesInheritingFrom(childID)
		for _, childDesc := range childrenDesc {
			result = append(result, fmt.Sprintf("  %v", childDesc))
		}
	}

	return result
}

// SearchForTemplates takes one or multiple keywords and searches for templates
// where all these keywords are included in the description or the id.
// The search is case-insensitive.
// Result is returned as multi-line string.
func (s *Store) SearchForTemplates(keywords ...string) (result string, foundMatch bool) {
	if len(keywords) == 0 {
		return "No keywords provided", false
	}

	// convert all keywords to lower-case
	lowerKW := make([]string, 0)
	for _, kw := range keywords {
		lowerKW = append(lowerKW, strings.ToLower(kw))
	}

	var sb strings.Builder
	foundSomething := false
	for key, template := range *s {
		if termsContainAllKeywords(strings.ToLower(key), strings.ToLower(template.Description), lowerKW...) {
			foundSomething = true
			fmt.Fprintf(&sb, "- %v: %v\n", template.ID, template.Description)
		}
	}

	return sb.String(), foundSomething
}

func termsContainAllKeywords(termA, termB string, keywords ...string) bool {
	for _, kw := range keywords {
		if !strings.Contains(termA, kw) && !strings.Contains(termB, kw) {
			return false
		}
	}

	return true
}
