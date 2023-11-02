package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// must is a helper function that panics if err is not nil.
func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

// findGoFilesInDir returns a list of all go files in the directory.
func findGoFilesInDir(dir string) ([]string, error) {
	absoluteGoSourceFilePaths := make([]string, 0)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		absoluteGoSourceFilePaths = append(absoluteGoSourceFilePaths, must(filepath.Abs(path)))

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk dir: %w", err)
	}

	return absoluteGoSourceFilePaths, nil
}
