package content

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	regexParamValue = regexp.MustCompile(`^\s*param:\s*(\S*)$`)
)

// getValue returns the value that should be used for the current content.
func getValue(valueField string, as *args.Store) (result *string) {
	// No input? No result!
	if !utils.IsSet(valueField) {
		return nil
	}

	// check whether a parameter reference was provided, i.e. something like "param:<name>"
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

// getMultiValue returns an array of values that should be used for the current content.
func getMultiValue(contentValueField string, as *args.Store) (result []string) {
	// No input? No result!
	if !utils.IsSet(contentValueField) {
		return nil
	}

	// check whether a parameter reference was provided, i.e. something like "param:<name>"
	paramName := regexParamValue.FindStringSubmatch(contentValueField)
	if len(paramName) > 0 {
		utils.Assert(len(paramName) == 2, "Should contain the matching text plus a single capturing group")

		argArray := as.GetArray(paramName[1])
		if len(argArray) > 0 {
			return argArray
		}
		return nil
	}

	// else assume that provided value was a static text
	return []string{contentValueField}
}

func parseColor(color string) (r, g, b int, err error) {
	regexHexColorCode := regexp.MustCompile(`^[0-9a-f]{6}$`)

	color = strings.ToLower(strings.TrimSpace(color))

	switch color {
	case "white":
		return 255, 255, 255, nil
	case "black":
		return 0, 0, 0, nil
	case "blue":
		return 0, 0, 255, nil
	case "red":
		return 255, 0, 0, nil
	case "green":
		return 0, 255, 0, nil
	}

	colorCode := regexHexColorCode.FindString(color)
	if utils.IsSet(colorCode) {
		colorCodeBytes := []byte(colorCode)
		decoded := make([]byte, hex.DecodedLen(len(colorCodeBytes)))
		_, err := hex.Decode(decoded, colorCodeBytes)
		utils.Assert(err == nil, fmt.Sprintf("Valid input should have been guaranteed by regexp, but instead got error: %v", err))
		utils.Assert(len(decoded) == 3, fmt.Sprintf("Number of resultint entries should be guaranteed by regexp, was %v instead", len(decoded)))

		r, g, b = int(decoded[0]), int(decoded[1]), int(decoded[2])
		return r, g, b, nil
	}

	return 0, 0, 0, fmt.Errorf("Unknown color: '%v'", color)
}
