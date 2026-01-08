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
		return a.generateTree(in)
	}
}

// generateTree generates ASCII tree output.
func (a *analyzer) generateTree(in Input) string {
	var sb strings.Builder
	pkgs := a.getSortedPackages()

	for _, pkgPath := range pkgs {
		pkg := a.packages[pkgPath]
		a.visited = make(map[string]bool)
		shortName := a.shortName(pkg.ImportPath, in.Group)
		sb.WriteString(shortName + "\n")
		a.printTreeNode(&sb, pkg, "", in, 0)
	}
	return sb.String()
}

func (a *analyzer) printTreeNode(sb *strings.Builder, pkg *goPackage, prefix string, in Input, depth int) {
	if in.Depth > 0 && depth >= in.Depth {
		return
	}

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
	allDeps := make(map[string]bool)

	for _, pkg := range a.packages {
		for _, dep := range a.filterDeps(pkg.Imports, in) {
			allDeps[dep] = true
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
	nodes := make([]*depNode, 0)
	for _, pkgPath := range a.getSortedPackages() {
		pkg := a.packages[pkgPath]
		a.visited = make(map[string]bool)
		node := a.buildDepNode(pkg, in, 0)
		nodes = append(nodes, node)
	}

	data, err := json.MarshalIndent(nodes, "", "  ")
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
