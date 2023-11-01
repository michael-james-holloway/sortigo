package main

import (
	"github.com/michael-james-holloway/sortigo/cmd"
	"github.com/spf13/cobra"
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
