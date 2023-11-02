package main

import (
	"github.com/spf13/cobra"

	"github.com/michael-james-holloway/sortigo/cmd"
)

var rootCMD = &cobra.Command{
	Use: "sortigo",
}

func init() {
	rootCMD.AddCommand(cmd.FormatCMD)
}

func main() {
	if err := rootCMD.Execute(); err != nil {
		panic("failed to execute root command")
	}
}
