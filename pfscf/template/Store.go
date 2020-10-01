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

	if err = store.resolve(); err != nil {
		return nil, err
	}

	if err = store.isValid(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) resolve() (err error) {
	// resolve references between templates
	for _, ct := range *s {
		if utils.IsSet(ct.Inherit) {
			// check if parent actually exists, and if so, add chronicle reference
			parentCt, exists := s.Get(ct.Inherit)
			if !exists {
				return fmt.Errorf("Template '%v' inherits from template '%v', but that template cannot be found", ct.ID, ct.Inherit)
			}

			if err = parentCt.addChild(ct); err != nil {
				return err
			}
		}
	}

	// inherit and resolve, starting at root nodes
	if err = s.performPreOrder(func(ct *Chronicle) error {
		for _, childCt := range ct.children {
			// ensure that children inherit before the current chronicle is resolved
			if err = childCt.inheritFrom(ct); err != nil {
				return err
			}
		}

		// resolve each chronicle template
		if err = ct.resolve(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Store) isValid() (err error) {
	// get deterministic template order for validation. The order itself is not relevant
	// for the validation. But if a parent template has invalid entries, then the error
	// message should referr to that template, not to some other template that inherits it.
	// Basically this function guarantees that parent templates are validated before their
	// child templates.

	if err = s.performPreOrder(func(ct *Chronicle) error {
		if err = ct.IsValid(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// performPreOrder traverses the hierarchie structure in a preorder way,
// i.e. first resolves the current node, then child subtrees from left to right.
func (s *Store) performPreOrder(workerFct func(*Chronicle) error) error {
	worklist := s.getRootTemplates()
	for len(worklist) > 0 {
		ct := worklist[0]

		if err := workerFct(ct); err != nil {
			return err
		}

		// add new nodes at beginning for preorder traversal
		worklist = append(ct.children, worklist[1:]...)
	}
	return nil
}

func (s *Store) getRootTemplates() (result []*Chronicle) {
	result = make([]*Chronicle, 0)

	for _, ct := range *s {
		if ct.parent == nil {
			result = append(result, ct)
		}
	}

	sortChronicleList(result)
	return result
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

	s.performPreOrder(func(ct *Chronicle) error {
		if !ct.hasFlag("hidden") {
			// print hierarchie indentation
			for level := ct.getHierarchieLevel(true); level > 0; level-- {
				fmt.Fprint(&sb, "  ")
			}
			fmt.Fprintf(&sb, "- %v: %v\n", ct.ID, ct.Description)
		}
		return nil
	})

	return sb.String()
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

func sortChronicleList(input []*Chronicle) {
	sort.Slice(input, func(i, j int) bool {
		return input[i].ID < input[j].ID
	})
}
