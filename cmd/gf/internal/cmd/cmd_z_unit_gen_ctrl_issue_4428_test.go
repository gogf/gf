package cmd

import (
	"path/filepath"
	"testing"
	"context"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genctrl"
)

func Test_Gen_Ctrl_Issue4428(t *testing.T) {
    var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		var (
			rootPath  = gfile.Temp(guid.S())
			apiFolder = filepath.Join(rootPath, "api")
			ctrlPath  = filepath.Join(rootPath, "internal", "controller")
			in        = genctrl.CGenCtrlInput{
				SrcFolder: apiFolder,
				DstFolder: ctrlPath,
				Merge:     true,
			}
		)
		err := gfile.Mkdir(rootPath)
		t.AssertNil(err)
		defer gfile.Remove(rootPath)

		// Create go.mod
		err = gfile.PutContents(filepath.Join(rootPath, "go.mod"), "module test\n\ngo 1.20\n")
		t.AssertNil(err)

		// 1. Create initial V2 API
		var apiV2Content = `
package v2

import "github.com/gogf/gf/v2/frame/g"

type HelloReq struct {
	g.Meta ` + "`path:\"/hello\" method:\"get\"`" + `
}

type HelloRes struct {}
`
		err = gfile.PutContents(filepath.Join(apiFolder, "hello", "v2", "hello.go"), apiV2Content)
		t.AssertNil(err)

		// 2. Generate controller
		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		// Check generated file
		ctrlFile := filepath.Join(ctrlPath, "hello", "hello_v2_hello.go")
		t.Assert(gfile.Exists(ctrlFile), true)
		content := gfile.GetContents(ctrlFile)
		t.Assert(gstr.Count(content, "func (c *ControllerV2) Hello("), 1)

		// 3. Add conflicting import to controller and save it
		// Simulating user adding "excelize" import
		newContent := gstr.Replace(content, `import (`, `import (
	excelize "github.com/xuri/excelize/v2"`, 1)
		err = gfile.PutContents(ctrlFile, newContent)
		t.AssertNil(err)

		// 4. Add new API to V2
		var apiV2NewContent = `
type WorldReq struct {
	g.Meta ` + "`path:\"/world\" method:\"get\"`" + `
}

type WorldRes struct {}
`
		err = gfile.PutContentsAppend(filepath.Join(apiFolder, "hello", "v2", "hello.go"), apiV2NewContent)
		t.AssertNil(err)

		// 5. Generate controller again
		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		// 6. Check for duplication
		content = gfile.GetContents(ctrlFile)
		// Hello should still appear exactly once
		t.Assert(gstr.Count(content, "func (c *ControllerV2) Hello("), 1)
		// World should appear exactly once
		t.Assert(gstr.Count(content, "func (c *ControllerV2) World("), 1)
	})
}
