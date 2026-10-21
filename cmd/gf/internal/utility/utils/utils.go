// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// GoFmt formats the source file and removes unused import statements as necessary.
//
// It uses the standard library parser/formatter with astutil instead of
// golang.org/x/tools/imports.Process, which spawns a `go env` subprocess on
// every single call since x/tools v0.28.0, making code generation about 20x
// slower for each file.
func GoFmt(path string) {
	replaceFunc := func(path, content string) string {
		res, err := goFormatContent(path, content)
		if err != nil {
			mlog.Printf(`error format "%s" go files: %v`, path, err)
			return content
		}
		return res
	}

	var err error
	if gfile.IsFile(path) {
		// File format.
		if gfile.ExtName(path) != "go" {
			return
		}
		err = gfile.ReplaceFileFunc(replaceFunc, path)
	} else {
		// Folder format.
		err = gfile.ReplaceDirFunc(replaceFunc, path, "*.go", true)
	}
	if err != nil {
		mlog.Printf(`error format "%s" go files: %v`, path, err)
	}
}

// goFormatContent formats the content of given file and removes its unused imports.
func goFormatContent(path, content string) (string, error) {
	var (
		fset      = token.NewFileSet()
		file, err = parser.ParseFile(fset, path, content, parser.ParseComments|parser.AllErrors)
	)
	if err != nil {
		return "", err
	}
	// Collect all the identifiers used as package qualifiers, like `pkg` in `pkg.Name`.
	// Cross package references must be qualified by the package name, so this covers
	// all the usages of the imported packages (dot imports are kept unconditionally).
	var usedNames = make(map[string]struct{})
	ast.Inspect(file, func(n ast.Node) bool {
		if selector, ok := n.(*ast.SelectorExpr); ok {
			if ident, ok := selector.X.(*ast.Ident); ok {
				usedNames[ident.Name] = struct{}{}
			}
		}
		return true
	})
	// Collect the imports to delete firstly, as astutil.DeleteImport modifies
	// `file.Imports` that we are iterating, which skips imports and leaves
	// unused ones behind.
	type importToDelete struct {
		name string
		path string
	}
	var importsToDelete []importToDelete
	for _, imp := range file.Imports {
		// Never remove blank or dot imports, as their usages cannot be determined.
		if imp.Name != nil && (imp.Name.Name == "_" || imp.Name.Name == ".") {
			continue
		}
		var (
			importPath, unquoteErr = strconv.Unquote(imp.Path.Value)
			name, certain          = importPackageName(imp)
		)
		if unquoteErr != nil {
			// Keep the import if its path cannot be quoted, as it cannot be deleted safely.
			continue
		}
		if !certain {
			// Keep the import if its package name cannot be determined from the import
			// path, as removing it mistakenly would break the generated code.
			continue
		}
		if _, ok := usedNames[name]; ok {
			continue
		}
		var item = importToDelete{path: importPath}
		if imp.Name != nil {
			item.name = imp.Name.Name
		}
		importsToDelete = append(importsToDelete, item)
	}
	for _, item := range importsToDelete {
		if item.name != "" {
			astutil.DeleteNamedImport(fset, file, item.name, item.path)
		} else {
			astutil.DeleteImport(fset, file, item.path)
		}
	}
	// Remove the empty import declarations like `import ()` left by deleting imports.
	var decls = file.Decls[:0]
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT && len(genDecl.Specs) == 0 {
			continue
		}
		decls = append(decls, decl)
	}
	file.Decls = decls

	var buf = new(bytes.Buffer)
	if err = format.Node(buf, fset, file); err != nil {
		return "", err
	}
	return groupImportSections(buf.String()), nil
}

// groupImportSections inserts a blank line between the standard library imports and the
// other imports in every import block, as goimports (imports.Process) does. The standard
// gofmt does not group imports on its own, and gf generated files rely on this grouping.
//
// Import blocks containing comments are kept as is, as rewriting them by lines is risky.
func groupImportSections(content string) string {
	var (
		lines  = strings.Split(content, "\n")
		result = make([]string, 0, len(lines))
	)
	for i := 0; i < len(lines); i++ {
		if gstr.Trim(lines[i]) != "import (" {
			result = append(result, lines[i])
			continue
		}
		// Find the end of the current import block.
		var (
			end   = i + 1
			specs = make([]string, 0, 8)
		)
		for ; end < len(lines); end++ {
			if gstr.Trim(lines[end]) == ")" {
				break
			}
			specs = append(specs, lines[end])
		}
		if end >= len(lines) {
			// Unbalanced parentheses, keep the content as is.
			result = append(result, lines[i])
			continue
		}
		// Keep import blocks containing comments as is.
		var handleable = true
		for _, spec := range specs {
			if gstr.Contains(spec, "//") || gstr.Contains(spec, "/*") {
				handleable = false
				break
			}
		}
		if !handleable {
			result = append(result, lines[i])
			result = append(result, specs...)
			result = append(result, lines[end])
			i = end
			continue
		}
		// Split imports into standard library group and third party group.
		var (
			stdImports   = make([]string, 0, len(specs))
			thirdImports = make([]string, 0, len(specs))
		)
		for _, spec := range specs {
			if gstr.Trim(spec) == "" {
				continue
			}
			if isStandardLibraryImport(spec) {
				stdImports = append(stdImports, spec)
			} else {
				thirdImports = append(thirdImports, spec)
			}
		}
		sort.Slice(stdImports, func(i, j int) bool { return importPathOf(stdImports[i]) < importPathOf(stdImports[j]) })
		sort.Slice(thirdImports, func(i, j int) bool { return importPathOf(thirdImports[i]) < importPathOf(thirdImports[j]) })

		result = append(result, lines[i])
		result = append(result, stdImports...)
		if len(stdImports) > 0 && len(thirdImports) > 0 {
			result = append(result, "")
		}
		result = append(result, thirdImports...)
		result = append(result, lines[end])
		i = end
	}
	return strings.Join(result, "\n")
}

// isStandardLibraryImport checks and returns whether given import spec line is a
// standard library import, like `\t"fmt"` or `\tnet/http "net/http"`.
func isStandardLibraryImport(spec string) bool {
	var path = importPathOf(spec)
	if path == "" {
		return false
	}
	var firstSegment = path
	if index := strings.IndexByte(path, '/'); index >= 0 {
		firstSegment = path[:index]
	}
	// Standard library import paths never contain a dot in the first segment.
	return !gstr.Contains(firstSegment, ".")
}

// importPathOf retrieves the import path from given import spec line, like
// `\t"fmt"` returns `fmt`, `\talias "github.com/x/y"` returns `github.com/x/y`.
func importPathOf(spec string) string {
	var (
		start = strings.IndexByte(spec, '"')
		end   = strings.LastIndexByte(spec, '"')
	)
	if start < 0 || end <= start {
		return ""
	}
	return spec[start+1 : end]
}

// importPackageName retrieves the package name the import refers to, like "gdb" for
// import path "github.com/gogf/gf/v2/database/gdb".
//
// It returns false as the second result if the package name cannot be determined from
// the import path, for example the path "github.com/minio/minio-go/v7" whose package
// name is "minio", which is not any segment of the path.
func importPackageName(imp *ast.ImportSpec) (name string, certain bool) {
	if imp.Name != nil {
		return imp.Name.Name, true
	}
	path, err := strconv.Unquote(imp.Path.Value)
	if err != nil {
		return "", false
	}
	var paths = strings.Split(path, "/")
	name = paths[len(paths)-1]
	// The last segment may be a version suffix like "v2", in which case the
	// package name is usually the second-to-last segment.
	if len(paths) > 1 && len(name) > 1 && name[0] == 'v' && gstr.IsNumeric(name[1:]) {
		name = paths[len(paths)-2]
	}
	// The package name cannot be determined if it is not a valid go identifier,
	// like "minio-go" whose package name is "minio".
	if !isValidGoIdentifier(name) {
		return "", false
	}
	return name, true
}

// isValidGoIdentifier checks and returns whether given string is a valid go identifier.
func isValidGoIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		var c = name[i]
		if c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			if i == 0 && c >= '0' && c <= '9' {
				return false
			}
			continue
		}
		return false
	}
	return true
}

// GoModTidy executes `go mod tidy` at specified directory `dirPath`.
func GoModTidy(ctx context.Context, dirPath string) error {
	command := fmt.Sprintf(`cd %s && go mod tidy`, dirPath)
	err := gproc.ShellRun(ctx, command)
	return err
}

// IsFileDoNotEdit checks and returns whether file contains `do not edit` key.
func IsFileDoNotEdit(filePath string) bool {
	if !gfile.Exists(filePath) {
		return true
	}
	return gstr.Contains(gfile.GetContents(filePath), consts.DoNotEditKey)
}

// ReplaceGeneratedContentGFV2 replaces generated go content from goframe v1 to v2.
func ReplaceGeneratedContentGFV2(folderPath string) (err error) {
	return gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			content = gstr.Replace(content, `"github.com/gogf/gf/v2/contrib/`, `"github.com/gogf/gf/contrib/`)
			return content
		}
		return content
	}, folderPath, "*.go", true)
}

// GetImportPath calculates and returns the golang import path for given `dirPath`.
// Note that it needs a `go.mod` in current working directory or parent directories to detect the path.
func GetImportPath(dirPath string) string {
	// If `filePath` does not exist, create it firstly to find the import path.
	var realPath = gfile.RealPath(dirPath)
	if realPath == "" {
		_ = gfile.Mkdir(dirPath)
		realPath = gfile.RealPath(dirPath)
	}

	var (
		newDir     = gfile.Dir(realPath)
		oldDir     string
		suffix     = gfile.Basename(dirPath)
		goModName  = "go.mod"
		goModPath  string
		importPath string
	)
	for {
		goModPath = gfile.Join(newDir, goModName)
		if gfile.Exists(goModPath) {
			match, _ := gregex.MatchString(`^module\s+(.+)\s*`, gfile.GetContents(goModPath))
			importPath = gstr.Trim(match[1]) + "/" + suffix
			importPath = gstr.Replace(importPath, `\`, `/`)
			importPath = gstr.TrimRight(importPath, `/`)
			return importPath
		}
		oldDir = newDir
		newDir = gfile.Dir(oldDir)
		if newDir == oldDir {
			return ""
		}
		suffix = gfile.Basename(oldDir) + "/" + suffix
	}
}

// GetModPath retrieves and returns the file path of go.mod for current project.
func GetModPath() string {
	var (
		oldDir    = gfile.Pwd()
		newDir    = oldDir
		goModName = "go.mod"
		goModPath string
	)
	for {
		goModPath = gfile.Join(newDir, goModName)
		if gfile.Exists(goModPath) {
			return goModPath
		}
		newDir = gfile.Dir(oldDir)
		if newDir == oldDir {
			break
		}
		oldDir = newDir
	}
	return ""
}
