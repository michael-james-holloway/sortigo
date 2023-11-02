package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

type importPathAndSpec struct {
	importPath string
	spec       *ast.ImportSpec
}

func (f formatter) formatFile(filePath string) ([]byte, []byte, error) {
	// Create a new scanner and parse the file into an AST
	fileSet := token.NewFileSet()

	// Original file data.
	originalFileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}
	node, err := parser.ParseFile(fileSet, filePath, originalFileData, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Sort imports.
	f.sortImports(node)

	// Convert the modified AST back to a string
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, node); err != nil {
		return nil, nil, fmt.Errorf("failed to format AST: %w", err)
	}

	// Return the formatted string
	return originalFileData, buf.Bytes(), nil
}

func (f formatter) sortImports(node *ast.File) {
	for i, decl := range node.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok || decl.Tok != token.IMPORT {
			continue
		}

		// Sort the declaration.
		f.sortDeclNode(decl)

		// Set the declaration.
		node.Decls[i] = decl
	}
}

func (f formatter) sortDeclNode(decl *ast.GenDecl) {
	// Aggregate imports into groups.
	stdlib := make([]importPathAndSpec, 0, len(decl.Specs))
	thirdParty := make([]importPathAndSpec, 0, len(decl.Specs))
	firstParty := make([]importPathAndSpec, 0, len(decl.Specs))
	for _, spec := range decl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		imp := importPathAndSpec{
			importPath: strings.Trim(importSpec.Path.Value, "\""),
			spec:       importSpec,
		}

		// If we are consolidating multiple same-group import blocks into the same block, then
		// we need to set the value position to 0 so that the formatter doesn't try to preserve
		// the original, separate block formatting.
		if !f.dontConsolidateBlocks {
			imp.spec.Path.ValuePos = 0
		}

		switch {
		case isFirstParty(f.localPrefixes, imp.importPath):
			firstParty = append(firstParty, imp)
		case strings.Contains(imp.importPath, "."):
			thirdParty = append(thirdParty, imp)
		default:
			stdlib = append(stdlib, imp)
		}
	}

	// Sort each group.
	sort.Slice(stdlib, func(i, j int) bool {
		return stdlib[i].importPath < stdlib[j].importPath
	})
	sort.Slice(thirdParty, func(i, j int) bool {
		return thirdParty[i].importPath < thirdParty[j].importPath
	})
	sort.Slice(firstParty, func(i, j int) bool {
		return firstParty[i].importPath < firstParty[j].importPath
	})

	// Clear original import specs.
	decl.Specs = []ast.Spec{}

	// Append sorted imports back to the declaration.
	for _, imp := range stdlib {
		decl.Specs = append(decl.Specs, imp.spec)
	}
	if len(thirdParty) > 0 && len(stdlib) > 0 {
		decl.Specs = append(decl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{Value: ""},
		})
	}
	for _, imp := range thirdParty {
		decl.Specs = append(decl.Specs, imp.spec)
	}
	if len(firstParty) > 0 && (len(stdlib) > 0 || len(thirdParty) > 0) {
		decl.Specs = append(decl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{Value: ""},
		})
	}
	for _, imp := range firstParty {
		decl.Specs = append(decl.Specs, imp.spec)
	}
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
