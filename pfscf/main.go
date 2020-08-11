package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Blesmol/pfscf/pfscf/cfg"
	"github.com/Blesmol/pfscf/pfscf/cmd"
)

var (
	version = "dev"
)

func main() {

	RootCmd := &cobra.Command{
		Use:   "pfscf",
		Short: "The Pathfinder Society Chronicle Tagger v" + version,
	}

	RootCmd.PersistentFlags().BoolVarP(&cfg.Global.Verbose, "verbose", "v", false, "verbose output")

	RootCmd.AddCommand(cmd.GetFillCommand())
	RootCmd.AddCommand(cmd.GetTemplateCommand())
	RootCmd.AddCommand(cmd.GetBatchCommand())

	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
