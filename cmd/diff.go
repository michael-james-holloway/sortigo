package cmd

import (
	"fmt"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// printFileDiff prints the diff of the file iff the file is different.
//
// It first logs the absolute file path and then the diff.
func printFileDiff(absoluteFilePath string, formattedFileContents string) error {
	currentFileContents, err := os.ReadFile(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file at path %s: %w", absoluteFilePath, err)
	}
	currentFileContentsStr := string(currentFileContents)

	if currentFileContentsStr == formattedFileContents {
		return nil
	}

	// Calculate and log the diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(currentFileContentsStr, formattedFileContents, false)
	fmt.Println("File path:", absoluteFilePath)
	fmt.Println(dmp.DiffPrettyText(diffs))

	return nil
}
