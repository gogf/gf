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
	cmd := fmt.Sprintf("go list -json %s", pkgPath)
	result, err := gproc.ShellExec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to execute go list: %v", err)
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
	return nil
}

// filterDeps filters dependencies based on options.
func (a *analyzer) filterDeps(deps []string, in Input) []string {
	result := make([]string, 0)
	for _, dep := range deps {
		if a.shouldInclude(dep, in) {
			result = append(result, dep)
		}
	}
	return result
}

// shouldInclude checks if a dependency should be included.
func (a *analyzer) shouldInclude(dep string, in Input) bool {
	// Exclude standard library
	if in.NoStd && a.isStdLib(dep) {
		return false
	}
	// Only internal packages
	if in.Internal && a.modulePrefix != "" {
		if !gstr.HasPrefix(dep, a.modulePrefix) {
			return false
		}
	}
	return true
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
