package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func warnOnWrongFileExtension(filename, expectedExt string) {
	realExt := strings.ToLower(filepath.Ext(filename))
	if realExt != strings.ToLower("."+expectedExt) {
		fmt.Fprintf(os.Stderr, "Warning: File '%v' does not have expected extension '%v'\n", filename, expectedExt)
	}
}
