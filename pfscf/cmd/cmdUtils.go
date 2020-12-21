package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func warnOnWrongFileExtension(filename, expectedExt string) {
	realExt := strings.ToLower(filepath.Ext(filename))
	if realExt != strings.ToLower("."+expectedExt) {
		fmt.Fprintf(os.Stderr, "Warning: File '%v' does not have expected extension '%v'\n", filename, expectedExt)
	}
}

// setFlagsFromRecords searches for records that begin with "--", interprets anything behind that as
// command line flag name, and tries to set it to the value provided in the next record.
func setFlagsFromRecords(cmd *cobra.Command, records [][]string) error {
	// check whether there is any content at all
	if len(records) == 0 || len(records[0]) < 2 {
		return nil
	}

	// TODO add check if flag was provided but no value? Or at least warning

	for _, record := range records {
		flagCandidate := record[0]
		marker := "flag:--"
		if strings.HasPrefix(flagCandidate, marker) {
			flagName := flagCandidate[len(marker):]
			flagValue := record[1]

			if err := cmd.Flags().Set(flagName, flagValue); err != nil {
				return fmt.Errorf("Error setting flag '%v' with value '%v': %v", flagName, flagValue, err)
			}
		}
	}

	return nil
}
