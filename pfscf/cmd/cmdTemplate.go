package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// GetTemplateCommand returns the cobra command for the "fill" action.
func GetTemplateCommand() (cmd *cobra.Command) {
	templateCmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"t", "templates"},

		Short: "Various actions on templates: list, describe, etc",
		//Long:  "TBD",

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

	templateSearchCmd := &cobra.Command{
		Use:     "search <search term>",
		Aliases: []string{"s"},

		Short: "Search for templates",
		Long:  "Search for specific templates by listing all templates where the id or description match the provided search term. The search is case-insensitive.",

		Args: cobra.MinimumNArgs(1),

		Run: executeTemplateSearch,
	}
	templateCmd.AddCommand(templateSearchCmd)

	/*
		templateValidateCmd := &cobra.Command{
			Use:     "validate <template>",
			Aliases: []string{"v"},

			Short: "Validate a specific template",
			//Long:  "TBD",

			Args: cobra.ExactArgs(1),

			Run: executeTemplateValidate,
		}
		templateCmd.AddCommand(templateValidateCmd)
	*/

	/*
		templateUpdateCmd := &cobra.Command{
			Use:     "update",
			Aliases: []string{"u"},

			Short: "Update the locally available templates",
			Long:  "Update the locally available templates with the latest templates from the central github repository",

			Run: executeTemplateUpdate,
		}
		templateCmd.AddCommand(templateUpdateCmd)
	*/

	return templateCmd
}

func executeTemplateList(cmd *cobra.Command, args []string) {
	ts, err := template.GetStore()
	utils.ExitOnError(err, "Could not read templates")

	fmt.Printf("List of available templates:\n\n")
	fmt.Printf(ts.ListTemplates())
}

func executeTemplateDescribe(cmd *cobra.Command, args []string) {
	templateName := args[0]

	ts, err := template.GetStore()
	utils.ExitOnError(err, "Could not read templates")

	ct, exists := ts.Get(templateName)
	if !exists {
		utils.ExitWithMessage("Coud not find template '%v'", templateName)
	}

	fmt.Printf("Template '%v'\n\n", templateName)
	fmt.Printf(ct.DescribeParams(cfg.Global.Verbose))
}

func executeTemplateSearch(cmd *cobra.Command, args []string) {
	ts, err := template.GetStore()
	utils.ExitOnError(err, "Could not read templates")

	if matchDesc, foundMatch := ts.SearchForTemplates(args...); foundMatch {
		fmt.Printf("Matching Templates:\n")
		fmt.Print(matchDesc)
	} else {
		fmt.Println("Found no matching templates")
	}
}

func executeTemplateValidate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}

func executeTemplateUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}
