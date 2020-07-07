package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	drawGrid       bool
	drawCellBorder bool
	useDummyValues bool
)

// GetFillCommand returns the cobra command for the "fill" action.
func GetFillCommand() (cmd *cobra.Command) {
	fillCmd := &cobra.Command{
		Use:     "fill <template> <infile> <outfile>",
		Aliases: []string{"f"},

		Short: "Fill out a single chronicle sheet",
		Long:  "Fill out a single chronicle sheet with parameters provided on the command line.",

		Args: cobra.MinimumNArgs(3),

		Run: executeFill,
	}
	fillCmd.Flags().BoolVarP(&drawGrid, "grid", "g", false, "Draw a coordinate grid on the output file")
	fillCmd.Flags().BoolVarP(&drawCellBorder, "cellBorder", "c", false, "Draw the cell borders of all added fields")
	fillCmd.Flags().BoolVarP(&useDummyValues, "dummyValues", "d", false, "Use dummy values to fill out the chronicle.")

	return fillCmd
}

func executeFill(cmd *cobra.Command, args []string) {
	Assert(len(args) >= 3, "Number of arguments should be guaranteed by cobra settings")

	tmplName := args[0]
	inFile := args[1]
	outFile := args[2]

	Assert(inFile != outFile, "Input file and output file must not be identical")

	cTmpl, err := GetTemplate(tmplName)
	ExitOnError(err, "Error getting template")

	// parse remaining arguments
	var argStore ArgStore
	if !useDummyValues {
		argStore = ArgStoreFromArgs(args[3:])
	} else {
		argStore = ArgStoreFromTemplateExamples(cTmpl)
	}

	// prepare temporary working dir
	workDir := GetTempDir()
	defer os.RemoveAll(workDir)

	pdf, err := NewPdf(inFile)
	AssertNoError(err)

	// extract chronicle page from pdf
	extractedPage, err := pdf.ExtractPage(-1, workDir)
	ExitOnError(err, "Error extracing page from %v", pdf.filename)
	width, height := extractedPage.GetDimensionsInPoints()

	// create stamp
	stamp := NewStamp(width, height)

	if drawCellBorder {
		stamp.SetCellBorder(true)
	}

	// add content to stamp
	for key, value := range argStore {
		//fmt.Printf("Processing Key='%v', value='%v'\n", key, *value)

		content, exists := cTmpl.GetContent(key)
		Assert(exists, "No content with key="+key)
		err := stamp.AddContent(content, &value)
		ExitOnError(err, "Eror adding content: %w", err)
	}

	if drawGrid {
		stamp.CreateMeasurementCoordinates(25, 5)
	}

	// write stamp
	stampFile := filepath.Join(workDir, "stamp.pdf")
	err = stamp.WriteToFile(stampFile)
	ExitOnError(err, "Error while writing to file %v", stampFile)

	// add watermark/stamp to page
	extractedPage.StampIt(stampFile, outFile)
}
