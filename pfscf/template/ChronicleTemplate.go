package template

import (
	"fmt"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/content"
	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// ChronicleTemplate is the new approach for the Chronicle Template
type ChronicleTemplate struct {
	ID          string
	Description string
	Inherit     string
	Parameters  param.Store
	Presets     preset.Store
	Content     content.Store

	filename string // filename of the originating yaml file
}

// NewChronicleTemplate returns a new ChronicleTemplate object.
func NewChronicleTemplate(filename string) (ct ChronicleTemplate) {
	ct.filename = filename
	return ct
}

func templateErr(ct *ChronicleTemplate, errIn error) (errOut error) {
	return fmt.Errorf("Template %v: %v", ct.ID, errIn)
}

func templateErrf(ct *ChronicleTemplate, msg string, args ...interface{}) (errOut error) {
	return fmt.Errorf("Template '%v': "+msg, ct.ID, args)
}

// ensureStoresAreInitialized is a workaround for the behavior of the stupid f... yaml library.
// If a section like "parameters:" is present, but empty, it will not be unmarshalled, and the
// underlying data structure will be ZEROed. So the stores will be uninitialized. Even if they
// were initialized before the unmarshalling. Yeah, great.
// See https://github.com/go-yaml/yaml/issues/395 , might be fixed with go-yaml v3 in the future.
func (ct *ChronicleTemplate) ensureStoresAreInitialized() {
	if ct.Parameters == nil {
		ct.Parameters = param.NewStore()
	}
	if ct.Presets == nil {
		ct.Presets = preset.NewStore()
	}
	if ct.Content == nil {
		ct.Content = content.NewStore()
	}
}

// GetExampleArguments returns an array containing all keys and example values for all parameters.
// The result can be passed to the ArgStore.
func (ct *ChronicleTemplate) GetExampleArguments() (result []string) {
	return ct.Parameters.GetExampleArguments()
}

// inheritFrom inherits entries from multiple sections from another
// ChronicleTemplate object. An error is returned in case a content
// entry from sections 'parameters' or 'content' exists in both objects.
// In case a preset entry exists in both objects, then the one from the original
// object takes precedence.
func (ct *ChronicleTemplate) inheritFrom(otherCT *ChronicleTemplate) (err error) {
	err = ct.Parameters.InheritFrom(&otherCT.Parameters)
	if err != nil {
		return err
	}

	ct.Presets.InheritFrom(otherCT.Presets)

	ct.Content.InheritFrom(otherCT.Content)

	return nil
}

// resolve resolves this template. This means that preset dependencies are resolved
// and after that the preset dependencies on content side. Currently nothing needs
// to be done for parameters.
func (ct *ChronicleTemplate) resolve() (err error) {
	if err = ct.Presets.Resolve(); err != nil {
		return err
	}

	if err = ct.Content.Resolve(ct.Presets); err != nil {
		return err
	}
	return nil
}

// WriteToCsvFile creates a CSV file out of the current chronicle template than can be used
// as input for the "batch fill" command
func (ct *ChronicleTemplate) WriteToCsvFile(filename string, separator rune, as *args.Store) (err error) {
	const numPlayers = 7

	records := [][]string{
		{"#ID", ct.ID},
		{"#Description", ct.Description},
		{"#"},
		{"#Players"}, // will be filled below with labels
	}
	for idx := 1; idx <= numPlayers; idx++ {
		outerIdx := len(records) - 1
		records[outerIdx] = append(records[outerIdx], fmt.Sprintf("Player %d", idx))
	}

	for _, contentID := range ct.Parameters.GetSortedKeys() {
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
	// as we add new entries to the argStore, create a local store and set the
	// original store as parent.
	localArgStore, err := args.NewStore(args.StoreInit{Parent: argStore})
	if err != nil {
		return err
	}

	// check argStore values against parameter definitions
	if err = ct.Parameters.ValidateAndProcessArgs(localArgStore); err != nil {
		return err
	}

	// pass to content store to generate output
	if err = ct.Content.GenerateOutput(stamp, localArgStore); err != nil {
		return err
	}

	return nil
}

// IsValid checks whether a given chronicle is valid. This should only be called
// after resolve() was called on this template.
func (ct *ChronicleTemplate) IsValid() (err error) {
	if !utils.IsSet(ct.Description) {
		return templateErrf(ct, "Missing description")
	}

	if err = ct.Parameters.IsValid(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Content.IsValid(); err != nil {
		return templateErr(ct, err)
	}

	return nil
}

// Describe returns a short textual description of a single chronicle template.
// It returns the description as a multi-line string.
func (ct *ChronicleTemplate) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", ct.ID)
		if utils.IsSet(ct.Description) {
			fmt.Fprintf(&sb, ": %v", ct.Description)
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", ct.ID)
		fmt.Fprintf(&sb, "\tDesc: %v\n", ct.Description)
		fmt.Fprintf(&sb, "\tFile: %v", ct.filename)
	}

	return sb.String()
}

// DescribeParams returns a textual description of the parameters expected by
// this chronicle template. It returns the description as a multi-line string.
func (ct *ChronicleTemplate) DescribeParams(verbose bool) (result string) {
	return ct.Parameters.Describe(verbose)
}
