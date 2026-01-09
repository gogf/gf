// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

//go:embed static/*
var staticFiles embed.FS

// graphData represents the graph structure for visualization.
type graphData struct {
	Nodes []graphNode `json:"nodes"`
	Edges []graphEdge `json:"edges"`
}

type graphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Group string `json:"group,omitempty"`
}

type graphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// packageInfo represents package information for API response.
type packageInfo struct {
	Name         string   `json:"name"`
	FullPath     string   `json:"fullPath"`
	Dependencies []string `json:"dependencies"`
	UsedBy       []string `json:"usedBy"`
}

// packageSummary represents package summary for list API response.
type packageSummary struct {
	Name        string `json:"name"`
	DepCount    int    `json:"depCount"`
	UsedByCount int    `json:"usedByCount"`
}

// moduleInfo represents module information for API response.
type moduleInfo struct {
	Name string `json:"name"`
}

// versionInfo represents version list response.
type versionInfo struct {
	Versions []string `json:"versions,omitempty"`
	Error    string   `json:"error,omitempty"`
}

// analyzeResult represents analyze result response.
type analyzeResult struct {
	Success bool   `json:"success"`
	Module  string `json:"module,omitempty"`
	Error   string `json:"error,omitempty"`
}

// serverState holds the server state for remote module analysis.
type serverState struct {
	originalAnalyzer *analyzer
	currentAnalyzer  *analyzer
	originalInput    Input
	tempDir          string
}

// startServer starts an HTTP server to visualize dependencies.
func (a *analyzer) startServer(in Input) error {
	addr := fmt.Sprintf(":%d", in.Port)

	// Create server state
	state := &serverState{
		originalAnalyzer: a,
		currentAnalyzer:  a,
		originalInput:    in,
	}

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})

	// API endpoints
	http.HandleFunc("/api/module", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handleModuleAPI(w)
	})
	http.HandleFunc("/api/graph", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handleGraphAPI(w, r, in)
	})
	http.HandleFunc("/api/packages", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handlePackagesAPI(w, r, in)
	})
	http.HandleFunc("/api/package", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handlePackageAPI(w, r, in)
	})
	http.HandleFunc("/api/tree", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handleTreeAPI(w, r, in)
	})
	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		state.currentAnalyzer.handleListAPI(w, r, in)
	})
	http.HandleFunc("/api/versions", func(w http.ResponseWriter, r *http.Request) {
		handleVersionsAPI(w, r)
	})
	http.HandleFunc("/api/analyze", func(w http.ResponseWriter, r *http.Request) {
		state.handleAnalyzeAPI(w, r)
	})
	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		state.handleResetAPI(w)
	})

	mlog.Printf("Starting dependency viewer at http://localhost%s", addr)
	mlog.Print("Press Ctrl+C to stop")

	return http.ListenAndServe(addr, nil)
}

// handleModuleAPI returns module information.
func (a *analyzer) handleModuleAPI(w http.ResponseWriter) {
	info := moduleInfo{
		Name: a.modulePrefix,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// handleGraphAPI returns graph data as JSON.
func (a *analyzer) handleGraphAPI(w http.ResponseWriter, r *http.Request, in Input) {
	query := r.URL.Query()
	if g := query.Get("group"); g != "" {
		in.Group = g == "true"
	}
	if d := query.Get("depth"); d != "" {
		fmt.Sscanf(d, "%d", &in.Depth)
	}
	if rev := query.Get("reverse"); rev != "" {
		in.Reverse = rev == "true"
	}
	if i := query.Get("internal"); i != "" {
		in.Internal = i == "true"
	}
	if e := query.Get("external"); e != "" {
		in.External = e == "true"
	}
	if m := query.Get("main"); m != "" {
		in.MainOnly = m == "true"
	}
	pkg := query.Get("package")

	var data *graphData
	if pkg != "" {
		data = a.buildPackageGraphData(pkg, in)
	} else {
		data = a.buildGraphData(in)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handlePackagesAPI returns all packages list with dependency stats.
func (a *analyzer) handlePackagesAPI(w http.ResponseWriter, r *http.Request, in Input) {
	query := r.URL.Query()
	if i := query.Get("internal"); i != "" {
		in.Internal = i == "true"
	}
	if e := query.Get("external"); e != "" {
		in.External = e == "true"
	}
	if m := query.Get("main"); m != "" {
		in.MainOnly = m == "true"
	}

	// Build reverse dependency map (who uses each package)
	usedByMap := make(map[string]int)
	for fullPath, pkg := range a.packages {
		if !a.shouldInclude(fullPath, in) {
			continue
		}
		fromShort := a.shortName(fullPath, false)
		if fromShort == "" {
			continue
		}
		for _, dep := range a.filterDeps(pkg.Imports, in) {
			shortDep := a.shortName(dep, false)
			if shortDep != "" {
				usedByMap[shortDep]++
			}
		}
	}

	packages := make([]packageSummary, 0)
	for _, pkgPath := range a.getSortedPackages() {
		if !a.shouldInclude(pkgPath, in) {
			continue
		}
		shortName := a.shortName(pkgPath, false)
		if shortName == "" {
			continue
		}

		// Count dependencies (filtered)
		depCount := 0
		if pkg, ok := a.packages[pkgPath]; ok {
			for _, dep := range a.filterDeps(pkg.Imports, in) {
				shortDep := a.shortName(dep, false)
				if shortDep != "" {
					depCount++
				}
			}
		}

		packages = append(packages, packageSummary{
			Name:        shortName,
			DepCount:    depCount,
			UsedByCount: usedByMap[shortName],
		})
	}
	
	// Add statistics to response
	result := map[string]any{
		"packages":   packages,
		"statistics": a.getDependencyStats(in),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handlePackageAPI returns detailed info for a specific package.
func (a *analyzer) handlePackageAPI(w http.ResponseWriter, r *http.Request, in Input) {
	query := r.URL.Query()
	pkgName := query.Get("name")
	if pkgName == "" {
		http.Error(w, "package name required", http.StatusBadRequest)
		return
	}

	// Find the full package path
	var fullPath string
	for path := range a.packages {
		if a.shortName(path, false) == pkgName {
			fullPath = path
			break
		}
	}

	if fullPath == "" {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}

	pkg := a.packages[fullPath]
	info := packageInfo{
		Name:         pkgName,
		FullPath:     fullPath,
		Dependencies: make([]string, 0),
		UsedBy:       make([]string, 0),
	}

	// Get dependencies
	for _, dep := range a.filterDeps(pkg.Imports, in) {
		shortName := a.shortName(dep, false)
		if shortName != "" {
			info.Dependencies = append(info.Dependencies, shortName)
		}
	}
	sort.Strings(info.Dependencies)

	// Get reverse dependencies (who uses this package)
	for path, p := range a.packages {
		for _, dep := range p.Imports {
			if dep == fullPath {
				shortName := a.shortName(path, false)
				if shortName != "" {
					info.UsedBy = append(info.UsedBy, shortName)
				}
				break
			}
		}
	}
	sort.Strings(info.UsedBy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// handleTreeAPI returns tree format output.
func (a *analyzer) handleTreeAPI(w http.ResponseWriter, r *http.Request, in Input) {
	query := r.URL.Query()
	if d := query.Get("depth"); d != "" {
		fmt.Sscanf(d, "%d", &in.Depth)
	}
	if i := query.Get("internal"); i != "" {
		in.Internal = i == "true"
	}
	if e := query.Get("external"); e != "" {
		in.External = e == "true"
	}
	if m := query.Get("main"); m != "" {
		in.MainOnly = m == "true"
	}
	pkg := query.Get("package")

	var output string
	if pkg != "" {
		output = a.generatePackageTree(pkg, in)
	} else {
		output = a.generateTree(in)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(output))
}

// handleListAPI returns list format output.
func (a *analyzer) handleListAPI(w http.ResponseWriter, r *http.Request, in Input) {
	query := r.URL.Query()
	if i := query.Get("internal"); i != "" {
		in.Internal = i == "true"
	}
	if e := query.Get("external"); e != "" {
		in.External = e == "true"
	}
	if m := query.Get("main"); m != "" {
		in.MainOnly = m == "true"
	}
	pkg := query.Get("package")

	var output string
	if pkg != "" {
		output = a.generatePackageList(pkg, in)
	} else {
		output = a.generateList(in)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(output))
}

// buildGraphData builds graph data for visualization.
func (a *analyzer) buildGraphData(in Input) *graphData {
	data := &graphData{
		Nodes: make([]graphNode, 0),
		Edges: make([]graphEdge, 0),
	}

	nodeSet := make(map[string]bool)
	edges := a.collectEdges(in)

	for edge := range edges {
		parts := strings.Split(edge, " --> ")
		if len(parts) != 2 {
			continue
		}
		from, to := parts[0], parts[1]

		if !nodeSet[from] {
			nodeSet[from] = true
			data.Nodes = append(data.Nodes, graphNode{
				ID:    from,
				Label: strings.ReplaceAll(from, "_", "/"),
				Group: a.getNodeGroup(from),
			})
		}
		if !nodeSet[to] {
			nodeSet[to] = true
			data.Nodes = append(data.Nodes, graphNode{
				ID:    to,
				Label: strings.ReplaceAll(to, "_", "/"),
				Group: a.getNodeGroup(to),
			})
		}

		data.Edges = append(data.Edges, graphEdge{From: from, To: to})
	}

	return data
}

// buildPackageGraphData builds graph data for a specific package.
func (a *analyzer) buildPackageGraphData(pkgName string, in Input) *graphData {
	data := &graphData{
		Nodes: make([]graphNode, 0),
		Edges: make([]graphEdge, 0),
	}

	// Find the full package path
	var fullPath string
	for path := range a.packages {
		if a.shortName(path, false) == pkgName {
			fullPath = path
			break
		}
	}

	if fullPath == "" {
		return data
	}

	nodeSet := make(map[string]bool)
	nodeSet[pkgName] = true
	data.Nodes = append(data.Nodes, graphNode{
		ID:    a.sanitizeName(pkgName),
		Label: pkgName,
		Group: a.getNodeGroup(pkgName),
	})

	pkg := a.packages[fullPath]

	if in.Reverse {
		// Show packages that depend on this package
		for path, p := range a.packages {
			for _, dep := range p.Imports {
				if dep == fullPath {
					shortName := a.shortName(path, false)
					if shortName != "" && !nodeSet[shortName] {
						nodeSet[shortName] = true
						data.Nodes = append(data.Nodes, graphNode{
							ID:    a.sanitizeName(shortName),
							Label: shortName,
							Group: a.getNodeGroup(shortName),
						})
						data.Edges = append(data.Edges, graphEdge{
							From: a.sanitizeName(shortName),
							To:   a.sanitizeName(pkgName),
						})
					}
					break
				}
			}
		}
	} else {
		// Show dependencies of this package
		a.collectPackageDeps(pkg, pkgName, in, nodeSet, data, 0)
	}

	return data
}

// collectPackageDeps recursively collects dependencies for a package.
func (a *analyzer) collectPackageDeps(pkg *goPackage, pkgName string, in Input, nodeSet map[string]bool, data *graphData, depth int) {
	if in.Depth > 0 && depth >= in.Depth {
		return
	}

	deps := a.filterDeps(pkg.Imports, in)
	for _, dep := range deps {
		shortName := a.shortName(dep, false)
		if shortName == "" {
			continue
		}

		data.Edges = append(data.Edges, graphEdge{
			From: a.sanitizeName(pkgName),
			To:   a.sanitizeName(shortName),
		})

		if !nodeSet[shortName] {
			nodeSet[shortName] = true
			data.Nodes = append(data.Nodes, graphNode{
				ID:    a.sanitizeName(shortName),
				Label: shortName,
				Group: a.getNodeGroup(shortName),
			})

			// Recursively collect dependencies
			if depPkg, ok := a.packages[dep]; ok {
				a.collectPackageDeps(depPkg, shortName, in, nodeSet, data, depth+1)
			}
		}
	}
}

// generatePackageTree generates tree output for a specific package.
func (a *analyzer) generatePackageTree(pkgName string, in Input) string {
	var fullPath string
	for path := range a.packages {
		if a.shortName(path, false) == pkgName {
			fullPath = path
			break
		}
	}

	if fullPath == "" {
		return "Package not found: " + pkgName
	}

	var sb strings.Builder
	pkg := a.packages[fullPath]
	a.visited = make(map[string]bool)
	sb.WriteString(pkgName + "\n")
	a.printTreeNode(&sb, pkg, "", in, 0)
	return sb.String()
}

// generatePackageList generates list output for a specific package.
func (a *analyzer) generatePackageList(pkgName string, in Input) string {
	var fullPath string
	for path := range a.packages {
		if a.shortName(path, false) == pkgName {
			fullPath = path
			break
		}
	}

	if fullPath == "" {
		return "Package not found: " + pkgName
	}

	var sb strings.Builder
	pkg := a.packages[fullPath]
	deps := a.filterDeps(pkg.Imports, in)

	shortDeps := make([]string, 0, len(deps))
	for _, dep := range deps {
		shortName := a.shortName(dep, false)
		if shortName != "" {
			shortDeps = append(shortDeps, shortName)
		}
	}
	sort.Strings(shortDeps)

	for _, dep := range shortDeps {
		sb.WriteString(dep + "\n")
	}
	return sb.String()
}

// getNodeGroup returns the group (top-level directory) of a node.
func (a *analyzer) getNodeGroup(name string) string {
	name = strings.ReplaceAll(name, "_", "/")
	parts := strings.Split(name, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// handleVersionsAPI fetches available versions for a Go module from proxy.
func handleVersionsAPI(w http.ResponseWriter, r *http.Request) {
	modulePath := r.URL.Query().Get("module")
	if modulePath == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionInfo{Error: "module parameter required"})
		return
	}

	// Fetch versions from Go proxy
	versions, err := fetchModuleVersions(modulePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionInfo{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionInfo{Versions: versions})
}

// fetchModuleVersions fetches versions from Go proxy.
func fetchModuleVersions(modulePath string) ([]string, error) {
	// Create a temp directory with go.mod to run go list -m
	tempDir, err := os.MkdirTemp("", "gf-dep-versions-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a temp module using exec.Command with Dir (cross-platform)
	initCmd := exec.Command("go", "mod", "init", "temp")
	initCmd.Dir = tempDir
	if output, err := initCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to init temp module: %v, output: %s", err, string(output))
	}

	// Use go list to get available versions in temp directory
	listCmd := exec.Command("go", "list", "-m", "-versions", modulePath)
	listCmd.Dir = tempDir
	output, err := listCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %v, output: %s", err, string(output))
	}
	result := string(output)

	// Parse output: module@version version1 version2 ...
	result = strings.TrimSpace(result)
	if result == "" {
		return nil, fmt.Errorf("no versions found")
	}

	parts := strings.Fields(result)
	if len(parts) < 2 {
		// Only module name, try to get latest
		return []string{"latest"}, nil
	}

	// Reverse order (newest first)
	versions := parts[1:]
	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}

	return versions, nil
}

// handleAnalyzeAPI analyzes a remote module.
func (s *serverState) handleAnalyzeAPI(w http.ResponseWriter, r *http.Request) {
	modulePath := r.URL.Query().Get("module")
	version := r.URL.Query().Get("version")

	if modulePath == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: "module parameter required"})
		return
	}

	// Clean up previous temp directory
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gf-dep-*")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: "failed to create temp directory"})
		return
	}
	s.tempDir = tempDir

	// Download and analyze module
	moduleWithVersion := modulePath
	if version != "" && version != "latest" {
		moduleWithVersion = modulePath + "@" + version
	}

	// Initialize go module in temp directory (cross-platform)
	initCmd := exec.Command("go", "mod", "init", "temp")
	initCmd.Dir = tempDir
	if output, err := initCmd.CombinedOutput(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: fmt.Sprintf("failed to init module: %v, output: %s", err, string(output))})
		return
	}

	// Download the module (cross-platform)
	getCmd := exec.Command("go", "get", moduleWithVersion)
	getCmd.Dir = tempDir
	if output, err := getCmd.CombinedOutput(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: fmt.Sprintf("failed to download module: %v, output: %s", err, string(output))})
		return
	}

	// Find the module in GOPATH/pkg/mod
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, _ := os.UserHomeDir()
		gopath = filepath.Join(home, "go")
	}

	// Find the actual module directory
	modCacheDir := filepath.Join(gopath, "pkg", "mod")
	moduleDir, err := findModuleDir(modCacheDir, modulePath, version)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: "failed to find module: " + err.Error()})
		return
	}

	// Create new analyzer for the remote module
	newAnalyzer := newAnalyzer()
	newAnalyzer.modulePrefix = modulePath

	// Load packages from the module directory (cross-platform)
	// IMPORTANT: Must run in tempDir context where the module was downloaded,
	// otherwise it will use packages from the current project's dependencies
	listCmd := exec.Command("go", "list", "-json", modulePath+"/...")
	listCmd.Dir = tempDir
	output, err := listCmd.CombinedOutput()
	result := string(output)
	if err != nil {
		// Try loading from the module directory directly
		listCmd2 := exec.Command("go", "list", "-json", "./...")
		listCmd2.Dir = moduleDir
		output2, err2 := listCmd2.CombinedOutput()
		if err2 != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(analyzeResult{Error: fmt.Sprintf("failed to list packages: %v, output: %s", err2, string(output2))})
			return
		}
		result = string(output2)
	}

	// Parse packages
	decoder := json.NewDecoder(strings.NewReader(result))
	for decoder.More() {
		var pkg goPackage
		if err := decoder.Decode(&pkg); err != nil {
			continue
		}
		newAnalyzer.packages[pkg.ImportPath] = &pkg
	}

	if len(newAnalyzer.packages) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analyzeResult{Error: "no packages found in module"})
		return
	}

	s.currentAnalyzer = newAnalyzer

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analyzeResult{
		Success: true,
		Module:  moduleWithVersion,
	})
}

// findModuleDir finds the module directory in the module cache.
func findModuleDir(modCacheDir, modulePath, version string) (string, error) {
	// Convert module path to filesystem path
	escapedPath := strings.ReplaceAll(modulePath, "/", string(filepath.Separator))

	// Handle uppercase letters in module path (they're escaped in the cache)
	var escapedParts []string
	for _, part := range strings.Split(escapedPath, string(filepath.Separator)) {
		var escaped strings.Builder
		for _, c := range part {
			if c >= 'A' && c <= 'Z' {
				escaped.WriteRune('!')
				escaped.WriteRune(c + 32) // lowercase
			} else {
				escaped.WriteRune(c)
			}
		}
		escapedParts = append(escapedParts, escaped.String())
	}
	escapedPath = strings.Join(escapedParts, string(filepath.Separator))

	baseDir := filepath.Join(modCacheDir, escapedPath)

	// If version specified, look for exact match
	if version != "" && version != "latest" {
		versionDir := baseDir + "@" + version
		if _, err := os.Stat(versionDir); err == nil {
			return versionDir, nil
		}
	}

	// Find latest version
	parent := filepath.Dir(baseDir)
	base := filepath.Base(baseDir)

	entries, err := os.ReadDir(parent)
	if err != nil {
		return "", err
	}

	var latestDir string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), base+"@") {
			latestDir = filepath.Join(parent, entry.Name())
		}
	}

	if latestDir == "" {
		return "", fmt.Errorf("module not found in cache")
	}

	return latestDir, nil
}

// handleResetAPI resets to the original local analyzer.
func (s *serverState) handleResetAPI(w http.ResponseWriter) {
	// Clean up temp directory
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
		s.tempDir = ""
	}

	s.currentAnalyzer = s.originalAnalyzer

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
