package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

func formatFile(localPrefixes []string, filePath string) (string, error) {
	// Create a new scanner and parse the file into an AST
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	// Sort imports.
	for i, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		// Aggregate imports into groups.
		stdlib := make([]string, 0, len(genDecl.Specs))
		thirdParty := make([]string, 0, len(genDecl.Specs))
		firstParty := make([]string, 0, len(genDecl.Specs))
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}

			// Add the import to the appropriate group.
			importPath := strings.Trim(importSpec.Path.Value, "\"")
			switch {
			case isFirstParty(localPrefixes, importPath):
				firstParty = append(firstParty, importPath)
			case strings.Contains(importPath, "."):
				thirdParty = append(thirdParty, importPath)
			default:
				stdlib = append(stdlib, importPath)
			}
		}

		// Sort each group.
		sort.Strings(stdlib)
		sort.Strings(thirdParty)
		sort.Strings(firstParty)

		// Clear original import specs.
		genDecl.Specs = []ast.Spec{}

		// Append sorted imports back to the declaration.
		for _, path := range stdlib {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Value: fmt.Sprintf("\"%s\"", path)},
			})
		}
		if len(thirdParty) > 0 && len(stdlib) > 0 {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Value: ""},
			})
		}
		for _, path := range thirdParty {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Value: fmt.Sprintf("\"%s\"", path)},
			})
		}
		if len(firstParty) > 0 && (len(stdlib) > 0 || len(thirdParty) > 0) {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Value: ""},
			})
		}
		for _, path := range firstParty {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Value: fmt.Sprintf("\"%s\"", path)},
			})
		}

		// Set the declaration.
		node.Decls[i] = genDecl
	}

	// Convert the modified AST back to a string
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, node); err != nil {
		return "", fmt.Errorf("failed to format AST: %w", err)
	}

	// Return the formatted string
	return buf.String(), nil
}

// isFirstParty returns true if the import path is a first party import.
func isFirstParty(localPrefixes []string, importPath string) bool {
	for _, localPrefix := range localPrefixes {
		if strings.HasPrefix(importPath, localPrefix) {
			return true
		}
	}

	return false
}
