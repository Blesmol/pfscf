package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/pdf"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	cmdFillUseExampleValues bool
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
	fillCmd.Flags().BoolVarP(&cfg.Global.DrawGrid, "grid", "g", false, "Draw a coordinate grid on the output file")
	fillCmd.Flags().BoolVarP(&cfg.Global.DrawCellBorder, "cellBorder", "c", false, "Draw the cell borders of all added fields")
	fillCmd.Flags().BoolVarP(&cmdFillUseExampleValues, "exampleValues", "e", false, "Use example values to fill out the chronicle")

	return fillCmd
}

func executeFill(cmd *cobra.Command, cmdArgs []string) {
	utils.Assert(len(cmdArgs) >= 3, "Number of arguments should be guaranteed by cobra settings")

	tmplName := cmdArgs[0]
	inFile := cmdArgs[1]
	outFile := cmdArgs[2]

	if inFile == outFile {
		utils.ExitWithMessage("Input file and output file must not be identical")
	}

	cTmpl, err := template.Get(tmplName)
	utils.ExitOnError(err, "Error getting template")

	// parse remaining arguments
	var argStore *args.ArgStore
	if !cmdFillUseExampleValues {
		argStore, err = args.NewArgStore(args.ArgStoreInit{Args: cmdArgs[3:]})
	} else {
		argStore, err = args.NewArgStore(args.ArgStoreInit{Args: cTmpl.GetExampleArguments()})
	}
	utils.ExitOnError(err, "Eror processing command line arguments")

	pf, err := pdf.NewFile(inFile)
	utils.ExitOnError(err, "Error opening input file '%v'", inFile)

	err = pf.Fill(argStore, cTmpl, outFile)
	utils.ExitOnError(err, "Error when filling out chronicle")
}
