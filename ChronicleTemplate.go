package main

import (
	"fmt"
	"sort"
)

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	name  string
	yFile *YamlFile
}

// NewChronicleTemplate converts a YamlFile into a ChronicleTemplate. It returns
// an error if the YamlFile cannot be converted to a ChronicleTemplate, e.g. because
// it is missing required entries.
func NewChronicleTemplate(yFile YamlFile) (ct *ChronicleTemplate, err error) {
	err = checkValidityForChronicleTemplate(yFile)
	if err != nil {
		return nil, err
	}

	ct = new(ChronicleTemplate)
	ct.name = yFile.ID
	ct.yFile = &yFile

	// applying default values
	for key, value := range yFile.Content {
		value.applyDefaults(yFile.Default)
		yFile.Content[key] = value
	}

	return ct, nil
}

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleTemplate
func (ct *ChronicleTemplate) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = ct.yFile.Content[key]
	return
}

// Name returns the name of the chronicle template
func (ct *ChronicleTemplate) Name() string {
	return ct.name
}

// Filename returns the file name of the chronicle template
func (ct *ChronicleTemplate) Filename() string {
	return ct.yFile.GetFilename()
}

// Description returns the description of the chronicle template
func (ct *ChronicleTemplate) Description() string {
	return ct.yFile.Description
}

// checkYamlFileValidity checks
func checkValidityForChronicleTemplate(yFile YamlFile) (err error) {
	if !IsSet(yFile.GetFilename()) {
		return fmt.Errorf("No filename included in YamlFile object")
	}

	if !IsSet(yFile.ID) {
		return fmt.Errorf("Template file '%v' does not contain an ID", yFile.GetFilename())
	}

	if !IsSet(yFile.Description) {
		return fmt.Errorf("Template file '%v' does not contain a description", yFile.GetFilename())
	}

	return nil
}

// GetContentIDs returns a sorted list of content IDs contained in this chronicle template
func (ct *ChronicleTemplate) GetContentIDs() (idList []string) {
	idList = make([]string, 0, len(ct.yFile.Content))
	for id := range ct.yFile.Content {
		idList = append(idList, id)
	}
	sort.Strings(idList)
	return idList
}
