// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
)

// goPackage represents a Go package from go list -json output.
type goPackage struct {
	ImportPath string   `json:"ImportPath"`
	Imports    []string `json:"Imports"`
	Deps       []string `json:"Deps"`
	Standard   bool     `json:"Standard"`
	Module     struct {
		Path string `json:"Path"`
	} `json:"Module"`
}

// depNode represents a node in the dependency tree.
type depNode struct {
	Package      string     `json:"package"`
	Dependencies []*depNode `json:"dependencies,omitempty"`
}

// PackageKind indicates the kind of a Go package
type PackageKind int

const (
	KindInternal PackageKind = iota // Internal to main module
	KindExternal                      // External dependency
	KindStdLib                        // Standard library
)

// PackageInfo represents unified information about a Go package.
// This is the core data model for the refactored dependency analyzer.
// It consolidates package information from go list output with additional
// metadata for filtering and traversal.
type PackageInfo struct {
	ImportPath   string       // Full import path (e.g., github.com/gogf/gf/v2/os/gfile)
	ModulePath   string       // Module path (e.g., github.com/gogf/gf/v2)
	Kind         PackageKind  // Package classification (Internal/External/StdLib)
	Tier         int          // Package tier: 0=module root, 1=top-level, 2+=nested
	Imports      []string     // Direct imports of this package
	IsStdLib     bool         // Standard library marker (from go list)
	IsModuleRoot bool         // Is this the root package of its module
}

// FilterOptions represents filtering criteria for dependency analysis.
// It provides a clear, normalized representation of user filtering preferences.
// Usage:
//   opts := &FilterOptions{
//       IncludeInternal: true,
//       IncludeExternal: false,
//       IncludeStdLib:   false,
//       MainModuleOnly:  false,
//       Depth:           3,
//   }
type FilterOptions struct {
	IncludeInternal bool // Include internal packages from main module
	IncludeExternal bool // Include external dependencies
	IncludeStdLib   bool // Include standard library packages
	MainModuleOnly  bool // Only analyze packages from main module (exclude submodules)
	Depth           int  // Maximum traversal depth (0 = unlimited)
}

// TraversalContext manages state during dependency tree traversal.
// It centralizes visited tracking, depth management, and filtering logic
// to ensure consistent behavior across different output formats.
type TraversalContext struct {
	visited   map[string]bool // Track visited packages to prevent cycles
	depth     int             // Current traversal depth
	maxDepth  int             // Maximum traversal depth
	options   *FilterOptions  // Filtering criteria
	store     *PackageStore   // Reference to package store
}

// PackageStore manages a collection of packages and provides unified data access.
// This centralizes all package data and implements traversal algorithms.
// It replaces the scattered data access patterns in the original analyzer.
type PackageStore struct {
	packages      map[string]*PackageInfo // Package data indexed by import path
	modulePrefix  string                   // Main module path (from go.mod)
	sortedPkgs    []string                 // Cached sorted package list
	internalCount int                      // Cached count of internal packages
	externalCount int                      // Cached count of external packages
}

// analyzer handles dependency analysis.
type analyzer struct {
	packages     map[string]*goPackage
	modulePrefix string
	visited      map[string]bool
	edges        map[string]bool
	store        *PackageStore // New unified package store
}

// newAnalyzer creates a new dependency analyzer.
func newAnalyzer() *analyzer {
	return &analyzer{
		packages: make(map[string]*goPackage),
		visited:  make(map[string]bool),
		edges:    make(map[string]bool),
		store:    &PackageStore{},
	}
}

// newPackageStore creates a new package store.
func newPackageStore(modulePrefix string) *PackageStore {
	return &PackageStore{
		packages:     make(map[string]*PackageInfo),
		modulePrefix: modulePrefix,
	}
}

// identifyPackageKind determines the kind of a package.
func (ps *PackageStore) identifyPackageKind(pkg *PackageInfo) PackageKind {
	if pkg.IsStdLib {
		return KindStdLib
	}
	if ps.modulePrefix != "" && gstr.HasPrefix(pkg.ImportPath, ps.modulePrefix) {
		return KindInternal
	}
	return KindExternal
}

// detectModulePrefix reads go.mod to get the module path.
func (a *analyzer) detectModulePrefix() string {
	content := gfile.GetContents("go.mod")
	if content == "" {
		return ""
	}
	lines := gstr.Split(content, "\n")
	for _, line := range lines {
		line = gstr.Trim(line)
		if gstr.HasPrefix(line, "module ") {
			return gstr.Trim(line[7:])
		}
	}
	return ""
}

// loadPackages loads package information using go list with optimized approach.
// OPTIMIZATION: Reduced from 3 separate go list calls to 2 efficient calls:
// Previously:
//   1. go list -json %s                (target packages only)
//   2. go list -json -deps %s          (with dependencies)
//   3. go list -json -m all            (all modules)
// Now (optimized):
//   1. go list -json -m all            (all modules - fast, definitive)
//   2. go list -json -deps ./...       (all packages with dependencies)
func (a *analyzer) loadPackages(ctx context.Context, pkgPath string) error {
	// First, load module information - this is fast and provides module metadata
	// Load all module dependencies using go list -m all
	// This ensures we capture all modules declared in go.mod, including indirect ones
	moduleCmd := "go list -json -m all"
	moduleResult, err := gproc.ShellExec(ctx, moduleCmd)
	if err != nil {
		// Modules loading is optional, continue with package loading
		moduleResult = ""
	}

	// Parse module information if available
	if moduleResult != "" {
		moduleDecoder := json.NewDecoder(strings.NewReader(moduleResult))
		for moduleDecoder.More() {
			var mod struct {
				Path string `json:"Path"`
			}
			if err := moduleDecoder.Decode(&mod); err != nil {
				continue
			}
			// Create a virtual package entry for modules not found in code analysis
			// This ensures all declared dependencies are visible in the graph
			if _, exists := a.packages[mod.Path]; !exists {
				a.packages[mod.Path] = &goPackage{
					ImportPath: mod.Path,
					Imports:    []string{},
					Deps:       []string{},
					Standard:   false,
				}
			}
		}
	}

	// Second, load package information with all dependencies
	// Use go list -json -deps to get complete package dependency information
	cmd := fmt.Sprintf("go list -json -deps %s", pkgPath)
	result, err := gproc.ShellExec(ctx, cmd)
	if err != nil {
		// Try to get more detailed error information
		detailCmd := fmt.Sprintf("go list %s 2>&1", pkgPath)
		detailResult, _ := gproc.ShellExec(ctx, detailCmd)
		return fmt.Errorf("failed to execute go list: %v, details: %s", err, detailResult)
	}

	// Parse the package JSON stream (multiple JSON objects)
	decoder := json.NewDecoder(strings.NewReader(result))
	for decoder.More() {
		var pkg goPackage
		if err := decoder.Decode(&pkg); err != nil {
			continue
		}
		a.packages[pkg.ImportPath] = &pkg
	}

	return nil
}



// convertInputToFilterOptions converts legacy Input to new FilterOptions.
func (a *analyzer) convertInputToFilterOptions(in Input) *FilterOptions {
	opts := &FilterOptions{
		IncludeInternal: in.Internal,
		IncludeExternal: in.External,
		IncludeStdLib:   !in.NoStd,
		MainModuleOnly:  in.MainOnly,
		Depth:           in.Depth,
	}
	
	// Apply default: if neither internal nor external, include internal only
	if !in.Internal && !in.External {
		opts.IncludeInternal = true
		opts.IncludeExternal = false
	}
	
	return opts
}

// isStdLib checks if a package is from standard library.
func (a *analyzer) isStdLib(pkg string) bool {
	// Standard library packages don't contain dots in the first path segment
	if strings.Contains(pkg, ".") {
		return false
	}
	// Check if it's in our loaded packages and marked as standard
	if p, ok := a.packages[pkg]; ok {
		return p.Standard
	}
	return true
}

// isModuleRootPackage checks if a package path is the root package of its module.
// It prefers go list Module.Path metadata; if missing, it falls back to guessing by domain/repo segments.
func (a *analyzer) isModuleRootPackage(pkg string) bool {
	p, ok := a.packages[pkg]
	if ok {
		// Standard library has no module, treat as root
		if p.Module.Path == "" {
			return true
		}
		return p.Module.Path == p.ImportPath
	}

	// Fallback: derive a plausible module root from the import path
	return pkg == guessModuleRoot(pkg)
}

// getModuleRoot returns the module root path for a package, using Module metadata when available.
func (a *analyzer) getModuleRoot(pkg string) string {
	if p, ok := a.packages[pkg]; ok {
		if p.Module.Path != "" {
			return p.Module.Path
		}
	}
	return guessModuleRoot(pkg)
}

// guessModuleRoot tries to infer the module root path from an import path when Module metadata is missing.
// It keeps domain/owner/repo and also preserves a trailing /vN version segment if present.
func guessModuleRoot(pkg string) string {
	parts := strings.Split(pkg, "/")
	if len(parts) < 3 {
		return pkg
	}

	rootLen := 3 // domain/owner/repo
	// Handle semantic import path version like .../v2
	if len(parts) > 3 && strings.HasPrefix(parts[3], "v") {
		rootLen = 4
	}

	if rootLen > len(parts) {
		rootLen = len(parts)
	}
	return strings.Join(parts[:rootLen], "/")
}

// Normalize normalizes filter options based on default behavior.
func (opts *FilterOptions) Normalize(modulePrefix string) error {
	// If neither internal nor external is explicitly set to true,
	// use default behavior: internal only
	if !opts.IncludeInternal && !opts.IncludeExternal {
		opts.IncludeInternal = true
		opts.IncludeExternal = false
	}
	
	// Always include stdlib by default unless explicitly excluded
	if opts.IncludeStdLib == false {
		// This is the default (NoStd=true), stdlib is excluded
	} else {
		opts.IncludeStdLib = true
	}
	
	return nil
}

// ShouldInclude determines if a package should be included based on filter options.
func (opts *FilterOptions) ShouldInclude(pkg *PackageInfo) bool {
	// Filter by kind
	switch pkg.Kind {
	case KindStdLib:
		if !opts.IncludeStdLib {
			return false
		}
	case KindInternal:
		if !opts.IncludeInternal {
			return false
		}
	case KindExternal:
		if !opts.IncludeExternal {
			return false
		}
	}
	
	return true
}

// Visit marks a package as visited and returns whether it was already visited.
func (tc *TraversalContext) Visit(pkg string) bool {
	if tc.visited[pkg] {
		return true
	}
	tc.visited[pkg] = true
	return false
}

// GetDependencies returns the dependencies of a package according to filter options.
func (tc *TraversalContext) GetDependencies(pkg string) []string {
	pkgInfo, ok := tc.store.packages[pkg]
	if !ok {
		return []string{}
	}
	
	result := make([]string, 0)
	seen := make(map[string]bool)
	
	for _, dep := range pkgInfo.Imports {
		if seen[dep] {
			continue
		}
		
		depInfo, ok := tc.store.packages[dep]
		if !ok {
			continue
		}
		
		if tc.options.ShouldInclude(depInfo) {
			seen[dep] = true
			result = append(result, dep)
		}
	}
	
	return result
}

// isMainModulePackage checks if a package belongs to the main module (not a submodule).
func (a *analyzer) isMainModulePackage(pkg string) bool {
	if a.modulePrefix == "" {
		return true // If no module prefix, consider all as main module
	}
	
	if !gstr.HasPrefix(pkg, a.modulePrefix) {
		return false // Not even in our module
	}
	
	// Remove the module prefix to get the relative path
	relativePath := gstr.TrimLeft(pkg[len(a.modulePrefix):], "/")
	if relativePath == "" {
		return true // This is the root module itself
	}
	
	// Check if this path contains a go.mod file (indicating a submodule)
	// We check from the most specific path up to the root
	parts := gstr.Split(relativePath, "/")
	for i := len(parts); i > 0; i-- {
		subPath := gstr.Join(parts[:i], "/")
		if subPath != "" && gfile.Exists(subPath+"/go.mod") {
			// Found a go.mod file in a subdirectory, this indicates a submodule
			return false
		}
	}
	
	return true // This is part of the main module
}

// shortName returns a shortened package name.
func (a *analyzer) shortName(pkg string, group bool) string {
	if a.modulePrefix != "" && gstr.HasPrefix(pkg, a.modulePrefix) {
		short := gstr.TrimLeft(pkg[len(a.modulePrefix):], "/")
		if group {
			// Return only top-level directory
			parts := gstr.Split(short, "/")
			if len(parts) > 0 {
				return parts[0]
			}
		}
		return short
	}
	
	// For external packages, handle grouping differently
	if group {
		return a.getExternalGroup(pkg)
	}
	return pkg
}

// getExternalGroup returns the group name for external packages.
func (a *analyzer) getExternalGroup(pkg string) string {
	// For standard library packages
	if a.isStdLib(pkg) {
		return "stdlib"
	}
	
	// For external packages, group by domain/organization
	parts := gstr.Split(pkg, "/")
	if len(parts) > 0 {
		// Handle common patterns like github.com/user/repo
		if len(parts) >= 3 && (parts[0] == "github.com" || parts[0] == "gitlab.com" || parts[0] == "bitbucket.org") {
			return parts[0] + "/" + parts[1]
		}
		// For other domains, use the domain name
		if gstr.Contains(parts[0], ".") {
			return parts[0]
		}
		// For simple names, use the first part
		return parts[0]
	}
	return pkg
}

// sanitizeName makes a name safe for mermaid/dot output.
func (a *analyzer) sanitizeName(name string) string {
	return gstr.Replace(name, "/", "_")
}

// getSortedPackages returns sorted package paths.
func (a *analyzer) getSortedPackages() []string {
	pkgs := make([]string, 0, len(a.packages))
	for pkg := range a.packages {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)
	return pkgs
}

// collectEdges collects all dependency edges using new traversal system.
func (a *analyzer) collectEdges(in Input) map[string]bool {
	opts := a.convertInputToFilterOptions(in)
	opts.Normalize(a.modulePrefix)
	
	store := a.buildPackageStore()
	edges := make(map[string]bool)
	visited := make(map[string]bool)

	for pkgPath := range store.packages {
		a.collectEdgesRecursiveNew(pkgPath, opts, store, edges, visited, 0, in)
	}
	return edges
}

// collectEdgesRecursiveNew recursively collects edges using new system.
func (a *analyzer) collectEdgesRecursiveNew(pkgPath string, opts *FilterOptions, store *PackageStore, edges map[string]bool, visited map[string]bool, depth int, in Input) {
	if in.Depth > 0 && depth >= in.Depth {
		return
	}

	pkgInfo, ok := store.packages[pkgPath]
	if !ok || !opts.ShouldInclude(pkgInfo) {
		return
	}

	if visited[pkgPath] {
		return
	}
	visited[pkgPath] = true

	fromName := a.shortName(pkgPath, in.Group)
	
	for _, dep := range pkgInfo.Imports {
		depInfo, ok := store.packages[dep]
		if !ok || !opts.ShouldInclude(depInfo) {
			continue
		}

		toName := a.shortName(dep, in.Group)
		if fromName != toName && toName != "" && fromName != "" {
			edge := fmt.Sprintf("%s --> %s", a.sanitizeName(fromName), a.sanitizeName(toName))
			edges[edge] = true
		}

		a.collectEdgesRecursiveNew(dep, opts, store, edges, visited, depth+1, in)
	}
}

// getDependencyStats returns statistics about dependencies using new system.
func (a *analyzer) getDependencyStats(_ Input) map[string]any {
	stats := make(map[string]any)
	
	var internalCount, externalCount, stdlibCount int
	externalGroups := make(map[string]int)
	
	store := a.buildPackageStore()
	opts := &FilterOptions{
		IncludeInternal: true,
		IncludeExternal: true,
		IncludeStdLib:   true,
		MainModuleOnly:  false,
		Depth:           0,
	}
	
	for _, pkgInfo := range store.packages {
		if !opts.ShouldInclude(pkgInfo) {
			continue
		}
		
		if pkgInfo.IsStdLib {
			stdlibCount++
		} else if pkgInfo.Kind == KindInternal {
			internalCount++
		} else if pkgInfo.Kind == KindExternal {
			externalCount++
			group := a.getExternalGroup(pkgInfo.ImportPath)
			externalGroups[group]++
		}
	}
	
	stats["total"] = len(store.packages)
	stats["internal"] = internalCount
	stats["external"] = externalCount
	stats["stdlib"] = stdlibCount
	stats["external_groups"] = externalGroups
	
	return stats
}

// TraverseDependencies traverses dependencies starting from root using filter options.
func (ps *PackageStore) TraverseDependencies(
	root string,
	options *FilterOptions,
) []string {
	ctx := &TraversalContext{
		visited:  make(map[string]bool),
		maxDepth: options.Depth,
		options:  options,
		store:    ps,
	}
	
	result := make([]string, 0)
	ps.traverseRecursive(root, ctx, &result)
	return result
}

// traverseRecursive recursively traverses dependency tree.
func (ps *PackageStore) traverseRecursive(
	pkg string,
	ctx *TraversalContext,
	result *[]string,
) {
	if ctx.maxDepth > 0 && ctx.depth >= ctx.maxDepth {
		return
	}
	
	if ctx.Visit(pkg) {
		return // Already visited
	}
	
	pkgInfo, ok := ps.packages[pkg]
	if !ok {
		return
	}
	
	if !ctx.options.ShouldInclude(pkgInfo) {
		return // Filtered out
	}
	
	*result = append(*result, pkg)
	
	ctx.depth++
	for _, dep := range pkgInfo.Imports {
		ps.traverseRecursive(dep, ctx, result)
	}
	ctx.depth--
}

// TraverseReverse traverses reverse dependencies (reverse of dependency tree).
func (ps *PackageStore) TraverseReverse(
	target string,
	options *FilterOptions,
) []string {
	result := make([]string, 0)
	
	// Build reverse dependency map on-the-fly
	for _, pkg := range ps.packages {
		for _, dep := range pkg.Imports {
			if dep == target && options.ShouldInclude(pkg) {
				result = append(result, pkg.ImportPath)
			}
		}
	}
	
	sort.Strings(result)
	return result
}

// buildPackageStore converts current analyzer state to PackageStore for new traversal system.
func (a *analyzer) buildPackageStore() *PackageStore {
	store := newPackageStore(a.modulePrefix)
	
	// Convert go packages to PackageInfo
	for path, goPkg := range a.packages {
		pkgInfo := &PackageInfo{
			ImportPath: path,
			ModulePath: goPkg.Module.Path,
			IsStdLib:   goPkg.Standard,
			Imports:    goPkg.Imports,
		}
		pkgInfo.Kind = store.identifyPackageKind(pkgInfo)
		store.packages[path] = pkgInfo
	}
	
	return store
}

// getFilteredPackages returns all packages matching filter options using new system.
func (a *analyzer) getFilteredPackages(in Input) []*PackageInfo {
	opts := a.convertInputToFilterOptions(in)
	opts.Normalize(a.modulePrefix)
	
	store := a.buildPackageStore()
	result := make([]*PackageInfo, 0)
	
	for _, pkgInfo := range store.packages {
		if opts.ShouldInclude(pkgInfo) {
			result = append(result, pkgInfo)
		}
	}
	
	return result
}

// getFilteredDependencies returns filtered dependencies of a package using new system.
func (a *analyzer) getFilteredDependencies(pkgPath string, in Input) []string {
	_, ok := a.packages[pkgPath]
	if !ok {
		return []string{}
	}
	
	opts := a.convertInputToFilterOptions(in)
	opts.Normalize(a.modulePrefix)
	
	store := a.buildPackageStore()
	ctx := &TraversalContext{
		visited: make(map[string]bool),
		options: opts,
		store:   store,
	}
	
	return ctx.GetDependencies(pkgPath)
}
