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
		if err != nil {
			panic(err)
		}
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
		if err != nil {
			panic(err)
		}
		command := &Issue3390CommandCase2{root}
		value, err := command.RunWithSpecificArgs(
			gctx.New(),
			[]string{"main", "-a", "aaa", "-b", "bbb"},
		)
		t.AssertNil(err)
		t.Assert(value.(*Issue3390Case2Output).Content, `{"A":"bbb","Be":"aaa"}`)
	})
}
