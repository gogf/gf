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
	"sync"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

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
		var name, certain = importPackageName(imp, importPath, resolvedNames)
		if !certain {
			// Keep the import if its package name cannot be determined, as removing it
			// mistakenly would break the generated code.
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

// groupImportSpecs groups the given import spec lines into standard library section and
// third party section, sorted by import path within each section.
func groupImportSpecs(specs []string) []string {
	// Keep import blocks containing comments as is.
	for _, spec := range specs {
		if gstr.Contains(spec, "//") || gstr.Contains(spec, "/*") {
			return specs
		}
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

// importPackageName retrieves the package name the import refers to, like "gdb" for
// the import path "github.com/gogf/gf/v2/database/gdb".
//
// The real package name resolved by `resolvePackageNames` is preferred, as the package
// name is not always a segment of the import path, for example
// "github.com/minio/minio-go/v7" is package "minio", "k8s.io/api/core/v1" is package
// "v1" and "go.etcd.io/etcd/client/v3" is package "clientv3".
//
// It returns false as the second result if the package name cannot be determined, in
// which case the import should be kept.
func importPackageName(imp *ast.ImportSpec, importPath string, resolvedNames map[string]string) (name string, certain bool) {
	if imp.Name != nil {
		return imp.Name.Name, true
	}
	if resolvedName, ok := resolvedNames[importPath]; ok && resolvedName != "" {
		return resolvedName, true
	}
	// Fall back to the import path when the package name cannot be resolved, which
	// happens when the import is out of the module context, for example an invalid
	// import path that gf rewrites. The fallback is only used to detect the usages,
	// and the import is kept when the last segment is not a reliable package name.
	if name = packageNameFromPath(importPath); name != "" {
		return name, true
	}
	return "", false
}

// packageNameFromPath derives the package name from an import path, like "net/http"
// returns "http".
//
// It returns an empty name for the versioned paths, whose package name is not reliable,
// for example "k8s.io/api/core/v1" is package "v1" while
// "github.com/gogf/gf/contrib/drivers/mysql/v2" is package "mysql".
func packageNameFromPath(importPath string) string {
	var (
		segments = strings.Split(importPath, "/")
		name     = segments[len(segments)-1]
	)
	// A version suffix like "v1" makes the package name unpredictable.
	if len(name) > 1 && name[0] == 'v' && gstr.IsNumeric(name[1:]) {
		return ""
	}
	if !isValidGoIdentifier(name) {
		return ""
	}
	return name
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
