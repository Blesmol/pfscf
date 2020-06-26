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

		Args: cobra.ExactArgs(0),
	}

	templateListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Long:  "Provide a list of all locally available templates",

		Args: cobra.ExactArgs(0),

		Run: executeTemplateList,
	}
	templateCmd.AddCommand(templateListCmd)

	templateDescribeCmd := &cobra.Command{
		Use:   "describe <template>",
		Short: "Describe a specific template",
		Long:  "Describe a specific template by listing all available fields from this template along with their description",

		Args: cobra.ExactArgs(1),

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
	ts, err := GetTemplateStore()
	ExitOnError(err, "Could not read templates")

	templateNames := ts.GetTemplateNames()
	fmt.Printf("List of templates:\n")
	for _, templateName := range templateNames {
		template, _ := ts.GetTemplate(templateName)
		fmt.Printf("- %v: %v\n", template.Name(), template.Description())
	}
}

func executeTemplateDescribe(cmd *cobra.Command, args []string) {
	templateName := args[0]

	ct, err := GetTemplate(templateName)
	ExitOnError(err, "Could not get template '%v'", templateName)

	fmt.Printf("Template '%v'\n", templateName)
	idList := ct.GetContentIDs()
	for _, id := range idList {
		ce, _ := ct.GetContent(id)
		fmt.Printf("- %v: %v\n", id, ce.Desc)
		// TODO add example input to output. Example should contain the id
		// and a CE-specific example value that needs to be included in the yaml file
	}
}

func executeTemplateUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}
