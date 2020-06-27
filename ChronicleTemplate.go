package main

import (
	"fmt"
	"sort"
	"strings"
)

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	id    string
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
	ct.id = yFile.ID
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

// ID returns the ID of the chronicle template
func (ct *ChronicleTemplate) ID() string {
	return ct.id
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
func (ct *ChronicleTemplate) GetContentIDs(includeAliases bool) (idList []string) {
	idList = make([]string, 0, len(ct.yFile.Content))
	for id, entry := range ct.yFile.Content {
		if includeAliases || id == entry.id {
			idList = append(idList, id)
		}
	}
	sort.Strings(idList)
	return idList
}

// Describe describes a single chronicle template. It returns the
// description as a multi-line string
func (ct *ChronicleTemplate) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", ct.ID())
		if IsSet(ct.Description()) {
			fmt.Fprintf(&sb, ": %v", ct.Description())
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", ct.ID())
		fmt.Fprintf(&sb, "\tDesc: %v\n", ct.Description())
		fmt.Fprintf(&sb, "\tFile: %v", ct.Filename())
	}

	return sb.String()
}
