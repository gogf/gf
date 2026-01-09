// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// generate creates output based on format.
func (a *analyzer) generate(in Input) string {
	switch in.Format {
	case "tree":
		return a.generateTree(in)
	case "list":
		return a.generateList(in)
	case "mermaid":
		return a.generateMermaid(in)
	case "dot":
		return a.generateDot(in)
	case "json":
		return a.generateJSON(in)
	default:
		// Default to tree format
		return a.generateTree(in)
	}
}

// generateTree generates ASCII tree output.
func (a *analyzer) generateTree(in Input) string {
	var sb strings.Builder

	// Add statistics header if showing external dependencies
	if in.External {
		stats := a.getDependencyStats(in)
		sb.WriteString("Dependency Statistics:\n")
		fmt.Fprintf(&sb, "  Total packages: %v\n", stats["total"])
		fmt.Fprintf(&sb, "  Internal: %v\n", stats["internal"])
		fmt.Fprintf(&sb, "  External: %v\n", stats["external"])
		fmt.Fprintf(&sb, "  Standard library: %v\n", stats["stdlib"])
		
		if groups, ok := stats["external_groups"].(map[string]int); ok && len(groups) > 0 {
			sb.WriteString("  External groups:\n")
			for group, count := range groups {
				fmt.Fprintf(&sb, "    %s: %d\n", group, count)
			}
		}
		sb.WriteString("\nDependency Tree:\n")
	}

	// Find root packages (packages that are not imported by any other package)
	rootPkgs := a.findRootPackages()

	// Use a single visited map across all root packages to avoid duplicates
	a.visited = make(map[string]bool)

	for _, pkgPath := range rootPkgs {
		pkg := a.packages[pkgPath]
		if a.shouldInclude(pkg.ImportPath, in) {
			shortName := a.shortName(pkg.ImportPath, in.Group)
			sb.WriteString(shortName + "\n")
			a.printTreeNode(&sb, pkg, "", in, 0)
		}
	}
	return sb.String()
}

// findRootPackages finds packages that are not imported by any other internal package.
func (a *analyzer) findRootPackages() []string {
	// Build a set of all imported packages
	imported := make(map[string]bool)
	for _, pkg := range a.packages {
		for _, dep := range pkg.Imports {
			imported[dep] = true
		}
	}

	// Find packages that are not imported by others
	roots := make([]string, 0)
	for pkgPath := range a.packages {
		if !imported[pkgPath] {
			roots = append(roots, pkgPath)
		}
	}

	// If no roots found (circular dependencies), use all packages
	if len(roots) == 0 {
		roots = a.getSortedPackages()
	}

	sort.Strings(roots)
	return roots
}

func (a *analyzer) printTreeNode(sb *strings.Builder, pkg *goPackage, prefix string, in Input, depth int) {
	if in.Depth > 0 && depth >= in.Depth {
		return
	}

	// filterDeps already applies all filtering including main-only
	deps := a.filterDeps(pkg.Imports, in)
	sort.Strings(deps)

	for i, dep := range deps {
		if a.visited[dep] {
			continue
		}
		a.visited[dep] = true

		isLast := i == len(deps)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		shortName := a.shortName(dep, in.Group)
		sb.WriteString(prefix + connector + shortName + "\n")

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		// Recursively print dependencies
		if depPkg, ok := a.packages[dep]; ok {
			a.printTreeNode(sb, depPkg, newPrefix, in, depth+1)
		}
	}
}

// generateList generates simple list output.
func (a *analyzer) generateList(in Input) string {
	var sb strings.Builder
	
	// Add statistics header if showing external dependencies
	if in.External {
		stats := a.getDependencyStats(in)
		sb.WriteString("# Dependency Statistics\n")
		fmt.Fprintf(&sb, "# Total: %v, Internal: %v, External: %v, Stdlib: %v\n", 
			stats["total"], stats["internal"], stats["external"], stats["stdlib"])
		sb.WriteString("\n")
	}

	// Debug mainOnly state
	// sb.WriteString(fmt.Sprintf("# DEBUG mainOnly=%v\\n", in.MainOnly))
	
	allDeps := make(map[string]bool)

	// Collect dependencies from packages that should be included
	for _, pkg := range a.packages {
		if a.shouldInclude(pkg.ImportPath, in) {
			// Collect dependencies (filterDeps already applies all filtering including main-only)
			for _, dep := range a.filterDeps(pkg.Imports, in) {
				allDeps[dep] = true
			}
		}
		
		// Additionally, keep the package itself when it passes main-only check
		if in.MainOnly && a.isModuleRootPackage(pkg.ImportPath) && a.shouldInclude(pkg.ImportPath, in) {
			allDeps[pkg.ImportPath] = true
		}
	}


	deps := make([]string, 0, len(allDeps))
	for dep := range allDeps {
		deps = append(deps, a.shortName(dep, in.Group))
	}
	sort.Strings(deps)

	for _, dep := range deps {
		sb.WriteString(dep + "\n")
	}
	return sb.String()
}

// generateMermaid generates Mermaid diagram output.
func (a *analyzer) generateMermaid(in Input) string {
	var sb strings.Builder
	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TD\n")

	edges := a.collectEdges(in)
	sortedEdges := make([]string, 0, len(edges))
	for edge := range edges {
		sortedEdges = append(sortedEdges, edge)
	}
	sort.Strings(sortedEdges)

	for _, edge := range sortedEdges {
		sb.WriteString("    " + edge + "\n")
	}
	sb.WriteString("```\n")
	return sb.String()
}

// generateMermaidRaw generates Mermaid code without markdown wrapper.
func (a *analyzer) generateMermaidRaw(in Input) string {
	var sb strings.Builder
	sb.WriteString("graph TD\n")

	edges := a.collectEdges(in)
	sortedEdges := make([]string, 0, len(edges))
	for edge := range edges {
		sortedEdges = append(sortedEdges, edge)
	}
	sort.Strings(sortedEdges)

	for _, edge := range sortedEdges {
		sb.WriteString("    " + edge + "\n")
	}
	return sb.String()
}

// generateDot generates Graphviz DOT output.
func (a *analyzer) generateDot(in Input) string {
	var sb strings.Builder
	sb.WriteString("digraph deps {\n")
	sb.WriteString("    rankdir=TB;\n")
	sb.WriteString("    node [shape=box];\n")

	edges := a.collectEdges(in)
	sortedEdges := make([]string, 0, len(edges))
	for edge := range edges {
		sortedEdges = append(sortedEdges, edge)
	}
	sort.Strings(sortedEdges)

	for _, edge := range sortedEdges {
		parts := strings.Split(edge, " --> ")
		if len(parts) == 2 {
			fmt.Fprintf(&sb, "    \"%s\" -> \"%s\";\n", parts[0], parts[1])
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

// generateJSON generates JSON output.
func (a *analyzer) generateJSON(in Input) string {
	result := make(map[string]any)
	
	// Add dependency nodes
	nodes := make([]*depNode, 0)
	for _, pkgPath := range a.getSortedPackages() {
		pkg := a.packages[pkgPath]
		if a.shouldInclude(pkg.ImportPath, in) {
			a.visited = make(map[string]bool)
			node := a.buildDepNode(pkg, in, 0)
			nodes = append(nodes, node)
		}
	}
	result["dependencies"] = nodes
	
	// Add statistics
	result["statistics"] = a.getDependencyStats(in)
	
	// Add metadata
	result["metadata"] = map[string]any{
		"module":   a.modulePrefix,
		"format":   in.Format,
		"depth":    in.Depth,
		"group":    in.Group,
		"internal": in.Internal,
		"external": in.External,
		"nostd":    in.NoStd,
		"main":     in.MainOnly,
	}
	
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(data)
}

func (a *analyzer) buildDepNode(pkg *goPackage, in Input, depth int) *depNode {
	node := &depNode{
		Package: a.shortName(pkg.ImportPath, in.Group),
	}

	if in.Depth > 0 && depth >= in.Depth {
		return node
	}

	deps := a.filterDeps(pkg.Imports, in)
	sort.Strings(deps)

	for _, dep := range deps {
		if a.visited[dep] {
			continue
		}
		a.visited[dep] = true

		if depPkg, ok := a.packages[dep]; ok {
			childNode := a.buildDepNode(depPkg, in, depth+1)
			node.Dependencies = append(node.Dependencies, childNode)
		} else {
			node.Dependencies = append(node.Dependencies, &depNode{
				Package: a.shortName(dep, in.Group),
			})
		}
	}
	return node
}

// generateReverse generates reverse dependency output.
func (a *analyzer) generateReverse(in Input) string {
	// Build reverse dependency map
	reverseDeps := make(map[string][]string)
	for pkgPath, pkg := range a.packages {
		for _, dep := range pkg.Imports {
			if a.shouldInclude(dep, in) {
				reverseDeps[dep] = append(reverseDeps[dep], pkgPath)
			}
		}
	}

	var sb strings.Builder
	targets := a.getSortedPackages()

	for _, target := range targets {
		deps := reverseDeps[target]
		if len(deps) == 0 {
			continue
		}

		sort.Strings(deps)
		shortTarget := a.shortName(target, in.Group)
		if shortTarget == "" {
			continue
		}
		fmt.Fprintf(&sb, "%s (used by %d packages):\n", shortTarget, len(deps))

		for i, dep := range deps {
			isLast := i == len(deps)-1
			connector := "├── "
			if isLast {
				connector = "└── "
			}
			sb.WriteString(connector + a.shortName(dep, in.Group) + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
