package cmd

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	Fix = cFix{}
)

type cFix struct {
	g.Meta `name:"fix" brief:"auto fixing codes after upgrading to new GoFrame version" usage:"gf fix" `
}

type cFixInput struct {
	g.Meta `name:"fix"`
}

type cFixOutput struct{}

func (c cFix) Index(ctx context.Context, in cFixInput) (out *cFixOutput, err error) {
	mlog.Print(`start auto fixing...`)
	defer mlog.Print(`done!`)
	err = c.doFix(ctx)
	return
}

func (c cFix) doFix(ctx context.Context) (err error) {
	err = c.doFixV23(ctx)
	return
}

// doFixV23 fixes code when upgrading to GoFrame v2.3.
func (c cFix) doFixV23(ctx context.Context) error {
	replaceFunc := func(path, content string) string {
		content = gstr.Replace(content, "*gdb.TX", "gdb.TX")
		return content
	}
	return gfile.ReplaceDirFunc(replaceFunc, ".", "*.go", true)
}
