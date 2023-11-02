package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	localPrefixesFlagName         = "local-prefixes"
	dontConsolidateBlocksFlagName = "dont-consolidate-blocks"
	writeFlagName                 = "write"
	checkFlagName                 = "check"
	verboseFlagName               = "verbose"
)

func init() {
	FormatCMD.Flags().StringSliceP(
		localPrefixesFlagName,
		"l",
		[]string{},
		"Local prefix(es) to consider first party imports (e.g. github.com/michael-james-holloway/sortigo).",
	)
	FormatCMD.Flags().Bool(
		dontConsolidateBlocksFlagName,
		false,
		"Don't consolidate existing separate blocks of the same group type (e.g. multiple third party blocks).",
	)
	FormatCMD.Flags().BoolP(
		writeFlagName,
		"w",
		false,
		"Write the formatted file back to the original file.",
	)
	FormatCMD.Flags().BoolP(
		checkFlagName,
		"c",
		false,
		"Checks if any files are different, and exits with a non-zero exit code if so.",
	)
	FormatCMD.Flags().BoolP(
		verboseFlagName,
		"v",
		false,
		"Verbose output.",
	)
}

var FormatCMD = &cobra.Command{
	Use:  "format",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument, got %d", len(args))
		}

		if err := cmd.ParseFlags(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		localPrefixes := must(cmd.Flags().GetStringSlice(localPrefixesFlagName))
		dontConsolidateBlocks := must(cmd.Flags().GetBool(dontConsolidateBlocksFlagName))
		write := must(cmd.Flags().GetBool(writeFlagName))
		check := must(cmd.Flags().GetBool(checkFlagName))
		verbose := must(cmd.Flags().GetBool(verboseFlagName))

		// Validate flag values.
		if len(localPrefixes) == 0 {
			return fmt.Errorf("no local prefixes passed")
		}
		if write && check {
			return fmt.Errorf("cannot set both %q and %q", writeFlagName, checkFlagName)
		}

		fileOrDirToFormat := args[0]
		if fileOrDirToFormat == "" {
			return fmt.Errorf("no file or directory to format")
		}

		var absoluteFilePathsToFormat []string
		if strings.HasSuffix(fileOrDirToFormat, ".go") {
			absoluteFilePathsToFormat = []string{must(filepath.Abs(fileOrDirToFormat))}
		} else {
			var err error
			absoluteFilePathsToFormat, err = findGoFilesInDir(fileOrDirToFormat)
			if err != nil {
				return fmt.Errorf("failed to get go files in dir: %w", err)
			}
		}

		if err := (formatter{
			absoluteFilePaths:     absoluteFilePathsToFormat,
			localPrefixes:         localPrefixes,
			dontConsolidateBlocks: dontConsolidateBlocks,
			write:                 write,
			check:                 check,
			verbose:               verbose,
		}).run(); err != nil {
			return fmt.Errorf("failed to run formatter: %w", err)
		}

		return nil
	},
}
