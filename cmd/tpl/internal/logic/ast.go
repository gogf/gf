package logic

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

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
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
		g.Log().Debugf(ctx, "Failed to parse %s: %v", filePath, err)
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
					newPath := strings.Replace(importPath, r.oldModule, r.newModule, 1)
					x.Path.Value = `"` + newPath + `"`
					changed = true
					g.Log().Debugf(ctx, "Replaced import: %s -> %s in %s", importPath, newPath, filePath)
				}
			}
		}
		return true
	})

	if !changed {
		return nil
	}

	// Write back to file
	var buf bytes.Buffer
	cfg := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 4,
	}
	if err := cfg.Fprint(&buf, r.fset, file); err != nil {
		return err
	}

	return gfile.PutContents(filePath, buf.String())
}

// ReplaceInDir replaces import paths in all Go files in a directory (recursively)
func (r *ASTReplacer) ReplaceInDir(ctx context.Context, dir string) error {
	g.Log().Infof(ctx, "Replacing imports: %s -> %s", r.oldModule, r.newModule)

	// Find all .go files
	files, err := findGoFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := r.ReplaceInFile(ctx, file); err != nil {
			g.Log().Warningf(ctx, "Failed to process %s: %v", file, err)
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
