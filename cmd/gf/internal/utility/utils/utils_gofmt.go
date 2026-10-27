// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements Go source formatting and unused-import cleanup used by
// the gf CLI code generators.

package utils

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// GoFmt formats the source file and removes unused import statements as necessary.
//
// It uses the standard library parser/formatter with astutil instead of
// golang.org/x/tools/imports.Process, which spawns a `go env` subprocess on
// every single call since x/tools v0.28.0, making code generation about 20x
// slower for each file.
//
// Unlike imports.Process, GoFmt never adds the missing imports of the source file, as
// guessing them requires the package name index of the whole module. The code generators
// that need extra imports declare them explicitly with EnsureImports.
//
// The only subprocess involved is a single `go list`, executed by packages.Load at most
// once per process when the real package name of an unnamed import has to be resolved
// (see resolvePackageNames). Everything else, including the unused-import cleanup and the
// import grouping, is done in process.
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

// ImportSpec describes an import statement to be ensured in a Go source file.
// An empty Name means the package is imported with its default package name.
type ImportSpec struct {
	// Name is the optional import alias, like `v1` in `v1 "path/to/api/v1"`.
	Name string
	// Path is the import path, like `github.com/gogf/gf/v2/net/grpcx`.
	Path string
}

// EnsureImports ensures that the given imports exist in the Go source file specified by
// `filePath`. It returns true if the file content is changed, and the file is left untouched
// (not written) when all the given imports already exist.
//
// It is the counterpart of GoFmt, which only removes the unused imports. The code generators
// use EnsureImports to declare the imports the generated code requires, and then call GoFmt
// to format the file and to remove the imports that are not used by the generated content.
func EnsureImports(filePath string, specs ...ImportSpec) (changed bool, err error) {
	if len(specs) == 0 {
		return false, nil
	}
	var content = gfile.GetContents(filePath)
	if gstr.Trim(content) == "" {
		return false, gerror.NewCodef(
			gcode.CodeInvalidParameter, `the file "%s" is empty or does not exist`, filePath,
		)
	}
	var fset = token.NewFileSet()
	var file, parseErr = parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if parseErr != nil {
		return false, parseErr
	}
	for _, spec := range specs {
		if spec.Path == "" {
			return false, gerror.NewCodef(gcode.CodeInvalidParameter, `the import path cannot be empty`)
		}
		if astutil.AddNamedImport(fset, file, spec.Name, spec.Path) {
			changed = true
		}
	}
	if !changed {
		return false, nil
	}
	var buf = new(bytes.Buffer)
	if err = format.Node(buf, fset, file); err != nil {
		return false, err
	}
	if err = gfile.PutContents(filePath, groupImportSections(buf.String())); err != nil {
		return false, err
	}
	return true, nil
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
	var usedNames = collectUsedPackageNames(file)
	// Resolve the real package names of the unnamed imports, as the package name is
	// not always a segment of the import path.
	var unnamedImportPaths = make([]string, 0, len(file.Imports))
	for _, imp := range file.Imports {
		if imp.Name != nil {
			continue
		}
		if importPath, unquoteErr := strconv.Unquote(imp.Path.Value); unquoteErr == nil {
			unnamedImportPaths = append(unnamedImportPaths, importPath)
		}
	}
	var resolvedNames = resolvePackageNames(gfile.Dir(path), unnamedImportPaths)

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
		var importPath, unquoteErr = strconv.Unquote(imp.Path.Value)
		if unquoteErr != nil {
			// Keep the import if its path cannot be quoted, as it cannot be deleted safely.
			continue
		}
		// The package name is not always a segment of the import path, for example
		// "github.com/minio/minio-go/v7" is package "minio" and "k8s.io/api/core/v1"
		// is package "v1", so the real name resolved by `resolvePackageNames` is used
		// for the unnamed imports.
		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			name = resolvedNames[importPath]
		}
		if name == "" {
			// The real package name is unavailable, for example when the file is generated into
			// a directory that is outside of any go module. Fall back to the name assumed from
			// the import path, which is what goimports does, so that the unused imports are
			// removed anyway: a generated file with an unused import left behind does not compile.
			name = assumedPackageName(importPath)
			if name == "" {
				continue
			}
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

// collectUsedPackageNames collects and returns the names used as package qualifiers in the
// given file, like `pkg` in `pkg.Name`.
//
// The identifier of a package qualifier is resolved by the parser to no object, as imported
// packages are not part of the file scope. An identifier that is resolved to an object is
// therefore bound to a declaration inside the file, for example a local variable shadowing
// the imported package name (`strings := T{}; return strings.Name`), and it must not count
// as a usage of the import.
func collectUsedPackageNames(file *ast.File) map[string]struct{} {
	var usedNames = make(map[string]struct{})
	ast.Inspect(file, func(n ast.Node) bool {
		selector, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Obj != nil {
			return true
		}
		usedNames[ident.Name] = struct{}{}
		return true
	})
	return usedNames
}

// groupImportSections inserts a blank line between the standard library imports and the
// other imports in every import block, as goimports (imports.Process) does. The standard
// gofmt does not group imports on its own, and gf generated files rely on this grouping.
//
// The import blocks are located by parsing the content, so that lines inside raw strings
// that look like an import block are never rewritten. Import blocks containing comments
// are kept as is, as rewriting them by lines is risky.
func groupImportSections(content string) string {
	var (
		fset      = token.NewFileSet()
		file, err = parser.ParseFile(fset, "", content, parser.ParseComments)
	)
	if err != nil {
		return content
	}
	// Collect the line ranges of the import blocks, 1-based and inclusive.
	type lineRange struct {
		start, end int
	}
	var ranges = make([]lineRange, 0, 2)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT || !genDecl.Lparen.IsValid() {
			continue
		}
		ranges = append(ranges, lineRange{
			start: fset.Position(genDecl.Lparen).Line,
			end:   fset.Position(genDecl.Rparen).Line,
		})
	}
	if len(ranges) == 0 {
		return content
	}

	var (
		lines  = strings.Split(content, "\n")
		result = make([]string, 0, len(lines))
		cursor = 0 // The index of the first line not handled yet.
	)
	for _, r := range ranges {
		var (
			start = r.start - 1
			end   = r.end - 1
		)
		if start < cursor || start >= end || end >= len(lines) {
			continue
		}
		result = append(result, lines[cursor:start]...)
		result = append(result, lines[start])
		result = append(result, groupImportSpecs(lines[start+1:end])...)
		result = append(result, lines[end])
		cursor = end + 1
	}
	result = append(result, lines[cursor:]...)
	return strings.Join(result, "\n")
}

// groupImportSpecs groups the given import spec lines and sorts them by import path.
//
// The existing blank-line-separated groups are respected: when the block already contains
// more than one group, as the three groups import block written by hand in a user's main.go
// (blank imports / gf packages / project packages), the groups and their order are kept and
// only the specs inside every group are sorted. A block with a single group is split into
// the standard library section and the third party section, as goimports does.
func groupImportSpecs(specs []string) []string {
	// Keep import blocks containing comments as is.
	for _, spec := range specs {
		if gstr.Contains(spec, "//") || gstr.Contains(spec, "/*") {
			return specs
		}
	}
	var existingGroups = splitImportGroups(specs)
	if len(existingGroups) > 1 {
		// Keep the blank-line-separated groups, sort the specs inside every group.
		var result = make([]string, 0, len(specs)+len(existingGroups))
		for index, group := range existingGroups {
			if index > 0 {
				result = append(result, "")
			}
			sort.Slice(group, func(i, j int) bool { return importPathOf(group[i]) < importPathOf(group[j]) })
			result = append(result, group...)
		}
		return result
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
		if isStandardLibraryImportPath(importPathOf(spec)) {
			stdImports = append(stdImports, spec)
		} else {
			thirdImports = append(thirdImports, spec)
		}
	}
	sort.Slice(stdImports, func(i, j int) bool { return importPathOf(stdImports[i]) < importPathOf(stdImports[j]) })
	sort.Slice(thirdImports, func(i, j int) bool { return importPathOf(thirdImports[i]) < importPathOf(thirdImports[j]) })

	var result = make([]string, 0, len(specs)+1)
	result = append(result, stdImports...)
	if len(stdImports) > 0 && len(thirdImports) > 0 {
		result = append(result, "")
	}
	result = append(result, thirdImports...)
	return result
}

// splitImportGroups splits the given import spec lines into the groups separated by blank
// lines. Blank lines are never part of a returned group.
func splitImportGroups(specs []string) [][]string {
	var (
		groups = make([][]string, 0, 2)
		group  []string
	)
	for _, spec := range specs {
		if gstr.Trim(spec) == "" {
			if len(group) > 0 {
				groups = append(groups, group)
				group = nil
			}
			continue
		}
		group = append(group, spec)
	}
	if len(group) > 0 {
		groups = append(groups, group)
	}
	return groups
}

// isStandardLibraryImportPath checks and returns whether given import path is of the
// standard library, like `fmt` and `net/http`.
func isStandardLibraryImportPath(path string) bool {
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

// assumedPackageName returns the package name assumed from the given import path, which is
// used only when the real package name cannot be resolved by `go list`. It follows the same
// rules as goimports: the trailing version suffix like "/v2" is skipped, the "go-" prefix is
// trimmed, and the name is cut at the first "-" or ".", so "github.com/minio/minio-go/v7"
// is assumed as "minio".
func assumedPackageName(importPath string) string {
	var base = path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			var dir = path.Dir(importPath)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	base = strings.TrimPrefix(base, "go-")
	if index := strings.IndexByte(base, '-'); index >= 0 {
		base = base[:index]
	}
	if index := strings.IndexByte(base, '.'); index >= 0 {
		base = base[:index]
	}
	return base
}

// resolvePackageNames resolves and returns the real package names of the given import
// paths in directory `dir`. Every import path is resolved only once and cached, so that
// the resolving happens once instead of for every file.
//
// It is a variable so that tests can stub the resolving.
var resolvePackageNames = func(dir string, importPaths []string) map[string]string {
	var (
		missing []string
		result  = make(map[string]string, len(importPaths))
	)
	packageNamesCacheMu.Lock()
	for _, importPath := range importPaths {
		if name, ok := packageNamesCache[importPath]; ok {
			result[importPath] = name
			continue
		}
		missing = append(missing, importPath)
	}
	packageNamesCacheMu.Unlock()
	if len(missing) == 0 {
		return result
	}

	var resolved, loadErr = loadPackageNames(dir, missing)
	if loadErr != nil || len(resolved) < len(missing) {
		// Fall back to the current working directory, as the directory may not be in a
		// go module, for example a temporary directory that gf generates files into.
		var fallback, fallbackErr = loadPackageNames(gfile.Pwd(), missing)
		if fallbackErr == nil {
			if resolved == nil {
				resolved = make(map[string]string, len(missing))
			}
			for importPath, name := range fallback {
				resolved[importPath] = name
			}
		}
	}

	packageNamesCacheMu.Lock()
	for _, importPath := range missing {
		// An empty name means the package name is unavailable, which is also cached to
		// avoid resolving it for every file.
		packageNamesCache[importPath] = resolved[importPath]
		result[importPath] = resolved[importPath]
	}
	packageNamesCacheMu.Unlock()
	return result
}

// packageNamesCache caches the resolved package names of the import paths.
var (
	packageNamesCache   = make(map[string]string)
	packageNamesCacheMu sync.Mutex
)

// loadPackageNames loads the real package names of the given import paths in directory
// `dir` using go/packages, which runs a single `go list` for all the paths.
func loadPackageNames(dir string, importPaths []string) (map[string]string, error) {
	var pkgs, err = packages.Load(&packages.Config{
		Mode: packages.NeedName,
		Dir:  dir,
	}, importPaths...)
	if err != nil {
		return nil, err
	}
	var names = make(map[string]string, len(pkgs))
	for _, pkg := range pkgs {
		if pkg.PkgPath == "" || pkg.Name == "" {
			continue
		}
		names[pkg.PkgPath] = pkg.Name
	}
	return names, nil
}
