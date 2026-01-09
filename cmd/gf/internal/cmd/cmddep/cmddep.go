// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

var (
	Dep = cDep{}
)

type cDep struct {
	g.Meta `name:"dep" brief:"{cDepBrief}" eg:"{cDepEg}"`
}

const (
	cDepBrief = `analyze and display Go package dependencies`
	cDepEg    = `
gf dep
gf dep ./...
gf dep ./internal/...
gf dep -f list
gf dep -f mermaid
gf dep -f mermaid -g
gf dep -f dot -d 5
gf dep -f json -d 0
gf dep -g
gf dep -r
gf dep -i=false
gf dep -e
gf dep -e -i=false
gf dep -m
gf dep -m -e
gf dep -s
gf dep -s -p 8080
gf dep ./internal/... -f tree -d 2
gf dep --external --group -f mermaid
gf dep --main --external -f json
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cDepBrief`: cDepBrief,
		`cDepEg`:    cDepEg,
	})
}

// Input defines the input parameters for dep command.
type Input struct {
	g.Meta   `name:"dep"`
	Package  string `name:"PACKAGE" arg:"true" brief:"package path to analyze, default is ./..." d:"./..."`
	Format   string `name:"format"   short:"f" brief:"output format: tree/list/mermaid/dot/json" d:"tree"`
	Depth    int    `name:"depth"    short:"d" brief:"dependency depth limit, 0 means unlimited" d:"3"`
	Group    bool   `name:"group"    short:"g" brief:"group by top-level directory" d:"false"`
	Internal bool   `name:"internal" short:"i" brief:"show only internal packages" d:"true"`
	External bool   `name:"external" short:"e" brief:"show external packages" d:"false"`
	MainOnly bool   `name:"main"     short:"m" brief:"analyze only main module packages (exclude submodules)" d:"false"`
	NoStd    bool   `name:"nostd"    short:"n" brief:"exclude standard library" d:"true"`
	Reverse  bool   `name:"reverse"  short:"r" brief:"show reverse dependencies" d:"false"`
	Serve    bool   `name:"serve"    short:"s" brief:"start HTTP server to view dependencies" d:"false" orphan:"true"`
	Port     int    `name:"port"     short:"p" brief:"HTTP server port" d:"8888"`
}

// Output defines the output for dep command.
type Output struct{}

// Index is the main entry point for the dep command.
func (c cDep) Index(ctx context.Context, in Input) (out *Output, err error) {
	analyzer := newAnalyzer()

	// Detect module prefix from go.mod
	analyzer.modulePrefix = analyzer.detectModulePrefix()

	// Get package information
	loadErr := analyzer.loadPackages(ctx, in.Package)

	// Start HTTP server if requested
	// In server mode, allow starting even without local Go module
	// because users may want to analyze remote modules
	if in.Serve {
		if loadErr != nil {
			mlog.Print("Warning: No local Go module found, you can analyze remote modules in the web UI")
		}
		return nil, analyzer.startServer(in)
	}

	// For non-server mode, return error if loading failed
	if loadErr != nil {
		return nil, loadErr
	}

	if len(analyzer.packages) == 0 {
		mlog.Print("No packages found")
		return
	}

	// Generate output based on format
	var output string
	if in.Reverse {
		output = analyzer.generateReverse(in)
	} else {
		output = analyzer.generate(in)
	}

	mlog.Print(output)
	return
}
