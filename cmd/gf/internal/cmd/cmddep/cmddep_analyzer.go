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

// analyzer handles dependency analysis.
type analyzer struct {
	packages     map[string]*goPackage
	modulePrefix string
	visited      map[string]bool
	edges        map[string]bool
}

// newAnalyzer creates a new dependency analyzer.
func newAnalyzer() *analyzer {
	return &analyzer{
		packages: make(map[string]*goPackage),
		visited:  make(map[string]bool),
		edges:    make(map[string]bool),
	}
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

// loadPackages loads package information using go list.
func (a *analyzer) loadPackages(ctx context.Context, pkgPath string) error {
	// Load main packages first
	cmd := fmt.Sprintf("go list -json %s", pkgPath)
	result, err := gproc.ShellExec(ctx, cmd)
	if err != nil {
		// Try to get more detailed error information
		detailCmd := fmt.Sprintf("go list %s 2>&1", pkgPath)
		detailResult, _ := gproc.ShellExec(ctx, detailCmd)
		return fmt.Errorf("failed to execute go list: %v, details: %s", err, detailResult)
	}

	// Parse JSON stream (multiple JSON objects)
	decoder := json.NewDecoder(strings.NewReader(result))
	for decoder.More() {
		var pkg goPackage
		if err := decoder.Decode(&pkg); err != nil {
			continue
		}
		a.packages[pkg.ImportPath] = &pkg
	}

	// For external dependency analysis, also load dependencies
	// This is optional and won't fail the entire operation
	cmd = fmt.Sprintf("go list -json -deps %s", pkgPath)
	result, err = gproc.ShellExec(ctx, cmd)
	if err == nil {
		// Parse dependency JSON stream
		decoder = json.NewDecoder(strings.NewReader(result))
		for decoder.More() {
			var pkg goPackage
			if err := decoder.Decode(&pkg); err != nil {
				continue
			}
			// Only add if not already present
			if _, exists := a.packages[pkg.ImportPath]; !exists {
				a.packages[pkg.ImportPath] = &pkg
			}
		}
	}
	return nil
}

// filterDeps filters dependencies based on options.
func (a *analyzer) filterDeps(deps []string, in Input) []string {
	result := make([]string, 0)
	seen := make(map[string]bool)
	for _, original := range deps {
		dep := original
		if in.MainOnly {
			dep = a.getModuleRoot(original)
		}

		if a.shouldInclude(dep, in) && !seen[dep] {
			seen[dep] = true
			result = append(result, dep)
		}
	}
	return result
}

// shouldInclude checks if a dependency should be included.
func (a *analyzer) shouldInclude(dep string, in Input) bool {
	// Exclude standard library if requested
	if in.NoStd && a.isStdLib(dep) {
		return false
	}

	isInternal := a.modulePrefix != "" && gstr.HasPrefix(dep, a.modulePrefix)
	
	// Handle main-only filtering - only keep module root packages
	if in.MainOnly {
		if dep != a.getModuleRoot(dep) {
			return false
		}
	}
	
	// Handle internal/external filtering
	if in.Internal && in.External {
		// Show both internal and external
		return true
	} else if in.Internal && !in.External {
		// Show only internal packages
		return isInternal
	} else if !in.Internal && in.External {
		// Show only external packages
		return !isInternal
	} else {
		// Default behavior: show internal packages only
		return isInternal
	}
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

// collectEdges collects all dependency edges.
func (a *analyzer) collectEdges(in Input) map[string]bool {
	edges := make(map[string]bool)
	a.visited = make(map[string]bool)

	for _, pkg := range a.packages {
		a.collectEdgesRecursive(pkg, in, edges, 0)
	}
	return edges
}

func (a *analyzer) collectEdgesRecursive(pkg *goPackage, in Input, edges map[string]bool, depth int) {
	if in.Depth > 0 && depth >= in.Depth {
		return
	}

	fromName := a.shortName(pkg.ImportPath, in.Group)
	deps := a.filterDeps(pkg.Imports, in)

	for _, dep := range deps {
		toName := a.shortName(dep, in.Group)
		if fromName != toName && toName != "" && fromName != "" {
			edge := fmt.Sprintf("%s --> %s", a.sanitizeName(fromName), a.sanitizeName(toName))
			edges[edge] = true
		}

		if !a.visited[dep] {
			a.visited[dep] = true
			if depPkg, ok := a.packages[dep]; ok {
				a.collectEdgesRecursive(depPkg, in, edges, depth+1)
			}
		}
	}
}

// getDependencyStats returns statistics about dependencies.
func (a *analyzer) getDependencyStats(_ Input) map[string]any {
	stats := make(map[string]any)
	
	var internalCount, externalCount, stdlibCount int
	externalGroups := make(map[string]int)
	
	for _, pkg := range a.packages {
		if !a.shouldInclude(pkg.ImportPath, Input{
			Internal: true,
			External: true,
			NoStd:    false,
		}) {
			continue
		}
		
		if a.isStdLib(pkg.ImportPath) {
			stdlibCount++
		} else if a.modulePrefix != "" && gstr.HasPrefix(pkg.ImportPath, a.modulePrefix) {
			internalCount++
		} else {
			externalCount++
			group := a.getExternalGroup(pkg.ImportPath)
			externalGroups[group]++
		}
	}
	
	stats["total"] = len(a.packages)
	stats["internal"] = internalCount
	stats["external"] = externalCount
	stats["stdlib"] = stdlibCount
	stats["external_groups"] = externalGroups
	
	return stats
}
