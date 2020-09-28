package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

// GetOpenCommand returns the cobra command for the "open" action
func GetOpenCommand() (cmd *cobra.Command) {
	openCmd := &cobra.Command{
		Use:     "open <file>",
		Aliases: []string{"o"},

		Hidden: false,

		Short: "Open the specified file (e.g. PDF, CSV) with its default app",
		Long:  "Open the specified file with its default app",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, cmdArgs []string) {
			err := utils.OpenWithDefaultViewer(cmdArgs[0])
			utils.ExitOnError(err, "Error opening file")
		},
	}

	return openCmd
}
