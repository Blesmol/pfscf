package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var fillCmd *cobra.Command

// GetFillCommand returns the cobra command for the "fill" action.
func GetFillCommand() (cmd *cobra.Command) {
	Assert(fillCmd == nil, "FillCmd already initialized")

	fillCmd = &cobra.Command{
		Use:   "fill [INFILE] [OUTFILE]",
		Short: "Fill a single chronicle sheet",
		Long:  "Fill a single chronicle sheet with parameters provided on the command line.",

		Args: cobra.MinimumNArgs(2),

		Run: executeFill,
	}
	//fillCmd.Flags().BoolVar(&myFlag, "varBool", false, "Do something with a bool flag.")

	return fillCmd
}

func executeFill(cmd *cobra.Command, args []string) {
	Assert(len(args) >= 2, "Number of arguments should be guaranteed by cobra settings")

	inFile := args[0]
	outFile := args[1]

	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	pdf := NewPdf(inFile)

	// extract chronicle page from pdf
	extractedPage := pdf.ExtractPage(-1, workDir)
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := NewStamp(width, height)

	// demo text
	stamp.AddText(433, 107, "Grand Archive", "Helvetica", 8)
	stamp.AddCellText(227, 730, 305, 716, "05.06.2020", "Helvetica", 14)

	//stamp.CreateMeasurementCoordinates(25, 5)

	// write stamp
	stampFile := filepath.Join(workDir, "stamp.pdf")
	stamp.WriteToFile(stampFile)

	// add watermark/stamp to page
	extractedPage.StampIt(stampFile, outFile)

}
