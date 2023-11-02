package cmd

import (
	"fmt"
	"os"
)

type formatter struct {
	absoluteFilePaths []string
	localPrefixes     []string
	write             bool
	diff              bool
	verbose           bool
}

// run runs the formatter.
func (f formatter) run() error {
	for _, absoluteFilePath := range f.absoluteFilePaths {
		if f.verbose {
			fmt.Printf("Formatting %q...\n", absoluteFilePath)
		}

		formattedFile, err := formatFile(f.localPrefixes, absoluteFilePath)
		if err != nil {
			return fmt.Errorf("failed to format file %q: %w", absoluteFilePath, err)
		}

		if f.write {
			if err := os.WriteFile(absoluteFilePath, []byte(formattedFile), 0o644); err != nil {
				return fmt.Errorf("failed to write file %q: %w", absoluteFilePath, err)
			}
		}
		if f.diff {
			if err := printFileDiff(absoluteFilePath, formattedFile); err != nil {
				return fmt.Errorf("failed to diff file %q: %w", absoluteFilePath, err)
			}
		}
	}

	return nil
}
