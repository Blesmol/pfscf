package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	actionBatchCreateUseExampleValues bool
	actionBatchCreateSeparator        string
)

// GetBatchCommand returns the cobra command for the "batch" action.
func GetBatchCommand() (cmd *cobra.Command) {
	batchCmd := &cobra.Command{
		Use:     "batch",
		Aliases: []string{"b"},

		Short: "Fill out multiple chronicles in one go",
		Long:  "The batch operation can fill out multiple chronicles in one go by reading all necessary input from a csv file.",

		Args: cobra.ExactArgs(0),
	}

	batchCreateCmd := &cobra.Command{
		Use:     "create <template> <output> [<content_id>=<value> ...]",
		Aliases: []string{"c"},

		Short: "Create ready-to-fill csv file based on selected template",
		//Long:  "TBD",

		Args: cobra.MinimumNArgs(2),

		Run: executeBatchCreate,
	}
	batchCreateCmd.Flags().BoolVarP(&actionBatchCreateUseExampleValues, "exampleValues", "e", false, "Use example values to fill out the chronicle")
	batchCreateCmd.Flags().StringVarP(&actionBatchCreateSeparator, "separator", "s", ";", "Field separator character for resulting CSV file")

	batchCmd.AddCommand(batchCreateCmd)

	batchFillCmd := &cobra.Command{
		Use:     "fill <template> <csv_file> <input_pdf> <output_dir> [<content_id>=<value> ...]",
		Aliases: []string{"f"},

		Short: "Fill multiple templates with values read from a csv file",
		//Long:  "TBD",

		Args: cobra.ExactArgs(4),

		Run: executeBatchFill,
	}
	batchCmd.AddCommand(batchFillCmd)

	return batchCmd
}

func executeBatchCreate(cmd *cobra.Command, args []string) {
	utils.Assert(len(args) >= 2, "Number of arguments should be guaranteed by cobra settings")

	tmplName := args[0]
	outFile := args[1]

	var separator rune
	switch actionBatchCreateSeparator {
	case ";":
		separator = ';'
	case ",":
		separator = ','
	default:
		utils.ExitWithMessage("Currently only ';' and ',' are accepted as separators")
	}
	fmt.Printf("Separator: %v", separator)

	cTmpl, err := GetTemplate(tmplName)
	utils.ExitOnError(err, "Error getting template")

	// parse remaining arguments
	var argStore *ArgStore
	if !actionBatchCreateUseExampleValues {
		argStore = ArgStoreFromArgs(args[2:])
	} else {
		argStore = ArgStoreFromTemplateExamples(cTmpl)
	}

	err = cTmpl.WriteTemplateToCsvFile(outFile, argStore, separator)
	utils.ExitOnError(err, "Error writing CSV file for template %v", tmplName)
}

func executeBatchFill(cmd *cobra.Command, args []string) {
	utils.Assert(len(args) >= 4, "Number of arguments should be guaranteed by cobra settings")

	tmplName := args[0]
	inCsv := args[1]
	inPdf := args[2]
	outDir := args[3]

	cTmpl, err := GetTemplate(tmplName)
	utils.ExitOnError(err, "Error getting template")

	batchArgStores, err := GetFillInformationFromCsvFile(inCsv)
	utils.ExitOnError(err, "Error reading csv file")

	// parse remaining arguments
	cmdLineArgStore := ArgStoreFromArgs(args[3:])

	for idx, batchArgStore := range batchArgStores {
		cmdLineArgStore.SetParent(batchArgStore) // command line arguments have priority

		pdf, err := NewPdf(inPdf)
		utils.ExitOnError(err, "Error opening input file '%v'", inPdf)

		playerNumber := idx + 1
		baseOutfile := fmt.Sprintf("Chronicle_Player_%d.pdf", playerNumber)
		outfile := filepath.Join(outDir, baseOutfile)

		err = pdf.Fill(cmdLineArgStore, cTmpl, outfile)
		utils.ExitOnError(err, "Error when filling out chronicle for player %d", playerNumber)
	}

}
