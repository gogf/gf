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
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

type TestCmdObject struct{}

type TestCmdObjectInput struct {
	g.Meta `name:"gf" usage:"gf env/test" brief:"gf env command" dc:"description" ad:"ad"`
}
type TestCmdObjectOutput struct{}

type TestCmdObjectEnvInput struct {
	g.Meta `name:"env" usage:"gf env/test" brief:"gf env command" dc:"description" ad:"ad"`
	Name   string `v:"required" name:"name" short:"n" orphan:"false" brief:"name for command"`
}
type TestCmdObjectEnvOutput struct{}

type TestCmdObjectTestInput struct {
	g.Meta `name:"test" usage:"gf env/test" brief:"gf test command" dc:"description" ad:"ad"`
	Name   string `v:"required" name:"name" short:"n" orphan:"false" brief:"name for command"`
}
type TestCmdObjectTestOutput struct{}

func (TestCmdObject) Root(ctx context.Context, in TestCmdObjectInput) (out *TestCmdObjectOutput, err error) {
	return
}

func (TestCmdObject) Env(ctx context.Context, in TestCmdObjectEnvInput) (out *TestCmdObjectEnvOutput, err error) {
	return
}

func (TestCmdObject) Test(ctx context.Context, in TestCmdObjectTestInput) (out *TestCmdObjectTestOutput, err error) {
	return
}

func Test_Command_AddObject(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = gctx.New()
			err error
		)
		commandRoot := &gcmd.Command{
			Name: "gf",
		}
		// env
		commandEnv := gcmd.Command{
			Name: "env",
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				fmt.Println("env")
				return nil
			},
		}
		// test
		commandTest := gcmd.Command{
			Name:        "test",
			Brief:       "test brief",
			Description: "test description current Golang environment variables",
			Examples: `
gf get github.com/gogf/gf
gf get github.com/gogf/gf@latest
gf get github.com/gogf/gf@master
gf get golang.org/x/sys
`,
			Options: []gcmd.Option{
				{
					Name:   "my-option",
					Short:  "o",
					Brief:  "It's my custom option",
					Orphan: false,
				},
				{
					Name:   "another",
					Short:  "a",
					Brief:  "It's my another custom option",
					Orphan: false,
				},
			},
			Func: func(ctx context.Context, parser *gcmd.Parser) error {
				fmt.Println("test")
				return nil
			},
		}
		err = commandRoot.AddCommand(
			commandEnv,
			commandTest,
		)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}

		if err = commandRoot.Run(ctx); err != nil {
			g.Log().Fatal(ctx, err)
		}
	})
}
