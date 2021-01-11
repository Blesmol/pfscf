package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/canvas"
	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/content"
	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	floatGroupPattern  = `(\d+(?:\.(?:(?:\d*))?)?)` // should match floating point numbers, e.g. 2, 1., 45.123
	aspectRatioPattern = `^\s*` + floatGroupPattern + `\s*:\s*` + floatGroupPattern + `\s*$`
)

var (
	regexAspectRatio = regexp.MustCompile(aspectRatioPattern)
	validFlags       = []string{"hidden"}
)

// Chronicle is the new approach for the Chronicle Template
type Chronicle struct {
	ID          string
	Description string
	Parent      string
	Aspectratio string
	Flags       []string
	Parameters  param.Store
	Presets     preset.Store
	Canvas      canvas.Store
	Content     content.ListStore

	filename string // filename of the originating yaml file

	parent   *Chronicle
	children []*Chronicle
}

// NewChronicleTemplate returns a new ChronicleTemplate object.
func NewChronicleTemplate(filename string) (ct Chronicle) {
	ct.filename = filename
	return ct
}

func templateErr(ct *Chronicle, errIn error) (errOut error) {
	return fmt.Errorf("Template '%v': %v", ct.ID, errIn)
}

func templateErrf(ct *Chronicle, msg string, args ...interface{}) (errOut error) {
	return fmt.Errorf("Template '%v': "+msg, ct.ID, args)
}

// ensureStoresAreInitialized is a workaround for the behavior of the stupid f... yaml library.
// If a section like "parameters:" is present, but empty, it will not be unmarshalled, and the
// underlying data structure will be ZEROed. So the stores will be uninitialized. Even if they
// were initialized before the unmarshalling. Yeah, great.
// See https://github.com/go-yaml/yaml/issues/395 , might be fixed with go-yaml v3 in the future.
func (ct *Chronicle) ensureStoresAreInitialized() {
	if ct.Parameters == nil {
		ct.Parameters = param.NewStore()
	}
	if ct.Presets == nil {
		ct.Presets = preset.NewStore()
	}
	if ct.Canvas == nil {
		ct.Canvas = canvas.NewStore()
	}
	if ct.Content == nil {
		ct.Content = content.NewListStore()
	}

	if ct.children == nil {
		ct.children = make([]*Chronicle, 0)
	}
}

// GetExampleArguments returns an array containing all keys and example values for all parameters.
// The result can be passed to the ArgStore.
func (ct *Chronicle) GetExampleArguments() (result []string) {
	return ct.Parameters.GetExampleArguments()
}

// inheritFrom inherits entries from multiple sections from another
// ChronicleTemplate object. An error is returned in case a content
// entry from sections 'parameters' or 'content' exists in both objects.
// In case a preset entry exists in both objects, then the one from the original
// object takes precedence.
func (ct *Chronicle) inheritFrom(otherCT *Chronicle) (err error) {
	err = ct.Parameters.InheritFrom(&otherCT.Parameters)
	if err != nil {
		return templateErr(ct, err)
	}

	ct.Presets.InheritFrom(otherCT.Presets)

	ct.Canvas.InheritFrom(otherCT.Canvas)

	ct.Content.InheritFrom(otherCT.Content)

	if !utils.IsSet(ct.Aspectratio) {
		ct.Aspectratio = otherCT.Aspectratio
	}

	return nil
}

// resolve resolves this template. This means that preset dependencies are resolved
// and after that the preset dependencies on content side. Currently nothing needs
// to be done for parameters.
func (ct *Chronicle) resolve() (err error) {
	if err = ct.Presets.Resolve(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Canvas.Resolve(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Content.Resolve(ct.Presets); err != nil {
		return templateErr(ct, err)
	}
	return nil
}

// GenerateCsvFile creates a CSV file out of the current chronicle template than can be used
// as input for the "batch fill" command
func (ct *Chronicle) GenerateCsvFile(filename string, separator rune, argStore *args.Store, cmdFlags [][]string) (err error) {
	const numPlayers = 7
	const numChronicles = numPlayers + 1     // GM also wants a chronicle
	const numColumns = 1 + numChronicles + 1 // identifiers + chronicles + example column

	// file header
	records := [][]string{
		{"# ID", ct.ID},
		{"# Description", ct.Description},
		{""},
	}

	// add command line flags
	records = append(records, []string{"# Command line arguments"})
	for _, entry := range cmdFlags {
		utils.Assert(len(entry) == 2, "Number of entries is wrong")
		records = append(records, []string{entry[0], entry[1]})
	}
	records = append(records, []string{""})

	// Add players section with "Player <nr>" labels
	records = append(records, []string{"# Players"})
	outerIdx := len(records) - 1
	for idx := 1; idx <= numPlayers; idx++ {
		records[outerIdx] = append(records[outerIdx], fmt.Sprintf("Player %d", idx))
	}
	records[outerIdx] = append(records[outerIdx], "GM")        // Add "GM" label as well
	records[outerIdx] = append(records[outerIdx], "# Example") // Add "GM" label as well

	// fill from parameters
	for _, groupName := range ct.Parameters.GetGroupsSortedByRank() {
		// Header for the current group
		records = append(records, []string{""})
		records = append(records, []string{fmt.Sprintf("# %v", groupName)})

		// add parameters from current group
		for _, paramID := range ct.Parameters.GetKeysForGroupSortedByRank(groupName) {
			paramEntry, _ := ct.Parameters.Get(paramID)

			// parameters can have multiple identifiers, e.g. for splitlines.
			for _, argStoreID := range paramEntry.ArgStoreIDs() {
				// entry should be large enough for id column + number of chronicles
				row := make([]string, numColumns)

				row[0] = argStoreID // first column is always parameter name

				// check if some value was provided on the cmd line that should be filled in all columns for this parameter
				if val, exists := argStore.Get(argStoreID); exists {
					for colIdx := 1; colIdx <= numChronicles; colIdx++ {
						row[colIdx] = val
					}
				}

				// add example text
				row[len(row)-1] = fmt.Sprintf("# %v", paramEntry.Example())

				records = append(records, row)
			}
		}
	}

	// add parameter legend to end of file
	records = append(records, []string{""})
	records = append(records, []string{"# Legend for input values:"})
	records = append(records, []string{"# Name", "Accepted values", "Example", "Description"})
	for _, paramName := range ct.Parameters.GetKeysSortedByName() {
		param, _ := ct.Parameters.Get(paramName)

		entry := make([]string, 4) // Comment char, Name, type, example, description

		entry[0] = "# " + param.ID()
		entry[1] = utils.ToCommaSeparatedString(param.AcceptedValues())
		entry[2] = param.Example()
		entry[3] = param.Description()

		records = append(records, entry)
	}

	err = csv.WriteFile(filename, separator, records)
	if err != nil {
		return err
	}

	return nil
}

// GenerateOutput adds the content of this chronicle template to the provided stamp.
func (ct *Chronicle) GenerateOutput(stamp *stamp.Stamp, argStore *args.Store) (err error) {
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

	if utils.IsSet(ct.Aspectratio) {
		xMarginPct, yMarginPct, err := ct.guessMarginsFromAspectRatio(stamp)
		if err != nil {
			return err
		}

		stamp.SetPageCanvas(0.0+xMarginPct/2, 0.0+yMarginPct/2, 100.0-xMarginPct/2, 100-yMarginPct/2)
	}

	ct.Canvas.AddCanvasesToStamp(stamp)

	// pass to content store to generate output
	if err = ct.Content.GenerateOutput(stamp, localArgStore); err != nil {
		return err
	}

	// draw canvas borders as last action to be visible over other content
	if cfg.Global.DrawCanvas {
		stamp.DrawCanvases()
	}

	return nil
}

// IsValid checks whether a given chronicle is valid. This should only be called
// after resolve() was called on this template.
func (ct *Chronicle) IsValid() (err error) {
	if utils.IsSet(ct.Aspectratio) {
		if _, _, err = parseAspectRatio(ct.Aspectratio); err != nil {
			return templateErr(ct, err)
		}
	}

	if !utils.IsSet(ct.Description) {
		return templateErrf(ct, "Missing description")
	}

	if err = ct.hasValidFlags(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Parameters.IsValid(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Canvas.IsValid(); err != nil {
		return templateErr(ct, err)
	}

	if err = ct.Content.IsValid(&ct.Parameters, &ct.Canvas); err != nil {
		return templateErr(ct, err)
	}

	return nil
}

// Describe returns a short textual description of a single chronicle template.
// It returns the description as a multi-line string.
func (ct *Chronicle) Describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v", ct.ID)
		if utils.IsSet(ct.Description) {
			fmt.Fprintf(&sb, ": %v", ct.Description)
		}
	} else {
		fmt.Fprintf(&sb, "- %v\n", ct.ID)
		fmt.Fprintf(&sb, "\tDescription: %v\n", ct.Description)
		fmt.Fprintf(&sb, "\tFile: %v\n", ct.filename)
	}

	return sb.String()
}

// DescribeParams returns a textual description of the parameters expected by
// this chronicle template. It returns the description as a multi-line string.
func (ct *Chronicle) DescribeParams(verbose bool) (result string) {
	return ct.Parameters.Describe(verbose)
}

func parseAspectRatio(input string) (x, y float64, err error) {
	match := regexAspectRatio.FindStringSubmatch(input)
	if len(match) == 0 {
		err = fmt.Errorf("Provided aspect ratio does not follow pattern '<x>:<y>': %v", input)
		return
	}

	if x, err = strconv.ParseFloat(match[1], 64); err != nil {
		err = fmt.Errorf("Error parsing X part of aspect ratio '%v': %v", match[1], err)
	}

	if y, err = strconv.ParseFloat(match[2], 64); err != nil {
		err = fmt.Errorf("Error parsing Y part of aspect ratio '%v': %v", match[2], err)
	}

	return x, y, nil
}

// guessMarginsFromAspectRatio tries to calculate possible document margins from the provided
// aspect ratio. Assumption is that the PDF content will not be squeezed or stretched in any
// direction. So if the aspect ration differs this must mean than margins were added somewhere.
// The following function should correctly calculate the margins from the aspect ratio if
// margins were only added on the x axis OR the y axis, not on both.
func (ct *Chronicle) guessMarginsFromAspectRatio(stamp *stamp.Stamp) (xMarginPct, yMarginPct float64, err error) {
	sx, sy := stamp.GetDimensionsWithOffset()
	arx, ary, err := parseAspectRatio(ct.Aspectratio)
	if err != nil {
		return
	}

	haveAR := sx / sy
	wantAR := arx / ary

	switch {
	case wantAR > haveAR: // y axis has a margin
		f := sx * ary / arx
		g := (100.0 * f) / sy
		marginPct := 100 - g
		return 0.0, marginPct, nil
	case wantAR < haveAR: // x axis has a margin
		f := arx * sy / ary
		g := (100.0 * f) / sx
		marginPct := 100.0 - g
		//fmt.Printf("f: %.6f, g: %.6f, h: %.6f, test: %.6f\n", f, g, h, 100.0*609.9/612.0)
		return marginPct, 0.0, nil
	}

	return 0.0, 0.0, nil // no margins, fits perfect
}

func (ct *Chronicle) addChild(childCt *Chronicle) (err error) {
	// add references to chronicle templates
	utils.Assert(childCt.parent == nil, "Chronicle can only have one 'inherit' entry, thus can only have one parent")
	childCt.parent = ct
	ct.children = append(ct.children, childCt)

	// check for cyclic dependencies
	depList := make([]string, 0)
	for curCt := ct; curCt.parent != nil; curCt = curCt.parent {
		if utils.Contains(depList, curCt.ID) {
			depList = append(depList, curCt.ID) // add entry before printing to have complete cycle in output
			return templateErrf(curCt, "Found cyclic inheritance dependencies. Inheritance chain is %v", depList)
		}
		depList = append(depList, curCt.ID)
	}

	// ensure that list of children is sorted lexically
	sortChronicleList(ct.children)

	return nil
}

func (ct *Chronicle) getHierarchieLevel(excludeHidden bool) (level uint) {
	for curCt := ct.parent; curCt != nil; curCt = curCt.parent {
		if !curCt.hasFlag("hidden") {
			level++
		}
	}
	return level
}

func (ct *Chronicle) hasValidFlags() (err error) {
	for _, flag := range ct.Flags {
		if !utils.Contains(validFlags, flag) {
			return fmt.Errorf("Unknown flag: '%v'", flag)
		}
	}
	return nil
}

func (ct *Chronicle) hasFlag(flag string) bool {
	return utils.Contains(ct.Flags, flag)
}
