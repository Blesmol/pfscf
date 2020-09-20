package content

import (
	"fmt"
	"regexp"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/canvas"
	"github.com/Blesmol/pfscf/pfscf/param"
	"github.com/Blesmol/pfscf/pfscf/preset"
	"github.com/Blesmol/pfscf/pfscf/stamp"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	regexParamValue = regexp.MustCompile(`^\s*param:\s*(\S*)$`)
)

// Entry is an interface for the content. D'oh!
type Entry interface {
	isValid(*param.Store, *canvas.Store) (err error)
	resolve(ps preset.Store) (err error)
	generateOutput(s *stamp.Stamp, as *args.Store) (err error)
	deepCopy() Entry
}

// TODO now with no ID we should print all fields of the respective entry instead (don't forget the type)
func contentValErr(ce Entry, errIn error) (errOut error) {
	return fmt.Errorf("Error validating content: %v; complete content entry is: %v", errIn, ce)
}

// getValue returns the value that should be used for the current content.
func getValue(valueField string, as *args.Store) (result *string) {
	// No input? No result!
	if !utils.IsSet(valueField) {
		return nil
	}

	// check whether a parameter reference was provided
	paramName := regexParamValue.FindStringSubmatch(valueField)
	if len(paramName) > 0 {
		utils.Assert(len(paramName) == 2, "Should contain the matching text plus a single capturing group")
		argValue, exists := as.Get(paramName[1])
		if exists {
			return &argValue
		}
		return nil
	}

	// else assume that provided value was a static text
	return &valueField
}
