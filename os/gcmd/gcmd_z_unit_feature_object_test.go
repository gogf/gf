// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

type TestCmdObject struct {
	g.Meta `name:"root" usage:"root env/test" brief:"root env command" dc:"description" ad:"ad"`
}

type TestCmdObjectEnvInput struct {
	g.Meta `name:"env" usage:"root env" brief:"root env command" dc:"root env command description" ad:"root env command ad"`
}
type TestCmdObjectEnvOutput struct{}

type TestCmdObjectTestInput struct {
	g.Meta `name:"test" usage:"root test" brief:"root test command" dc:"root test command description" ad:"root test command ad"`
	Name   string `v:"required" short:"n" orphan:"false" brief:"name for test command"`
}
type TestCmdObjectTestOutput struct {
	Content string
}

func (TestCmdObject) Env(ctx context.Context, in TestCmdObjectEnvInput) (out *TestCmdObjectEnvOutput, err error) {
	return
}

func (TestCmdObject) Test(ctx context.Context, in TestCmdObjectTestInput) (out *TestCmdObjectTestOutput, err error) {
	out = &TestCmdObjectTestOutput{
		Content: in.Name,
	}
	return
}

func Test_Command_NewFromObject_Help(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx      = gctx.New()
			cmd, err = gcmd.NewFromObject(&TestCmdObject{})
		)
		t.AssertNil(err)
		t.Assert(cmd.Name, "root")

		os.Args = []string{"root"}
		value, err := cmd.RunWithValue(ctx)
		t.AssertNil(err)
		t.Assert(value, nil)
	})
}

func Test_Command_NewFromObject_RunWithValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx      = gctx.New()
			cmd, err = gcmd.NewFromObject(&TestCmdObject{})
		)
		t.AssertNil(err)
		t.Assert(cmd.Name, "root")

		os.Args = []string{"root", "test", "-n=john"}
		value, err := cmd.RunWithValue(ctx)
		t.AssertNil(err)
		t.Assert(value, `{"Content":"john"}`)
	})
}

func Test_Command_AddObject(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			command = gcmd.Command{
				Name: "start",
			}
		)
		err := command.AddObject(&TestCmdObject{})
		t.AssertNil(err)

		os.Args = []string{"start", "root", "test", "-n=john"}
		value, err := command.RunWithValue(ctx)
		t.AssertNil(err)
		t.Assert(value, `{"Content":"john"}`)
	})
}
