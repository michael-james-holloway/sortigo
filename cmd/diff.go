package cmd

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// printFileDiff prints the diff of the file iff the file is different.
//
// It first logs the absolute file path and then the diff.
func printFileDiff(originalFile, formattedFile []byte) error {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(originalFile), string(formattedFile), false)
	fmt.Println(prettyGitDiff(diffs))

	return nil
}

func prettyGitDiff(diffs []diffmatchpatch.Diff) string {
	var result strings.Builder
	lines := strings.Split(strings.Join(diffTexts(diffs), ""), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "++") || strings.HasPrefix(line, "--") {
			result.WriteString(fmt.Sprintf("   %4d %s\n", i+1, line))
		} else {
			result.WriteString(fmt.Sprintf("   %4d %s\n", i+1, line))
		}
	}

	return result.String()
}

func diffTexts(diffs []diffmatchpatch.Diff) []string {
	var result []string
	for _, diff := range diffs {
		fmt.Println(diff)
		text := diff.Text
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if i < len(lines)-1 {
				line += "\n"
			}
			var newLine string
			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				newLine = fmt.Sprintf("++%s", line)
			case diffmatchpatch.DiffDelete:
				newLine = fmt.Sprintf("--%s", line)
			case diffmatchpatch.DiffEqual:
				newLine = fmt.Sprintf("  %s", line)
			}

			result = append(result, newLine)
		}
	}

	return result
}
