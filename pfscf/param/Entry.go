package param

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Entry defines the interface for all param entries
type Entry interface {
	ID() string
	setID(string)
	Type() string
	Example() string
	Description() string
	AcceptedValues() []string
	Group() string
	setGroup(string)
	rank() int
	setRank(int)
	deepCopy() Entry
	isValid() error
	validateAndProcessArgs(*args.Store) error
	describe(bool) string
}

func genericContentUsageExample(id, exampleValue string) (result string) {
	return fmt.Sprintf("%v=%v", id, utils.QuoteStringIfRequired(exampleValue))
}
