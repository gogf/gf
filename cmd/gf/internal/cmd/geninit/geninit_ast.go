// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// ASTReplacer handles import path replacement using Go AST
type ASTReplacer struct {
	oldModule string
	newModule string
	fset      *token.FileSet
}

// NewASTReplacer creates a new AST-based import replacer
func NewASTReplacer(oldModule, newModule string) *ASTReplacer {
	return &ASTReplacer{
		oldModule: oldModule,
		newModule: newModule,
		fset:      token.NewFileSet(),
	}
}

// ReplaceInFile replaces import paths in a single Go file
func (r *ASTReplacer) ReplaceInFile(ctx context.Context, filePath string) error {
	// Read file content
	content := gfile.GetContents(filePath)
	if content == "" {
		return nil
	}

	// Parse the file
	file, err := parser.ParseFile(r.fset, filePath, content, parser.ParseComments)
	if err != nil {
		mlog.Debugf("Failed to parse %s: %v", filePath, err)
		return nil // Skip files that can't be parsed
	}

	// Track if any changes were made
	changed := false

	// Traverse and modify imports
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if x.Path != nil {
				importPath := strings.Trim(x.Path.Value, `"`)
				if strings.HasPrefix(importPath, r.oldModule) {
					// Replace only the leading module prefix for clarity and correctness.
					newPath := r.newModule + strings.TrimPrefix(importPath, r.oldModule)
					x.Path.Value = `"` + newPath + `"`
					changed = true
					mlog.Debugf("Replaced import: %s -> %s in %s", importPath, newPath, filePath)
				}
			}
		}
		return true
	})

	if !changed {
		return nil
	}

	// Write back to file without formatting.
	// Formatting will be handled by formatGoFiles after all replacements are done.
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, r.fset, file); err != nil {
		return err
	}

	return gfile.PutContents(filePath, buf.String())
}

// ReplaceInDir replaces import paths in all Go files in a directory (recursively)
func (r *ASTReplacer) ReplaceInDir(ctx context.Context, dir string) error {
	mlog.Printf("Replacing imports: %s -> %s", r.oldModule, r.newModule)

	// Find all .go files
	files, err := findGoFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := r.ReplaceInFile(ctx, file); err != nil {
			mlog.Printf("Failed to process %s: %v", file, err)
		}
	}

	return nil
}

// findGoFiles recursively finds all .go files in a directory
func findGoFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
