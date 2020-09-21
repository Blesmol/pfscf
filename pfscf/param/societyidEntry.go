package param

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	typeSocietyID = "societyid"
)

var (
	regexSocietyID = regexp.MustCompile(`^\s*(\d*)\s*-\s*(\d*)\s*$`)
)

type societyidEntry struct {
	id   string
	theRank int

	TheExample     string `yaml:"example"`
	TheDescription string `yaml:"description"`
}

func (e *societyidEntry) Type() string {
	return typeSocietyID
}

func (e *societyidEntry) ID() string {
	return e.id
}

func (e *societyidEntry) setID(id string) {
	utils.Assert(!utils.IsSet(e.id), "Should only be called once per object")
	e.id = id
}

func (e *societyidEntry) Example() string {
	return e.TheExample
}

func (e *societyidEntry) Description() string {
	return e.TheDescription
}

func (e *societyidEntry) rank() int {
	return e.theRank
}

func (e *societyidEntry) setRank(rank int) {
	e.theRank = rank
}

func (e *societyidEntry) deepCopy() Entry {
	copy := *e
	return &copy
}

func (e *societyidEntry) isValid() (err error) {
	if !utils.IsSet(e.TheExample) {
		return fmt.Errorf("Missing example")
	}
	if !utils.IsSet(e.TheDescription) {
		return fmt.Errorf("Missing description")
	}
	return nil
}

func (e *societyidEntry) validateAndProcessArgs(as *args.Store) (err error) {
	argValue, exists := as.Get(e.ID())
	utils.Assert(exists, "Existence of entry should have been validated by caller")

	// check and split up provided society id value
	societyID := regexSocietyID.FindStringSubmatch(argValue)
	if len(societyID) == 0 {
		return fmt.Errorf("Provided society ID does not follow the pattern '<player_id>-<char_id>': '%v'", argValue)
	}
	utils.Assert(len(societyID) == 3, "Should contain the matching text plus the capturing groups")
	playerID := societyID[1]
	charID := societyID[2]

	// add to arg store
	as.Set(e.ID()+".player", playerID)
	as.Set(e.ID()+".char", charID)
	// TODO validate that no such entries yet exist in the argStore

	return nil
}

func (e *societyidEntry) describe(verbose bool) (result string) {
	var sb strings.Builder

	if !verbose {
		fmt.Fprintf(&sb, "- %v: %v\n", e.id, e.TheDescription)
	} else {
		fmt.Fprintf(&sb, "- %v\n", e.id)
		fmt.Fprintf(&sb, "\tDesc: %v\n", e.TheDescription)
		fmt.Fprintf(&sb, "\tType: %v\n", e.Type())
		fmt.Fprintf(&sb, "\tExample: %v\n", genericContentUsageExample(e.id, e.TheExample))
	}

	return sb.String()
}
