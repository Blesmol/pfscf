package main

import (
	"fmt"
	"sort"
	"strings"
)

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	id          string
	description string
	inherit     string
	yFilename   string // filename of the originating yaml file
	content     map[string]ContentEntry
	presets     map[string]ContentEntry
}

// NewChronicleTemplate converts a YamlFile into a ChronicleTemplate. It returns
// an error if the YamlFile cannot be converted to a ChronicleTemplate, e.g. because
// it is missing required entries.
func NewChronicleTemplate(yFilename string, yFile *YamlFile) (ct *ChronicleTemplate, err error) {
	if !IsSet(yFilename) {
		return nil, fmt.Errorf("No filename provided")
	}
	if yFile == nil {
		return nil, fmt.Errorf("Provided YamlFile object is nil")
	}

	if !IsSet(yFile.ID) {
		return nil, fmt.Errorf("Template file '%v' does not contain an ID", yFilename)
	}
	if !IsSet(yFile.Description) {
		return nil, fmt.Errorf("Template file '%v' does not contain a description", yFilename)
	}

	ct = new(ChronicleTemplate)

	ct.id = yFile.ID
	ct.description = yFile.Description
	ct.inherit = yFile.Inherit
	ct.yFilename = yFilename

	ct.content = make(map[string]ContentEntry, len(yFile.Content))
	for id, entry := range yFile.Content {
		ct.content[id] = NewContentEntry(id, entry)
	}

	ct.presets = make(map[string]ContentEntry, len(yFile.Presets))
	for id, entry := range yFile.Presets {
		ct.presets[id] = NewContentEntry(id, entry)
	}

	// TODO temporary workaround until presets are working properly
	// TODO remove
	if defPreset, exists := ct.presets["default"]; exists {
		for id, entry := range ct.content {
			entry.AddMissingValuesFrom(&defPreset)
			ct.content[id] = entry
		}
	}

	return ct, nil
}

// ID returns the ID of the chronicle template
func (ct *ChronicleTemplate) ID() string {
	return ct.id
}

// Description returns the description of the chronicle template
func (ct *ChronicleTemplate) Description() string {
	return ct.description
}

// Inherit returns the ID of the template from which this template inherits
func (ct *ChronicleTemplate) Inherit() string {
	return ct.inherit
}

// Filename returns the file name of the chronicle template
func (ct *ChronicleTemplate) Filename() string {
	return ct.yFilename
}

// GetPreset returns the preset ContentEntry matching the provided id from
// the current ChronicleTemplate
func (ct *ChronicleTemplate) GetPreset(id string) (ce ContentEntry, exists bool) {
	ce, exists = ct.presets[id]
	return
}

// GetContent returns the ContentEntry object matching the provided id
// from the current ChronicleTemplate
func (ct *ChronicleTemplate) GetContent(id string) (ce ContentEntry, exists bool) {
	ce, exists = ct.content[id]
	return
}

// GetContentIDs returns a sorted list of content IDs contained in this chronicle template
func (ct *ChronicleTemplate) GetContentIDs(includeAliases bool) (idList []string) {
	idList = make([]string, 0, len(ct.content))
	for id, entry := range ct.content {
		if includeAliases || id == entry.ID() {
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
