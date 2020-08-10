package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/args"
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

func executeBatchCreate(cmd *cobra.Command, cmdArgs []string) {
	utils.Assert(len(cmdArgs) >= 2, "Number of arguments should be guaranteed by cobra settings")

	tmplName := cmdArgs[0]
	outFile := cmdArgs[1]

	var separator rune
	// TODO remove check completely, just use first rune in separator string. Or forward as string instead of rune, and there check length and values.
	switch actionBatchCreateSeparator {
	case ";":
		separator = ';'
	case ",":
		separator = ','
	default:
		utils.ExitWithMessage("Currently only ';' and ',' are accepted as separators")
	}

	cTmpl, err := GetTemplate(tmplName)
	utils.ExitOnError(err, "Error getting template")

	// parse remaining arguments
	var argStore *args.ArgStore
	if !actionBatchCreateUseExampleValues {
		argStore, err = args.NewArgStore(args.ArgStoreInit{Args: cmdArgs[2:]})
	} else {
		argStore, err = args.NewArgStore(args.ArgStoreInit{Args: cTmpl.GetExampleArguments()})
	}
	utils.ExitOnError(err, "Error processing command line arguments")

	err = cTmpl.WriteToCsvFile(outFile, separator, argStore)
	utils.ExitOnError(err, "Error writing CSV file for template %v", tmplName)
}

func executeBatchFill(cmd *cobra.Command, cmdArgs []string) {
	utils.Assert(len(cmdArgs) >= 4, "Number of arguments should be guaranteed by cobra settings")

	tmplName := cmdArgs[0]
	inCsv := cmdArgs[1]
	inPdf := cmdArgs[2]
	outDir := cmdArgs[3]

	cTmpl, err := GetTemplate(tmplName)
	utils.ExitOnError(err, "Error getting template")

	batchArgStores, err := args.GetArgStoresFromCsvFile(inCsv)
	utils.ExitOnError(err, "Error reading csv file")

	// parse remaining arguments
	cmdLineArgStore, err := args.NewArgStore(args.ArgStoreInit{Args: cmdArgs[3:]})
	utils.ExitOnError(err, "Error processing command line arguments")

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
