package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Blesmol/pfscf/pfscf/args"
	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/csv"
	"github.com/Blesmol/pfscf/pfscf/pdf"
	"github.com/Blesmol/pfscf/pfscf/template"
	"github.com/Blesmol/pfscf/pfscf/utils"
)

const (
	namingPlaceholderPattern = `<(\w+)>`
)

var (
	actionBatchCreateUsageExampleValues  bool
	actionBatchCreateSeparator           string
	actionBatchCreateSuppressOpenOutfile bool

	actionBatchOutputPattern  string
	actionBatchTemplate       string
	actionBatchInputChronicle string
	actionBatchOutputDir      string

	regexNamingPlaceholder = regexp.MustCompile(namingPlaceholderPattern)
)

// GetBatchCommand returns the cobra command for the "batch" action.
func GetBatchCommand() (cmd *cobra.Command) {
	cmdBatch := &cobra.Command{
		Use:     "batch",
		Aliases: []string{"b"},

		Short: "Fill out multiple chronicles in one go",
		Long:  "The batch operation can fill out multiple chronicles in one go by reading all necessary input from a csv file.",

		Args: cobra.ExactArgs(0),
	}

	cmdBatch.PersistentFlags().StringVarP(&actionBatchOutputPattern, "output-pattern", "p", "Chronicle_<char>_<societyid>.pdf", "Naming pattern for the generated chronicle files")
	cmdBatch.PersistentFlags().StringVarP(&actionBatchTemplate, "template", "t", "", "Name of the template to use, e.g. pfs2.s1-22")
	cmdBatch.PersistentFlags().StringVarP(&actionBatchInputChronicle, "input-chronicle", "i", "", "Filename of the empty input scenario chronicle")
	cmdBatch.PersistentFlags().StringVarP(&actionBatchOutputDir, "output-dir", "o", ".", "Directory in which the generated chronicles should be stored")

	cmdCreate := &cobra.Command{
		Use:     "create <csv_file> [<content_id>=<value> ...]",
		Aliases: []string{"c"},

		Short: "Create ready-to-fill csv file based on selected template",
		//Long:  "TBD",

		Args: cobra.MinimumNArgs(1),

		Run: executeBatchCreate,
	}
	cmdCreate.Flags().BoolVarP(&actionBatchCreateUsageExampleValues, "examples", "e", false, "Use example values to fill out the chronicle")
	cmdCreate.Flags().StringVarP(&actionBatchCreateSeparator, "separator", "s", ";", "Field separator character for resulting CSV file")
	cmdCreate.Flags().BoolVarP(&actionBatchCreateSuppressOpenOutfile, "no-auto-open", "n", false, "Suppress auto-opening the created CSV file")

	cmdBatch.AddCommand(cmdCreate)

	cmdFill := &cobra.Command{
		Use:     "fill <csv_file> [<param_id>=<value> ...]",
		Aliases: []string{"f"},

		Short: "Fill multiple templates with values read from a csv file",

		Args: cobra.ExactArgs(1),

		Run: executeBatchFill,
	}
	cmdFill.Flags().Float64VarP(&cfg.Global.OffsetX, "offset-x", "x", 0, "Assume an additional offset for the X axis of the chronicle")
	cmdFill.Flags().Float64VarP(&cfg.Global.OffsetY, "offset-y", "y", 0, "Assume an additional offset for the Y axis of the chronicle")

	cmdBatch.AddCommand(cmdFill)

	return cmdBatch
}

func executeBatchCreate(cmd *cobra.Command, cmdArgs []string) {
	utils.Assert(len(cmdArgs) >= 1, "Number of arguments should be guaranteed by cobra settings")

	outFile := cmdArgs[0]
	remainingArgs := cmdArgs[1:]

	tmplName := getFlagOrExit(cmd, "template")

	warnOnWrongFileExtension(outFile, "csv")

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

	ts, err := template.GetStore()
	utils.ExitOnError(err, "Error retrieving templates")
	cTmpl, exists := ts.Get(tmplName)
	if !exists {
		utils.ExitWithMessage("Template '%v' not found", tmplName)
	}

	// parse remaining arguments
	var argStore *args.Store
	if !actionBatchCreateUsageExampleValues {
		argStore, err = args.NewStore(args.StoreInit{Args: remainingArgs})
	} else {
		argStore, err = args.NewStore(args.StoreInit{Args: cTmpl.GetExampleArguments()})
	}
	utils.ExitOnError(err, "Error processing command line arguments")

	// prepare cmd line flags to be included in csv file.
	// we only check flags which were defined on the parent command to not have to keep
	// flags synchronous between different subcommands
	cmdFlags := make([][]string, 0)
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if f.Value.Type() == "string" {
			fName := fmt.Sprintf("flag:--%v", f.Name)
			fVal := f.Value.String()
			if !utils.IsSet(fVal) {
				fName = "#" + fName
			}
			cmdFlags = append(cmdFlags, []string{fName, fVal})
		}
	})

	err = cTmpl.GenerateCsvFile(outFile, separator, argStore, cmdFlags)
	utils.ExitOnError(err, "Error writing CSV file for template %v", tmplName)

	if !actionBatchCreateSuppressOpenOutfile {
		fmt.Printf("Trying to open file '%v' in standard viewer\n", outFile)
		err = utils.OpenWithDefaultViewer(outFile)
		utils.ExitOnError(err, "Error opening file")
	}
}

func executeBatchFill(cmd *cobra.Command, cmdArgs []string) {
	utils.Assert(len(cmdArgs) >= 1, "Number of arguments should be guaranteed by cobra settings")

	inCsv := cmdArgs[0]
	remainingArgs := cmdArgs[1:]
	warnOnWrongFileExtension(inCsv, "csv")

	csvRecords, err := csv.ReadCsvFile(inCsv)
	utils.ExitOnError(err, "Cannot read CSV file")

	err = setFlagsFromRecords(cmd, csvRecords)
	utils.ExitOnError(err, "Error parsing CSV file '%v'", inCsv)

	tmplName := getFlagOrExit(cmd, "template")
	outDir := getFlagOrExit(cmd, "output-dir")
	inPdf := getFlagOrExit(cmd, "input-chronicle")
	warnOnWrongFileExtension(inPdf, "pdf")

	// get templates
	ts, err := template.GetStore()
	utils.ExitOnError(err, "Error retrieving templates")
	cTmpl, exists := ts.Get(tmplName)
	if !exists {
		utils.ExitWithMessage("Template '%v' not found", tmplName)
	}

	// get arg value stores from CSV data
	batchArgStores, err := args.GetArgStoresFromCsvRecords(csvRecords)
	utils.ExitOnError(err, "Error parsing CSV file")
	if len(batchArgStores) == 0 {
		utils.ExitWithMessage("No output files were created as CSV file '%v' does not contain any player values", inCsv)
	}

	// parse remaining command line arguments
	cmdLineArgStore, err := args.NewStore(args.StoreInit{Args: remainingArgs})
	utils.ExitOnError(err, "Error processing command line arguments")

	// ensure output directory exists
	err = os.MkdirAll(outDir, os.ModePerm)
	utils.ExitOnError(err, "Error creating output directory")

	for idx, batchArgStore := range batchArgStores {
		cmdLineArgStore.SetParent(batchArgStore) // command line arguments have priority

		pdf, err := pdf.NewFile(inPdf)
		utils.ExitOnError(err, "Error opening input file '%v'", inPdf)

		playerNumber := idx + 1
		baseOutfile, err := getOutputFilenameForPlayer(actionBatchOutputPattern, cmdLineArgStore)
		utils.ExitOnError(err, "Error getting output filename")
		outfile := filepath.Join(outDir, baseOutfile)

		fmt.Printf("Creating file %v\n", outfile)
		err = pdf.Fill(cmdLineArgStore, cTmpl, outfile)
		utils.ExitOnError(err, "Error when filling out chronicle for player %d", playerNumber)
	}
}

func getOutputFilenameForPlayer(pattern string, as *args.Store) (outfile string, err error) {
	if !utils.IsSet(pattern) {
		return "", fmt.Errorf("No naming pattern for output file provided")
	}
	outfile = pattern

	matches := regexNamingPlaceholder.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		utils.Assert(len(match) == 2, "Should have exactly 2 elements: the complete matching string and the subgroup")

		// lookup requested key in argStore
		resolvedValue, exists := as.Get(match[1])
		if !exists {
			return "", fmt.Errorf("Cannot find value for filename placeholder %v", match[0])
		}

		// replace key in filename with key from argStore
		re := regexp.MustCompile(match[0])
		outfile = re.ReplaceAllString(outfile, resolvedValue)
	}

	return outfile, nil
}

func getFlagOrExit(cmd *cobra.Command, flagName string) string {
	flag := cmd.Flags().Lookup(flagName)
	if flag == nil || !utils.IsSet(flag.Value.String()) {
		fmt.Fprintf(os.Stderr, "Error: required flag \"%v\" not set\n%v", flagName, cmd.UsageString())
		os.Exit(1)
	}
	return flag.Value.String()
}
