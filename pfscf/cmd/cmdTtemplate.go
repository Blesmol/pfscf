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

	templateValidateCmd := &cobra.Command{
		Use:     "validate <template>",
		Aliases: []string{"v"},

		Short: "Validate a specific template",
		//Long:  "TBD",

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
	ts, err := template.GetStore()
	utils.ExitOnError(err, "Could not read templates")

	templateNames := ts.GetTemplateIDs(false)
	fmt.Printf("\n")
	fmt.Printf("List of available templates:\n\n")
	for _, templateName := range templateNames {
		template, _ := ts.Get(templateName)
		fmt.Println(template.Describe(cfg.Global.Verbose))
	}
}

func executeTemplateDescribe(cmd *cobra.Command, args []string) {
	templateName := args[0]

	ct, err := template.Get(templateName)
	utils.ExitOnError(err, "Could not get template '%v'", templateName)

	fmt.Printf("Template '%v'\n\n", templateName)
	idList := ct.GetContentIDs(false)
	for _, id := range idList {
		ce, _ := ct.GetContent(id)
		fmt.Println(ce.Describe(cfg.Global.Verbose))
	}
}

func executeTemplateValidate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}

func executeTemplateUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Not yet implemented")
}
