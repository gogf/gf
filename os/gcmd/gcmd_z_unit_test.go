// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gcmd.Init([]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}...)
		t.Assert(len(gcmd.GetArgAll()), 2)
		t.Assert(gcmd.GetArg(1), "path")
		t.Assert(gcmd.GetArg(100, "test"), "test")
		t.Assert(gcmd.GetOpt("force"), "remove")
		t.Assert(gcmd.GetOpt("n"), "root")
		t.Assert(gcmd.GetOpt("fq").IsNil(), false)
		t.Assert(gcmd.GetOpt("p").IsNil(), false)
		t.Assert(gcmd.GetOpt("none").IsNil(), true)
		t.Assert(gcmd.GetOpt("none", "value"), "value")
	})
	gtest.C(t, func(t *gtest.T) {
		gcmd.Init([]string{"gf", "gen", "-h"}...)
		t.Assert(len(gcmd.GetArgAll()), 2)
		t.Assert(gcmd.GetOpt("h"), "")
		t.Assert(gcmd.GetOpt("h").IsNil(), false)
	})
}

func Test_BuildOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		})
		t.Assert(s, "-n=john")
	})

	gtest.C(t, func(t *gtest.T) {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		}, "-test")
		t.Assert(s, "-testn=john")
	})

	gtest.C(t, func(t *gtest.T) {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n1": "john",
			"n2": "huang",
		})
		t.Assert(strings.Contains(s, "-n1=john"), true)
		t.Assert(strings.Contains(s, "-n2=huang"), true)
	})
}

func Test_GetWithEnv(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		t.Assert(gcmd.GetOptWithEnv("test"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		gcmd.Init("-test", "2")
		t.Assert(gcmd.GetOptWithEnv("test"), 2)
	})
}

func Test_Command(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = gctx.New()
			err error
		)
		commandRoot := &gcmd.Command{
			Name: "gf",
		}
		// env
		commandEnv := &gcmd.Command{
			Name: "env",
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				fmt.Println("env")
				return nil
			},
		}
		// test
		commandTest := &gcmd.Command{
			Name:        "test",
			Brief:       "test brief",
			Description: "test description current Golang environment variables",
			Examples: `
gf get github.com/gogf/gf
gf get github.com/gogf/gf@latest
gf get github.com/gogf/gf@master
gf get golang.org/x/sys
`,
			Arguments: []gcmd.Argument{
				{
					Name:   "my-option",
					Short:  "o",
					Brief:  "It's my custom option",
					Orphan: true,
				},
				{
					Name:   "another",
					Short:  "a",
					Brief:  "It's my another custom option",
					Orphan: true,
				},
			},
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				fmt.Println("test")
				return nil
			},
		}
		err = commandRoot.AddCommand(
			commandEnv,
		)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		err = commandRoot.AddObject(
			commandTest,
		)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}

		if err = commandRoot.RunWithError(ctx); err != nil {
			if gerror.Code(err) == gcode.CodeNotFound {
				commandRoot.Print()
			}
		}
	})
}

func Test_Command_Print(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = gctx.New()
			err error
		)
		c := &gcmd.Command{
			Name:        "gf",
			Description: `GoFrame Command Line Interface, which is your helpmate for building GoFrame application with convenience.`,
			Additional: `
Use 'gf help COMMAND' or 'gf COMMAND -h' for detail about a command, which has '...' in the tail of their comments.`,
		}
		// env
		commandEnv := &gcmd.Command{
			Name:        "env",
			Brief:       "show current Golang environment variables, long brief.long brief.long brief.long brief.long brief.long brief.long brief.long brief.",
			Description: "show current Golang environment variables",
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(commandEnv); err != nil {
			g.Log().Fatal(ctx, err)
		}
		// get
		commandGet := &gcmd.Command{
			Name:        "get",
			Brief:       "install or update GF to system in default...",
			Description: "show current Golang environment variables",

			Examples: `
gf get github.com/gogf/gf
gf get github.com/gogf/gf@latest
gf get github.com/gogf/gf@master
gf get golang.org/x/sys
`,
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(commandGet); err != nil {
			g.Log().Fatal(ctx, err)
		}
		// build
		//-n, --name       output binary name
		//-v, --version    output binary version
		//-a, --arch       output binary architecture, multiple arch separated with ','
		//-s, --system     output binary system, multiple os separated with ','
		//-o, --output     output binary path, used when building single binary file
		//-p, --path       output binary directory path, default is './bin'
		//-e, --extra      extra custom "go build" options
		//-m, --mod        like "-mod" option of "go build", use "-m none" to disable go module
		//-c, --cgo        enable or disable cgo feature, it's disabled in default

		commandBuild := gcmd.Command{
			Name:  "build",
			Usage: "gf build FILE [OPTION]",
			Brief: "cross-building go project for lots of platforms...",
			Description: `
The "build" command is most commonly used command, which is designed as a powerful wrapper for
"go build" command for convenience cross-compiling usage.
It provides much more features for building binary:
1. Cross-Compiling for many platforms and architectures.
2. Configuration file support for compiling.
3. Build-In Variables.
`,
			Examples: `
gf build main.go
gf build main.go --swagger
gf build main.go --pack public,template
gf build main.go --cgo
gf build main.go -m none 
gf build main.go -n my-app -a all -s all
gf build main.go -n my-app -a amd64,386 -s linux -p .
gf build main.go -n my-app -v 1.0 -a amd64,386 -s linux,windows,darwin -p ./docker/bin
`,
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(&commandBuild); err != nil {
			g.Log().Fatal(ctx, err)
		}
		_ = c.RunWithError(ctx)
	})
}

func Test_Command_NotFound(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c0 := &gcmd.Command{
			Name: "c0",
		}
		c1 := &gcmd.Command{
			Name: "c1",
			FuncWithValue: func(ctx context.Context, parser *gcmd.Parser) (any, error) {
				return nil, nil
			},
		}
		c21 := &gcmd.Command{
			Name: "c21",
			FuncWithValue: func(ctx context.Context, parser *gcmd.Parser) (any, error) {
				return nil, nil
			},
		}
		c22 := &gcmd.Command{
			Name: "c22",
			FuncWithValue: func(ctx context.Context, parser *gcmd.Parser) (any, error) {
				return nil, nil
			},
		}
		t.AssertNil(c0.AddCommand(c1))
		t.AssertNil(c1.AddCommand(c21, c22))

		os.Args = []string{"c0", "c1", "c23", `--test="abc"`}
		err := c0.RunWithError(gctx.New())
		t.Assert(err.Error(), `command "c1 c23" not found for command "c0", command line: c0 c1 c23 --test="abc"`)
	})
}
