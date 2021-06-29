package get

import (
	"fmt"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"os"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
USAGE    
    gf get PACKAGE

ARGUMENT 
    PACKAGE  remote golang package path, eg: github.com/gogf/gf

EXAMPLES
    gf get github.com/gogf/gf
    gf get github.com/gogf/gf@latest
    gf get github.com/gogf/gf@master
    gf get golang.org/x/sys
`))
}

func Run() {
	if len(os.Args) > 2 {
		gproc.ShellRun(fmt.Sprintf(`go get -u %s`, gstr.Join(os.Args[2:], " ")))
	} else {
		mlog.Fatal("please input the package path for get")
	}
}
