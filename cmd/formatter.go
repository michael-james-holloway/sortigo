package cmd

import (
	"bytes"
	"fmt"
	"os"
)

type formatter struct {
	absoluteFilePaths     []string
	localPrefixes         []string
	dontConsolidateBlocks bool
	write                 bool
	verbose               bool
}

// run runs the formatter.
func (f formatter) run() error {
	diffFilePaths := make([]string, 0, len(f.absoluteFilePaths))
	for _, absoluteFilePath := range f.absoluteFilePaths {
		isDifferent, err := f.runOneFile(absoluteFilePath)
		if err != nil {
			return fmt.Errorf("failed to run formatter on file %q: %w", absoluteFilePath, err)
		}
		if isDifferent {
			diffFilePaths = append(diffFilePaths, absoluteFilePath)
		}
	}

	// If we have diffs, and we weren't in write mode, print the diffs and raise an error.
	if len(diffFilePaths) > 0 && !f.write {
		fmt.Printf("💥 Oh no! Diffs were found in the following files:\n")
		for _, diffFilePath := range diffFilePaths {
			fmt.Printf("  %s\n", diffFilePath)
		}
		fmt.Printf("\n")

		os.Exit(1)
	}
	if len(diffFilePaths) == 0 {
		fmt.Printf("🎉 No diffs were found!\n")
	}

	return nil
}

// runOneFile runs the formatter on one file.
func (f formatter) runOneFile(absoluteFilePath string) (bool, error) {
	if f.verbose {
		fmt.Printf("Formatting %q...\n", absoluteFilePath)
	}

	originalFile, formattedFile, err := f.formatFile(absoluteFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to format file %q: %w", absoluteFilePath, err)
	}

	isDifferent := bytes.Compare(originalFile, formattedFile) != 0

	if isDifferent && f.write {
		if err := os.WriteFile(absoluteFilePath, formattedFile, 0o644); err != nil {
			return false, fmt.Errorf("failed to write file %q: %w", absoluteFilePath, err)
		}
		if f.verbose {
			fmt.Printf("Wrote %q.\n", absoluteFilePath)
		}
	}
	if isDifferent && !f.write {
		// TODO(holloway): Re-enable once this looks better. Also gate on a flag.
		// printFileDiff(originalFile, formattedFile)
	}

	return isDifferent, nil
}
