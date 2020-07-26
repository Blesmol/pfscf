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
	presets     PresetStore
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

	ct.presets = NewPresetStore(len(yFile.Presets))
	for id, entry := range yFile.Presets {
		ct.presets.Set(id, NewPresetEntry(id, entry))
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
func (ct *ChronicleTemplate) GetPreset(id string) (pe PresetEntry, exists bool) {
	return ct.presets.Get(id)
}

// GetPresetIDs returns a sorted list of preset IDs contained in this chronicle template.
func (ct *ChronicleTemplate) GetPresetIDs() (idList []string) {
	return ct.presets.GetIDs()
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

// InheritFrom inherits the content and preset entries from another
// ChronicleTemplate object. An error is returned in case a content
// entry exists in both objects. In case a preset object exists in
// both objects, then the one from the original object takes precedence.
func (ct *ChronicleTemplate) InheritFrom(ctOther *ChronicleTemplate) (err error) {
	// get content from other object and throw error on duplicates
	for id, otherEntry := range ctOther.content {
		if _, exists := ct.content[id]; exists {
			return fmt.Errorf("Inheritance error: Content ID '%v' cannot be inherited from '%v', because it already exists in '%v'", id, ctOther.ID(), ct.ID())
		}
		ct.content[id] = otherEntry
	}

	ct.presets.InheritFrom(ctOther.presets)

	return nil
}

// ResolveContent resolves preset requirements for Content
func (ct *ChronicleTemplate) ResolveContent() (err error) {
	for _, ce := range ct.content {

		// check that required presets are not contradicting each other
		if err = ct.presets.PresetsAreNotContradicting(ce.Presets()...); err != nil {
			return fmt.Errorf("Error resolving content '%v.%v': %v", ct.ID(), ce.ID(), err)
		}

		for _, presetID := range ce.Presets() {
			preset, _ := ct.GetPreset(presetID)
			AddMissingValues(preset, &ce.data, "Presets", "ID")
		}

		ct.content[ce.ID()] = ce
	}

	return nil
}

// Resolve resolves the presets and content inside this template
func (ct *ChronicleTemplate) Resolve() (err error) {
	if err = ct.presets.Resolve(); err != nil {
		return err
	}
	if err = ct.ResolveContent(); err != nil {
		return err
	}
	return nil
}
