package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/cmd"
	"github.com/Blesmol/pfscf/pfscf/web"
)

var (
	version = "dev"
)

func main() {

	RootCmd := &cobra.Command{
		Use:   "pfscf",
		Short: "The Pathfinder Society Chronicle Filler (v" + version + ")",
	}

	RootCmd.PersistentFlags().BoolVarP(&cfg.Global.Verbose, "verbose", "v", false, "verbose output")

	RootCmd.AddCommand(cmd.GetFillCommand())
	RootCmd.AddCommand(cmd.GetTemplateCommand())
	RootCmd.AddCommand(cmd.GetBatchCommand())
	RootCmd.AddCommand(cmd.GetOpenCommand())

	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// WebCreateExample takes only a template name and some filenames and provides some PoC functionality.
func WebCreateExample(template, inFile, outFile string) {
	web.CreateExample(template, inFile, outFile)
}

// WebCreateExampleWithoutArgs is like WebCreateExample, but with hardcoded values.
func WebCreateExampleWithoutArgs() {
	web.CreateExample("pfs2.s1-14", "input.pdf", "output.pdf")
}
