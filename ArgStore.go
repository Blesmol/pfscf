package main

import (
	"regexp"
	"strings"
)

// ArgStore holds a mapping between argument keys and values
type ArgStore map[string]*string

const (
	// TODOs:
	// - Restrict input set for key: [a-zA-Z0-9-_]
	// - key should have shortest match, to allow value to contain "="
	argPattern = "(?P<key>.*)=(?P<value>.*)"
)

var (
	argRegex *regexp.Regexp // not used at the moment
)

func init() {
	argRegex = regexp.MustCompile(argPattern)
}

// ParseArgs takes a list of provided arguments and checks
// that they are in format "<key>=<value>". It will return
// an error on duplicate keys.
func ParseArgs(args []string) (as ArgStore) {
	as = make(ArgStore, len(args))

	for _, arg := range args {
		splitIdx := strings.Index(arg, "=")
		Assert(splitIdx != -1, "No '=' found in argument")
		Assert(splitIdx != 0, "No key found in argument")
		Assert(splitIdx != (len(arg)-1), "No value found in argument")

		key := arg[:splitIdx]
		value := arg[(splitIdx + 1):]

		if _, exists := as[key]; !exists {
			as[key] = &value
		} else {
			panic("Duplicate key found: " + key)
		}
	}

	return as
}
