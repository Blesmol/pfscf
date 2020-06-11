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

	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)

	// Configuration test run
	config := GetGlobalConfig()
	fmt.Printf("Config:\n%+v\n", *config)
	fmt.Printf("Content: %+v", *config.Content)

	/*
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
	*/
}
