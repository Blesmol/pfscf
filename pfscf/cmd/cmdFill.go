package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/pdf"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	cmdFillUseExampleValues    bool
	cmdFillSuppressOpenOutfile bool
)

// GetFillCommand returns the cobra command for the "fill" action.
func GetFillCommand() (cmd *cobra.Command) {
	fillCmd := &cobra.Command{
		Use:     "fill <template> <infile> <outfile> [<param_id>=<value> ...]",
		Aliases: []string{"f"},

		Short: "Fill out a single chronicle sheet",
		Long:  "Fill out a single chronicle sheet with parameters provided on the command line.",

		Args: cobra.MinimumNArgs(3),

		Run: executeFill,
	}
	fillCmd.Flags().StringVarP(&cfg.Global.DrawCanvasGrid, "canvas-grid", "g", "", "Draw a coordinate grid in the output file for the canvas with the provided name")
	fillCmd.Flags().BoolVarP(&cfg.Global.DrawCellBorder, "cell-border", "c", false, "Draw the cell borders of all added fields")
	fillCmd.Flags().BoolVarP(&cmdFillUseExampleValues, "examples", "e", false, "Use example values to fill out the chronicle")
	fillCmd.Flags().BoolVarP(&cmdFillSuppressOpenOutfile, "no-auto-open", "n", false, "Suppress auto-opening the filled out chronicle")
	fillCmd.Flags().BoolVarP(&cfg.Global.DrawCanvas, "draw-canvas", "d", false, "Draw a border around all defined canvases")
	fillCmd.Flags().Float64VarP(&cfg.Global.OffsetX, "offset-x", "x", 0, "Assume an additional offset for the X axis of the chronicle")
	fillCmd.Flags().Float64VarP(&cfg.Global.OffsetY, "offset-y", "y", 0, "Assume an additional offset for the Y axis of the chronicle")

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

	warnOnWrongFileExtension(inFile, "pdf")
	warnOnWrongFileExtension(outFile, "pdf")

	ts, err := template.GetStore()
	utils.ExitOnError(err, "Error retrieving templates")
	cTmpl, exists := ts.Get(tmplName)
	if !exists {
		utils.ExitWithMessage("Template '%v' not found", tmplName)
	}

	// parse remaining arguments
	var argStore *args.Store
	if !cmdFillUseExampleValues {
		argStore, err = args.NewStore(args.StoreInit{Args: cmdArgs[3:]})
	} else {
		argStore, err = args.NewStore(args.StoreInit{Args: cTmpl.GetExampleArguments()})
	}
	utils.ExitOnError(err, "Eror processing command line arguments")

	pf, err := pdf.NewFile(inFile)
	utils.ExitOnError(err, "Error opening input file '%v'", inFile)

	err = pf.Fill(argStore, cTmpl, outFile)
	utils.ExitOnError(err, "Error when filling out chronicle")

	if !cmdFillSuppressOpenOutfile {
		fmt.Printf("Trying to open file '%v' in standard PDF viewer\n", outFile)
		err = utils.OpenWithDefaultViewer(outFile)
		utils.ExitOnError(err, "Error opening PDF file")
	}
}
