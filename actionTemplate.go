package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GetTemplateCommand returns the cobra command for the "fill" action.
func GetTemplateCommand() (cmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"t", "templates"},

		Short: "Various actions on templates: list, describe, etc",
		Long:  "Long description of actions on templates.",

		Args: cobra.ExactArgs(0),
	}

	templateListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},

		Short: "List available templates",
		Long:  "Provide a list of all locally available templates",

		Args: cobra.ExactArgs(0),

		Run: executeTemplateList,
	}
	templateCmd.AddCommand(templateListCmd)

	templateDescribeCmd := &cobra.Command{
		Use:     "describe <template>",
		Aliases: []string{"d", "desc"},

		Short: "Describe a specific template",
		Long:  "Describe a specific template by listing all available fields from this template along with their description",

		Args: cobra.ExactArgs(1),

		Run: executeTemplateDescribe,
	}
	templateCmd.AddCommand(templateDescribeCmd)

	templateValidateCmd := &cobra.Command{
		Use:     "validate <template>",
		Aliases: []string{"v"},

		Short: "Validate a specific template",
		Long:  "Validate a specific template",

		Args: cobra.ExactArgs(1),

		Run: executeTemplateValidate,
	}
	templateCmd.AddCommand(templateValidateCmd)

	templateUpdateCmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},

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
	fmt.Printf("\n")
	fmt.Printf("List of available templates:\n\n")
	for _, templateName := range templateNames {
		template, _ := ts.GetTemplate(templateName)
		fmt.Println(template.Describe(flags.verbose))
	}
}

func executeTemplateDescribe(cmd *cobra.Command, args []string) {
	templateName := args[0]

	ct, err := GetTemplate(templateName)
	ExitOnError(err, "Could not get template '%v'", templateName)

	fmt.Printf("Template '%v'\n\n", templateName)
	idList := ct.GetContentIDs() // TODO only display non-alias IDs
	for _, id := range idList {
		ce, _ := ct.GetContent(id)
		fmt.Println(ce.Describe(flags.verbose))
		// TODO add aliases
	}
}

func executeTemplateValidate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}

func executeTemplateUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}
