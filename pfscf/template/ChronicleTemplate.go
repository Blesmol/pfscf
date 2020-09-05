package template

import (
	"fmt"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/content"
	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
	"github.com/Blesmol/pfscf/pfscf/yaml"
)

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	id          string
	description string
	inherit     string
	yFilename   string // filename of the originating yaml file
	content     content.Store
	presets     preset.Store
}

// NewChronicleTemplate converts a YamlFile into a ChronicleTemplate. It returns
// an error if the YamlFile cannot be converted to a ChronicleTemplate, e.g. because
// it is missing required entries.
func NewChronicleTemplate(yFilename string, yFile *yaml.File) (ct *ChronicleTemplate, err error) {
	if !utils.IsSet(yFilename) {
		return nil, fmt.Errorf("No filename provided")
	}
	if yFile == nil {
		return nil, fmt.Errorf("Provided YamlFile object is nil")
	}

	if !utils.IsSet(yFile.ID) {
		return nil, fmt.Errorf("Template file '%v' does not contain an ID", yFilename)
	}
	if !utils.IsSet(yFile.Description) {
		return nil, fmt.Errorf("Template file '%v' does not contain a description", yFilename)
	}

	ct = new(ChronicleTemplate)

	ct.id = yFile.ID
	ct.description = yFile.Description
	ct.inherit = yFile.Inherit
	ct.yFilename = yFilename

	ct.content = content.NewContentStore(len(yFile.Content))
	for id, entry := range yFile.Content {
		ct.content[id], err = content.NewContentEntry(id, entry)
		if err != nil {
			return nil, err
		}
	}

	ct.presets = preset.NewStore()
	for id, entry := range yFile.Presets {
		ct.presets.Add(preset.NewEntry(id, entry))
	}

	return ct, nil
}

// ID returns the ID of the chronicle template
func (ct ChronicleTemplate) ID() string {
	return ct.id
}

// Description returns the description of the chronicle template
func (ct ChronicleTemplate) Description() string {
	return ct.description
}

// Inherit returns the ID of the template from which this template inherits
func (ct ChronicleTemplate) Inherit() string {
	return ct.inherit
}

// Filename returns the file name of the chronicle template
func (ct ChronicleTemplate) Filename() string {
	return ct.yFilename
}

// GetContent returns the ContentEntry object matching the provided id
// from the current ChronicleTemplate
func (ct ChronicleTemplate) GetContent(id string) (ci content.Entry, exists bool) {
	return ct.content.Get(id)
}

// GetContentIDs returns a sorted list of content IDs contained in this chronicle template
func (ct ChronicleTemplate) GetContentIDs(includeAliases bool) (idList []string) {
	return ct.content.GetIDs(includeAliases)
}

// Describe describes a single chronicle template. It returns the
// description as a multi-line string
func (ct *ChronicleTemplate) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", ct.ID())
		if utils.IsSet(ct.Description()) {
			fmt.Fprintf(&sb, ": %v", ct.Description())
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", ct.ID())
		fmt.Fprintf(&sb, "\tDesc: %v\n", ct.Description())
		fmt.Fprintf(&sb, "\tFile: %v", ct.Filename())
	}

	return sb.String()
}

// GetExampleArguments returns an array containing all keys and example values for everything that
// has an example value. The result can be passed to the ArgStore.
// TODO test this
func (ct *ChronicleTemplate) GetExampleArguments() (result []string) {
	result = make([]string, 0)
	contentIDs := ct.GetContentIDs(false)

	for _, id := range contentIDs {
		ce, _ := ct.GetContent(id)
		if utils.IsSet(ce.ExampleValue()) {
			result = append(result, fmt.Sprintf("%v=%v", id, ce.ExampleValue()))
		}
	}
	return result
}

// InheritFrom inherits the content and preset entries from another
// ChronicleTemplate object. An error is returned in case a content
// entry exists in both objects. In case a preset object exists in
// both objects, then the one from the original object takes precedence.
func (ct *ChronicleTemplate) InheritFrom(ctOther *ChronicleTemplate) (err error) {
	err = ct.content.InheritFrom(ctOther.content)
	if err != nil {
		return err
	}

	ct.presets.InheritFrom(ctOther.presets)

	return nil
}

// Resolve resolves the presets and content inside this template
func (ct *ChronicleTemplate) Resolve() (err error) {
	if err = ct.presets.Resolve(); err != nil {
		return err
	}
	if err = ct.content.Resolve(ct.presets); err != nil {
		return err
	}
	return nil
}

// WriteToCsvFile creates a CSV file out of the current chronicle template than can be used
// as input for the "batch fill" command
func (ct *ChronicleTemplate) WriteToCsvFile(filename string, separator rune, as *args.Store) (err error) {
	const numPlayers = 7

	records := [][]string{
		{"#ID", ct.ID()},
		{"#Description", ct.Description()},
		{"#"},
		{"#Players"}, // will be filled below with labels
	}
	for idx := 1; idx <= numPlayers; idx++ {
		outerIdx := len(records) - 1
		records[outerIdx] = append(records[outerIdx], fmt.Sprintf("Player %d", idx))
	}

	for _, contentID := range ct.GetContentIDs(false) {
		// entry should be large enough for id column + 7 players
		entry := make([]string, numPlayers+1)

		entry[0] = contentID

		// check if some value was provided on the cmd line that should be filled in everywhere
		if val, exists := as.Get(contentID); exists {
			for colIdx := 1; colIdx <= numPlayers; colIdx++ {
				entry[colIdx] = val
			}
		}

		records = append(records, entry)
	}

	err = csv.WriteFile(filename, separator, records)
	if err != nil {
		return err
	}

	return nil
}

// GenerateOutput adds the content of this chronicle template to the provided stamp.
func (ct *ChronicleTemplate) GenerateOutput(stamp *stamp.Stamp, argStore *args.Store) (err error) {
	for _, key := range ct.GetContentIDs(false) {
		content, _ := ct.GetContent(key)

		err := content.GenerateOutput(stamp, argStore)
		if err != nil {
			return err
		}
	}

	return nil
}
