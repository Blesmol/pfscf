package web

import (
	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/pdf"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

// CreateExample is a simple PoC that takes the names of a blank
// chronicle PDF as input file and will produce an output file for PFS2.
func CreateExample(tmplName, inFile, outFile string) {
	if inFile == outFile {
		utils.ExitWithMessage("Input file and output file must not be identical")
	}

	ts, err := template.GetStore()
	utils.ExitOnError(err, "Error retrieving templates")
	cTmpl, exists := ts.Get(tmplName)
	if !exists {
		utils.ExitWithMessage("Template '%v' not found", tmplName)
	}

	// parse remaining arguments
	argStore, err := args.NewStore(args.StoreInit{Args: cTmpl.GetExampleArguments()})
	utils.ExitOnError(err, "Eror processing command line arguments")

	pf, err := pdf.NewFile(inFile)
	utils.ExitOnError(err, "Error opening input file '%v'", inFile)

	err = pf.Fill(argStore, cTmpl, outFile)
	utils.ExitOnError(err, "Error when filling out chronicle")
}
