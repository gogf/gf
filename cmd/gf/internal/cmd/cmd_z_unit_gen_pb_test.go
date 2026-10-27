// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpb"
)

// assertControllerImports asserts that the given controller file can be parsed, and that every
// package qualifier it uses, like `grpcx` in `grpcx.GrpcServer`, is provided by its import
// declaration. It is a lightweight proxy for "the generated controller compiles".
func assertControllerImports(t *gtest.T, filePath string) {
	var (
		fset    = token.NewFileSet()
		content = gfile.GetContents(filePath)
	)
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	t.AssertNil(err)

	// The local names bound by the import declaration of the file.
	var importedNames = make(map[string]struct{})
	for _, spec := range file.Imports {
		var name = gstr.Trim(spec.Path.Value, `"`)
		if spec.Name != nil {
			name = spec.Name.Name
		} else {
			name = gfile.Basename(name)
		}
		importedNames[name] = struct{}{}
	}
	// Every qualifier must be provided by the import declaration. Identifiers bound to a
	// declaration inside the file, like the local variable or the function parameter, are
	// not package qualifiers.
	ast.Inspect(file, func(n ast.Node) bool {
		selector, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := selector.X.(*ast.Ident)
		if !ok || ident.Obj != nil {
			return true
		}
		_, ok = importedNames[ident.Name]
		if !ok {
			t.Errorf(`%s: missing import for qualifier %q`, filePath, ident.Name)
		}
		return true
	})
}

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

// TestGenPb_ControllerImports asserts that the generated controller declares every import its
// content uses. GoFmt only removes the unused imports and never adds them, so the generator has
// to declare them explicitly, otherwise the generated controller does not compile.
func TestGenPb_ControllerImports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The output folder is created inside the module, so that the import path of the api
		// package can be resolved, as `gf gen pb` is used in a real project.
		var (
			outputPath     = gfile.Join(gfile.Pwd(), "tmp_genpb_"+guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")
			protobufFolder = gtest.DataPath("issue", "3953")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
		)
		defer gfile.Remove(outputPath)
		t.AssertNil(gfile.Mkdir(outputApiPath))
		t.AssertNil(gfile.Mkdir(outputCtrlPath))

		_, err := genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		var (
			controllers, scanErr = gfile.ScanDirFile(outputCtrlPath, "*.go", true)
		)
		t.AssertNil(scanErr)
		t.Assert(len(controllers), 1)
		var content = gfile.GetContents(controllers[0])
		// The controller embeds the unimplemented server of the api version package, which is
		// imported with the version name given by the generator as its alias.
		t.Assert(gstr.Contains(content, `api "github.com/gogf/gf/cmd/gf/v2/`), true)
		t.Assert(gstr.Contains(content, "api.UnimplementedAccountServer"), true)
		t.Assert(gstr.Contains(content, "grpcx.GrpcServer"), true)
		// The method body uses gerror and gcode, and its signature uses context.
		t.Assert(gstr.Contains(content, "gerror.NewCode(gcode.CodeNotImplemented)"), true)
		t.Assert(gstr.Contains(content, "context.Context"), true)
		assertControllerImports(t, controllers[0])

		// Re-generating an existing controller must not duplicate the imports.
		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)
		content = gfile.GetContents(controllers[0])
		t.Assert(gstr.Count(content, `"context"`), 1)
		t.Assert(gstr.Count(content, `"github.com/gogf/gf/v2/net/grpcx"`), 1)
		assertControllerImports(t, controllers[0])
	})

	gtest.C(t, func(t *gtest.T) {
		// When the method is appended to an existing controller that does not import gerror and
		// gcode, the generator ensures them, otherwise the appended method does not compile.
		var (
			outputPath     = gfile.Join(gfile.Pwd(), "tmp_genpb_"+guid.S())
			outputApiPath  = filepath.Join(outputPath, "api")
			outputCtrlPath = filepath.Join(outputPath, "controller")
			protobufFolder = gtest.DataPath("issue", "3953")
			in             = genpb.CGenPbInput{
				Path:       protobufFolder,
				OutputApi:  outputApiPath,
				OutputCtrl: outputCtrlPath,
			}
		)
		defer gfile.Remove(outputPath)
		t.AssertNil(gfile.Mkdir(outputApiPath))
		t.AssertNil(gfile.Mkdir(outputCtrlPath))

		_, err := genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)
		var (
			controllers, scanErr = gfile.ScanDirFile(outputCtrlPath, "*.go", true)
		)
		t.AssertNil(scanErr)
		t.Assert(len(controllers), 1)
		var controllerPath = controllers[0]

		// Drop the imports used by the methods only, simulating a controller written by the
		// user or generated by an older version of gf.
		var content = gfile.GetContents(controllerPath)
		content = gstr.Replace(content, "\t\"context\"\n", "")
		content = gstr.Replace(content, "\t\"github.com/gogf/gf/v2/errors/gcode\"\n", "")
		content = gstr.Replace(content, "\t\"github.com/gogf/gf/v2/errors/gerror\"\n", "")
		// Remove the generated methods as well, so that they are appended again.
		content = gstr.SubStr(content, 0, gstr.Pos(content, "func (*Controller)"))
		content = gstr.TrimRight(content) + "\n"
		t.AssertNil(gfile.PutContents(controllerPath, content))

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)
		content = gfile.GetContents(controllerPath)
		t.Assert(gstr.Contains(content, `"context"`), true)
		t.Assert(gstr.Contains(content, `"github.com/gogf/gf/v2/errors/gcode"`), true)
		t.Assert(gstr.Contains(content, `"github.com/gogf/gf/v2/errors/gerror"`), true)
		t.Assert(gstr.Contains(content, "gerror.NewCode(gcode.CodeNotImplemented)"), true)
		assertControllerImports(t, controllerPath)
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
