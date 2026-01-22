// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpb"
)

func TestGenPbIssue3882(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			outputPath     = gfile.Temp(guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")

			protobufFolder = gtest.DataPath("issue", "3882")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
			err error
		)
		err = gfile.Mkdir(outputApiPath)
		t.AssertNil(err)
		err = gfile.Mkdir(outputCtrlPath)
		t.AssertNil(err)
		defer gfile.Remove(outputPath)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		var (
			genContent = gfile.GetContents(filepath.Join(outputApiPath, "issue3882.pb.go"))
			exceptText = `dc:"Some comment on field with 'one' 'two' 'three' in the comment."`
		)
		t.Assert(gstr.Contains(genContent, exceptText), true)
	})
}

// This issue only occurs when executing multiple times
// and the subsequent OutputApi is the parent directory of the previous execution
func TestGenPbIssue3953(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			outputPath     = gfile.Temp("f" + guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")

			protobufFolder = gtest.DataPath("issue", "3953")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
			err error
		)
		err = gfile.Mkdir(outputApiPath)
		t.AssertNil(err)
		err = gfile.Mkdir(outputCtrlPath)
		t.AssertNil(err)
		defer gfile.Remove(outputPath)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		// do twice,and set outputApi to outputPath
		in.OutputApi = outputPath
		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		var (
			genContent = gfile.GetContents(filepath.Join(outputApiPath, "issue3953.pb.go"))
			// The old version would have appeared `v:"required" v:"required"`
			// but the new version of the code will appear `v:"required"` only once
			notExceptText = `v:"required" v:"required"`
		)
		t.Assert(gstr.Contains(genContent, notExceptText), false)
	})
}

func TestGenPb_MultipleTags(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			outputPath     = gfile.Temp(guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")

			protobufFolder = gtest.DataPath("genpb")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
			err error
		)
		err = gfile.Mkdir(outputApiPath)
		t.AssertNil(err)
		err = gfile.Mkdir(outputCtrlPath)
		t.AssertNil(err)
		defer gfile.Remove(outputPath)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		// Test multiple_tags.proto output
		genContent := gfile.GetContents(filepath.Join(outputApiPath, "multiple_tags.pb.go"))
		// Id field should have combined validation tags: v:"required#Id > 0"
		t.Assert(gstr.Contains(genContent, `v:"required#Id > 0"`), true)
		// Name field should have dc tag from plain comment
		t.Assert(gstr.Contains(genContent, `dc:"User name for login"`), true)
		// Email field should have combined validation and dc tag
		t.Assert(gstr.Contains(genContent, `v:"requiredemail"`), true)
		t.Assert(gstr.Contains(genContent, `dc:"User email address"`), true)
	})
}

func TestGenPb_NestedMessage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			outputPath     = gfile.Temp(guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")

			protobufFolder = gtest.DataPath("genpb")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
			err error
		)
		err = gfile.Mkdir(outputApiPath)
		t.AssertNil(err)
		err = gfile.Mkdir(outputCtrlPath)
		t.AssertNil(err)
		defer gfile.Remove(outputPath)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		// Test nested_message.proto output
		genContent := gfile.GetContents(filepath.Join(outputApiPath, "nested_message.pb.go"))
		// Order.OrderId should have v:"required"
		t.Assert(gstr.Contains(genContent, `v:"required"`), true)
		// Order.Detail should have dc:"Order details"
		t.Assert(gstr.Contains(genContent, `dc:"Order details"`), true)
		// OrderDetail.Quantity should have v:"min:1"
		t.Assert(gstr.Contains(genContent, `v:"min:1"`), true)
		// OrderDetail.Price should have v:"min:0.01"
		t.Assert(gstr.Contains(genContent, `v:"min:0.01"`), true)
	})
}
