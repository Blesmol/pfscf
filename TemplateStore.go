package main

import (
	"fmt"
	"sort"
)

// TemplateStore stores multiple ChronicleTemplates and provides means
// to retrieve them by name.
type TemplateStore struct {
	templates map[string]*ChronicleTemplate // Store as ptrs so that it is easier to modify them do things like aliasing
}

// GetTemplateStore returns a template store that is already filled with all templates
// contained in the main template directory. If some error showed up during reading and
// parsing files, resolving dependencies etc, then nil is returned together with an error.
func GetTemplateStore() (ts *TemplateStore, err error) {
	return getTemplateStoreForDir(GetTemplatesDir())
}

// getTemplateStoreForDir takes a directory and returns a template store
// for all entries in that directory, including its subdirectories
func getTemplateStoreForDir(dirName string) (ts *TemplateStore, err error) {
	yFiles, err := GetTemplateFilesFromDir(dirName)
	if err != nil {
		return nil, err
	}

	ts = new(TemplateStore)
	ts.templates = make(map[string]*ChronicleTemplate)

	for _, yFile := range yFiles {
		ct, err := NewChronicleTemplate(yFile)
		if err != nil {
			return nil, err
		}

		if otherEntry, exists := ts.templates[ct.Name()]; exists {
			return nil, fmt.Errorf("Found multiple templates with ID '%v':\n- %v\n- %v", ct.Name(), otherEntry.Filename(), ct.Filename())
		}
		ts.templates[ct.Name()] = ct
	}

	return ts, nil
}

// GetKeys returns a sorted list of keys contained in this TemplateStore
func (ts *TemplateStore) GetKeys() (keyList []string) {
	keyList = make([]string, 0, len(ts.templates))
	for key := range ts.templates {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	return keyList
}

// GetTemplate returns the template with the specified name from the TemplateStore, or
// an error if no template with that name exists
func (ts *TemplateStore) GetTemplate(templateID string) (ct *ChronicleTemplate, err error) {
	ct, exists := ts.templates[templateID]

	if !exists {
		return nil, fmt.Errorf("Could not find template with ID '%v'", templateID)
	}
	return ct, nil
}
