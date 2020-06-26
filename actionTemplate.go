package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GetTemplateCommand returns the cobra command for the "fill" action.
func GetTemplateCommand() (cmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Various actions on templates: list, describe, etc",
		Long:  "Long description of actions on templates.",

		//Args: cobra.MinimumNArgs(1),

		//Run: executeTemplate,
	}

	templateListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Long:  "Provide a list of all locally available templates",

		Run: executeTemplateList,
	}
	templateCmd.AddCommand(templateListCmd)

	templateDescribeCmd := &cobra.Command{
		Use:   "describe <template>",
		Short: "Describe a specific template",
		Long:  "Describe a specific template by listing all available fields from this template along with their description",

		Args: cobra.MinimumNArgs(1),

		Run: executeTemplateDescribe,
	}
	templateCmd.AddCommand(templateDescribeCmd)

	templateUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update the locally available templates",
		Long:  "Update the locally available templates with the latest templates from the central github repository",

		Run: executeTemplateUpdate,
	}
	templateCmd.AddCommand(templateUpdateCmd)

	return templateCmd
}

func executeTemplateList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		ExitWithMessage("Unrecognized command line arguments: %v", args)
	}

	ts, err := GetTemplateStore()
	ExitOnError(err, "Could not read templates")

	templateNames := ts.GetKeys()
	fmt.Printf("List of templates:\n")
	for _, templateName := range templateNames {
		fmt.Printf("- %v\n", templateName)
		// add Description if present
	}
}

func executeTemplateDescribe(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}

func executeTemplateUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}
