package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {

	RootCmd := &cobra.Command{
		Use:   "pfsct",
		Short: "The Pathfinder Society Chronicle Tagger v" + version,
	}

	RootCmd.AddCommand(GetFillCommand())
	RootCmd.AddCommand(GetTemplateCommand())

	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
