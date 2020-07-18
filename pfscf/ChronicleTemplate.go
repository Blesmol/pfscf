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

// GetPresetIDs returns a sorted list of preset IDs contained in this chronicle template.
func (ct *ChronicleTemplate) GetPresetIDs() (idList []string) {
	idList = make([]string, 0, len(ct.presets))
	for id := range ct.presets {
		idList = append(idList, id)
	}
	sort.Strings(idList)
	return idList
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

	// get presets from other object and intentionally ignore duplicates
	for id, otherEntry := range ctOther.presets {
		if _, exists := ct.presets[id]; !exists {
			ct.presets[id] = otherEntry
		}
	}

	return nil
}

// presetsAreNotContradicting takes an arbitrary number of preset IDs and
// checks each combination of them on whether they are contradicting or not.
func (ct *ChronicleTemplate) presetsAreNotContradicting(IDs ...string) (err error) {
	// ensure that all provided IDs exist. Even before checking the number of arguments
	for _, id := range IDs {
		_, exists := ct.GetPreset(id)
		if !exists {
			return fmt.Errorf("Preset '%v' does not exist", id)
		}
	}

	// with 0 or 1 entries, no contradictions are possible
	if len(IDs) <= 1 {
		return nil
	}

	firstID := IDs[0]
	remainingIDs := IDs[1:]

	firstEntry, _ := ct.GetPreset(firstID)

	// check first versus other elements
	for _, otherID := range remainingIDs {
		otherEntry, _ := ct.GetPreset(otherID)
		err = firstEntry.IsNotContradictingWith(otherEntry)
		if err != nil {
			return err
		}
	}

	// check for contradictions in remaining elements
	err = ct.presetsAreNotContradicting(remainingIDs...)
	if err != nil {
		return err
	}

	return nil
}

// ResolveContent resolves preset requirements for Content
func (ct *ChronicleTemplate) ResolveContent() (err error) {
	for _, ce := range ct.content {

		// check that required presets are not contradicting each other
		if err = ct.presetsAreNotContradicting(ce.Presets()...); err != nil {
			return fmt.Errorf("Error resolving content '%v.%v': %v", ct.ID(), ce.ID(), err)
		}

		for _, presetID := range ce.Presets() {
			preset, _ := ct.GetPreset(presetID)
			ce.AddMissingValuesFrom(&preset)
		}

		ct.content[ce.ID()] = ce
	}

	return nil
}

// ResolvePresets resolves inherited values between presets
func (ct *ChronicleTemplate) ResolvePresets() (err error) {
	resolved := make(map[string]bool)
	for _, ce := range ct.presets {
		if err := ct.resolvePresetsInternal(ce, &resolved); err != nil {
			return err
		}
	}
	return nil
}

// resolvePresetsInternal recursively resolves all presets
func (ct *ChronicleTemplate) resolvePresetsInternal(ce ContentEntry, resolved *map[string]bool, resolveChain ...string) (err error) {
	// check if already resolved
	if _, exists := (*resolved)[ce.ID()]; exists {
		return nil
	}

	// check that we do not have any cyclic dependencies
	for idx, otherID := range resolveChain {
		if ce.ID() == otherID {
			outputChain := append(resolveChain[idx:], otherID) // reduce to relevant part, include conflicting ID again
			return fmt.Errorf("Error resolving preset '%v.%v': Cyclic dependency, chain is %v", ct.ID(), ce.ID(), outputChain)
		}
	}

	// ensure that all required presets exist and are already resolved before continuing
	for _, requiredPresetID := range ce.Presets() {
		requiredPreset, exists := ct.GetPreset(requiredPresetID)
		if !exists {
			return fmt.Errorf("Error resolving preset '%v.%v': Consumed preset '%v' cannot be found", ct.ID(), ce.ID(), requiredPresetID)
		}

		tempResolveChain := append(resolveChain, ce.ID()) // prepare resolveChain for recursive call
		if err = ct.resolvePresetsInternal(requiredPreset, resolved, tempResolveChain...); err != nil {
			return err
		}
	}

	// check that required presets are not contradicting each other
	if err = ct.presetsAreNotContradicting(ce.Presets()...); err != nil {
		return fmt.Errorf("Error resolving preset '%v.%v': %v", ct.ID(), ce.ID(), err)
	}

	// now finally include values from presets into current entry
	for _, requiredPresetID := range ce.Presets() {
		requiredPreset, _ := ct.GetPreset(requiredPresetID)
		ce.AddMissingValuesFrom(&requiredPreset)
	}

	// update entry stored in ChronicleTemplate, record that we are ready, and thats it.
	ct.presets[ce.ID()] = ce
	(*resolved)[ce.ID()] = true

	return nil
}
