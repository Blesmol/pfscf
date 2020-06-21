package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	fillCmd        *cobra.Command
	drawGrid       bool
	drawCellBorder bool
)

// GetFillCommand returns the cobra command for the "fill" action.
func GetFillCommand() (cmd *cobra.Command) {
	Assert(fillCmd == nil, "FillCmd already initialized")

	fillCmd = &cobra.Command{
		Use:   "fill <config> <infile> <outfile>",
		Short: "Fill a single chronicle sheet",
		Long:  "Fill a single chronicle sheet with parameters provided on the command line.",

		Args: cobra.MinimumNArgs(3),

		Run: executeFill,
	}
	fillCmd.Flags().BoolVarP(&drawGrid, "grid", "g", false, "Draw a coordinate grid on the output file")
	fillCmd.Flags().BoolVarP(&drawCellBorder, "cellBorder", "c", false, "Draw the cell borders of all added fields")

	return fillCmd
}

func executeFill(cmd *cobra.Command, args []string) {
	Assert(len(args) >= 3, "Number of arguments should be guaranteed by cobra settings")

	cfgName := args[0]
	inFile := args[1]
	outFile := args[2]

	Assert(inFile != outFile, "Input file and output file must not be identical")

	yCfg, err := GetConfigByName(cfgName)
	AssertNoError(err) // TODO proper error message and exit

	cCfg := yCfg.GetChronicleConfig() // TODO assign to something and work with it

	// parse remaining arguments
	as := ParseArgs(args[3:])

	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	pdf := NewPdf(inFile)

	// extract chronicle page from pdf
	extractedPage := pdf.ExtractPage(-1, workDir)
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := NewStamp(width, height)

	if drawCellBorder {
		stamp.SetCellBorder(true)
	}

	// add content to stamp
	for key, value := range as {
		//fmt.Printf("Processing Key='%v', value='%v'\n", key, *value)

		content, exists := cCfg.GetContent(key)
		Assert(exists, "No content with key="+key)
		stamp.AddContent(content, value)
	}

	if drawGrid {
		stamp.CreateMeasurementCoordinates(25, 5)
	}

	// write stamp
	stampFile := filepath.Join(workDir, "stamp.pdf")
	stamp.WriteToFile(stampFile)

	// add watermark/stamp to page
	extractedPage.StampIt(stampFile, outFile)

}
