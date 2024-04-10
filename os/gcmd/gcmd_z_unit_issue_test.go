// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

type Issue3390CommandCase1 struct {
	*gcmd.Command
}

type Issue3390TestCase1 struct {
	g.Meta `name:"index" ad:"test"`
}

type Issue3390Case1Input struct {
	g.Meta `name:"index"`
	A      string `short:"a" name:"aa"`
	Be     string `short:"b" name:"bb"`
}

type Issue3390Case1Output struct {
	Content string
}

func (c Issue3390TestCase1) Index(ctx context.Context, in Issue3390Case1Input) (out *Issue3390Case1Output, err error) {
	out = &Issue3390Case1Output{
		Content: gjson.MustEncodeString(in),
	}
	return
}

func Test_Issue3390_Case1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root, err := gcmd.NewFromObject(Issue3390TestCase1{})
		t.AssertNil(err)
		command := &Issue3390CommandCase1{root}
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{"main", "-a", "aaa", "-b", "bbb"},
		)
		t.AssertNil(err)
		t.Assert(value.(*Issue3390Case1Output).Content, `{"A":"aaa","Be":"bbb"}`)
	})
}

type Issue3390CommandCase2 struct {
	*gcmd.Command
}

type Issue3390TestCase2 struct {
	g.Meta `name:"index" ad:"test"`
}

type Issue3390Case2Input struct {
	g.Meta `name:"index"`
	A      string `short:"b" name:"bb"`
	Be     string `short:"a" name:"aa"`
}

type Issue3390Case2Output struct {
	Content string
}

func (c Issue3390TestCase2) Index(ctx context.Context, in Issue3390Case2Input) (out *Issue3390Case2Output, err error) {
	out = &Issue3390Case2Output{
		Content: gjson.MustEncodeString(in),
	}
	return
}
func Test_Issue3390_Case2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root, err := gcmd.NewFromObject(Issue3390TestCase2{})
		t.AssertNil(err)
		command := &Issue3390CommandCase2{root}
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{"main", "-a", "aaa", "-b", "bbb"},
		)
		t.AssertNil(err)
		t.Assert(value.(*Issue3390Case2Output).Content, `{"A":"bbb","Be":"aaa"}`)
	})
}

type Issue3390CommandCase3 struct {
	*gcmd.Command
}

type Issue3390TestCase3 struct {
	g.Meta `name:"index" ad:"test"`
}

type Issue3390Case3Input struct {
	g.Meta `name:"index"`
	A      string `short:"b"`
	Be     string `short:"a"`
}

type Issue3390Case3Output struct {
	Content string
}

func (c Issue3390TestCase3) Index(ctx context.Context, in Issue3390Case3Input) (out *Issue3390Case3Output, err error) {
	out = &Issue3390Case3Output{
		Content: gjson.MustEncodeString(in),
	}
	return
}
func Test_Issue3390_Case3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root, err := gcmd.NewFromObject(Issue3390TestCase3{})
		t.AssertNil(err)
		command := &Issue3390CommandCase3{root}
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{"main", "-a", "aaa", "-b", "bbb"},
		)
		t.AssertNil(err)
		t.Assert(value.(*Issue3390Case3Output).Content, `{"A":"bbb","Be":"aaa"}`)
	})
}

type Issue3390CommandCase4 struct {
	*gcmd.Command
}

type Issue3390TestCase4 struct {
	g.Meta `name:"index" ad:"test"`
}

type Issue3390Case4Input struct {
	g.Meta `name:"index"`
	A      string `short:"a"`
	Be     string `short:"b"`
}

type Issue3390Case4Output struct {
	Content string
}

func (c Issue3390TestCase4) Index(ctx context.Context, in Issue3390Case4Input) (out *Issue3390Case4Output, err error) {
	out = &Issue3390Case4Output{
		Content: gjson.MustEncodeString(in),
	}
	return
}

func Test_Issue3390_Case4(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root, err := gcmd.NewFromObject(Issue3390TestCase4{})
		t.AssertNil(err)
		command := &Issue3390CommandCase4{root}
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{"main", "-a", "aaa", "-b", "bbb"},
		)
		t.AssertNil(err)
		t.Assert(value.(*Issue3390Case4Output).Content, `{"A":"aaa","Be":"bbb"}`)
	})
}

type Issue3417Test struct {
	g.Meta `name:"root"`
}

type Issue3417BuildInput struct {
	g.Meta        `name:"build" config:"gfcli.build"`
	File          string `name:"FILE" arg:"true"    brief:"building file path"`
	Name          string `short:"n"  name:"name"    brief:"output binary name"`
	Version       string `short:"v"  name:"version" brief:"output binary version"`
	Arch          string `short:"a"  name:"arch"    brief:"output binary architecture, multiple arch separated with ','"`
	System        string `short:"s"  name:"system"  brief:"output binary system, multiple os separated with ','"`
	Output        string `short:"o"  name:"output"  brief:"output binary path, used when building single binary file"`
	Path          string `short:"p"  name:"path"    brief:"output binary directory path, default is '.'" d:"."`
	Extra         string `short:"e"  name:"extra"   brief:"extra custom \"go build\" options"`
	Mod           string `short:"m"  name:"mod"     brief:"like \"-mod\" option of \"go build\", use \"-m none\" to disable go module"`
	Cgo           bool   `short:"c"  name:"cgo"     brief:"enable or disable cgo feature, it's disabled in default" orphan:"true"`
	VarMap        g.Map  `short:"r"  name:"varMap"  brief:"custom built embedded variable into binary"`
	PackSrc       string `short:"ps" name:"packSrc" brief:"pack one or more folders into one go file before building"`
	PackDst       string `short:"pd" name:"packDst" brief:"temporary go file path for pack, this go file will be automatically removed after built" d:"internal/packed/build_pack_data.go"`
	ExitWhenError bool   `short:"ew" name:"exitWhenError" brief:"exit building when any error occurs, specially for multiple arch and system buildings. default is false" orphan:"true"`
	DumpENV       bool   `short:"de" name:"dumpEnv" brief:"dump current go build environment before building binary" orphan:"true"`
}

type Issue3417BuildOutput struct {
	Content string
}

func (c *Issue3417Test) Build(ctx context.Context, in Issue3417BuildInput) (out *Issue3417BuildOutput, err error) {
	out = &Issue3417BuildOutput{
		Content: gjson.MustEncodeString(in),
	}
	return
}

func Test_Issue3417(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		command, err := gcmd.NewFromObject(Issue3417Test{})
		t.AssertNil(err)
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{
				"gf", "build",
				"-mod", "vendor",
				"-v", "0.0.19",
				"-n", "detect_hardware_os",
				"-a", "amd64,arm64",
				"-s", "linux",
				"-p", "./bin",
				"-e", "-trimpath -ldflags",
				"cmd/v3/main.go",
			},
		)
		t.AssertNil(err)
		t.Assert(
			value.(*Issue3417BuildOutput).Content,
			`{"File":"cmd/v3/main.go","Name":"detect_hardware_os","Version":"0.0.19","Arch":"amd64,arm64","System":"linux","Output":"","Path":"./bin","Extra":"-trimpath -ldflags","Mod":"vendor","Cgo":false,"VarMap":null,"PackSrc":"","PackDst":"internal/packed/build_pack_data.go","ExitWhenError":false,"DumpENV":false}`,
		)
	})
}
