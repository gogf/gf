package update

import (
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"runtime"
)

func Run() {
	goBinPath := gproc.SearchBinary("go")
	if goBinPath == "" {
		mlog.Fatal(`"go" command not found, install it first to step further`)
	}
	var err error
	if gstr.CompareVersionGo(runtime.Version(), "go1.16.0") >= 0 {
		err = gproc.ShellRun(`go install github.com/gogf/gf/tool/gf@latest`)
	} else {
		err = gproc.ShellRun(`go install github.com/gogf/gf/tool/gf`)
	}
	if err != nil {
		mlog.Fatalf(`gf binary installation failed: %+v`, err)
	}
	mlog.Print("gf binary is now updated to the latest version")
}
