package main

import (
	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	drawGrid                   bool
	drawCellBorder             bool
	actionFillUseExampleValues bool
)

// GetFillCommand returns the cobra command for the "fill" action.
func GetFillCommand() (cmd *cobra.Command) {
	fillCmd := &cobra.Command{
		Use:     "fill <template> <infile> <outfile> [<content_id>=<value> ...]",
		Aliases: []string{"f"},

		Short: "Fill out a single chronicle sheet",
		Long:  "Fill out a single chronicle sheet with parameters provided on the command line.",

		Args: cobra.MinimumNArgs(3),

		Run: executeFill,
	}
	fillCmd.Flags().BoolVarP(&drawGrid, "grid", "g", false, "Draw a coordinate grid on the output file")
	fillCmd.Flags().BoolVarP(&drawCellBorder, "cellBorder", "c", false, "Draw the cell borders of all added fields")
	fillCmd.Flags().BoolVarP(&actionFillUseExampleValues, "exampleValues", "e", false, "Use example values to fill out the chronicle")

	return fillCmd
}

func executeFill(cmd *cobra.Command, args []string) {
	utils.Assert(len(args) >= 3, "Number of arguments should be guaranteed by cobra settings")

	tmplName := args[0]
	inFile := args[1]
	outFile := args[2]

	if inFile == outFile {
		utils.ExitWithMessage("Input file and output file must not be identical")
	}

	cTmpl, err := GetTemplate(tmplName)
	utils.ExitOnError(err, "Error getting template")

	// parse remaining arguments
	var argStore *ArgStore
	if !actionFillUseExampleValues {
		argStore, err = NewArgStore(ArgStoreInit{args: args[3:]})
	} else {
		argStore, err = NewArgStore(ArgStoreInit{args: cTmpl.GetExampleArguments()})
	}
	utils.ExitOnError(err, "Eror processing command line arguments")

	pdf, err := NewPdf(inFile)
	utils.ExitOnError(err, "Error opening input file '%v'", inFile)

	err = pdf.Fill(argStore, cTmpl, outFile)
	utils.ExitOnError(err, "Error when filling out chronicle")
}
