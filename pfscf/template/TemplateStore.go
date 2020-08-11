package template

import (
	"fmt"
	"sort"

	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

// Store stores multiple ChronicleTemplates and provides means
// to retrieve them by name.
type Store struct {
	templates map[string]*ChronicleTemplate // Store as ptrs so that it is easier to modify them do things like aliasing
}

// GetStore returns a template store that is already filled with all templates
// contained in the main template directory. If some error showed up during reading and
// parsing files, resolving dependencies etc, then nil is returned together with an error.
func GetStore() (ts *Store, err error) {
	return getStoreForDir(cfg.GetTemplatesDir())
}

// getStoreForDir takes a directory and returns a template store
// for all entries in that directory, including its subdirectories
func getStoreForDir(dirName string) (ts *Store, err error) {
	yFiles, err := yaml.GetTemplateFilesFromDir(dirName)
	if err != nil {
		return nil, err
	}

	ts = new(Store)
	ts.templates = make(map[string]*ChronicleTemplate)

	// put basic yaml files as ChronicleTemplates into the store
	for yFilename, yFile := range yFiles {
		ct, err := NewChronicleTemplate(yFilename, yFile)
		if err != nil {
			return nil, err
		}

		if otherEntry, exists := ts.templates[ct.ID()]; exists {
			return nil, fmt.Errorf("Found multiple templates with ID '%v':\n- %v\n- %v", ct.ID(), otherEntry.Filename(), ct.Filename())
		}
		ts.templates[ct.ID()] = ct
	}

	// resolve inheritance between templates
	resolvedIDs := make(map[string]bool, 0) // stores IDs of all entries that are already resolved
	for _, entry := range ts.templates {
		err := resolveInheritance(ts, entry, &resolvedIDs)
		if err != nil {
			return nil, err
		}
	}

	// resolve presets and content
	for _, templateID := range ts.GetTemplateIDs(false) {
		template, _ := ts.Get(templateID)

		if err = template.Resolve(); err != nil {
			return nil, err
		}
	}

	return ts, nil
}

// resolveInheritance is responsible for resolving template inheritance by copying entries
// from the content and the presets section to other templates.
func resolveInheritance(ts *Store, ct *ChronicleTemplate, resolvedIDs *map[string]bool, resolveChain ...string) (err error) {
	// check if we have already seen that entry
	if _, exists := (*resolvedIDs)[ct.ID()]; exists {
		return nil
	}

	// check if we have a cyclic dependency
	for idx, inheritedID := range resolveChain {
		if inheritedID == ct.ID() {
			resolveChain = append(resolveChain, inheritedID) // add entry before printing to have complete cycle in output
			return fmt.Errorf("Error resolving dependencies of template '%v'. Inheritance chain is %v", ct.ID(), resolveChain[idx:])
		}
	}

	// entries without inheritance information can simply be added to the list of resolved IDs
	if ct.Inherit() == "" {
		(*resolvedIDs)[ct.ID()] = true
		return nil
	}

	inheritedID := ct.Inherit()
	inheritedCe, err := ts.Get(inheritedID)
	if err != nil {
		return err
	}

	// add current id to inheritance list for recursive call
	resolveChain = append(resolveChain, ct.ID())
	err = resolveInheritance(ts, inheritedCe, resolvedIDs, resolveChain...)
	if err != nil {
		return err
	}

	// now resolve chronicle inheritance
	err = ct.InheritFrom(inheritedCe)
	if err != nil {
		return err
	}

	// add to list of resolved entries
	(*resolvedIDs)[ct.ID()] = true
	return nil
}

// GetTemplateIDs returns a sorted list of keys contained in this TemplateStore
func (ts *Store) GetTemplateIDs(includeAliases bool) (keyList []string) {
	keyList = make([]string, 0, len(ts.templates))
	for key, entry := range ts.templates {
		if includeAliases || key == entry.ID() {
			keyList = append(keyList, key)
		}
	}
	sort.Strings(keyList)
	return keyList
}

// Get returns the template with the specified name from the TemplateStore, or
// an error if no template with that name exists
func (ts *Store) Get(templateID string) (ct *ChronicleTemplate, err error) {
	ct, exists := ts.templates[templateID]

	if !exists {
		return nil, fmt.Errorf("Could not find template with ID '%v'", templateID)
	}
	return ct, nil
}

// Get returns the template with the specified name, or
// an error if no template with that name exists. This is merely a
// convenience wrapper to avoid the need to create a TemplateStore
// object just for receiving a single template.
func Get(templateID string) (ct *ChronicleTemplate, err error) {
	ts, err := GetStore()
	if err != nil {
		return nil, err
	}

	return ts.Get(templateID)
}
