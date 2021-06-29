package pack

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/allyes"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"strings"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
USAGE 
    gf pack SRC DST

ARGUMENT
    SRC  source path for packing, which can be multiple source paths.
    DST  destination file path for packed file. if extension of the filename is ".go" and "-n" option is given, 
         it enables packing SRC to go file, or else it packs SRC into a binary file.

OPTION
    -n, --name      package name for output go file, it's set as its directory name if no name passed
    -p, --prefix    prefix for each file packed into the resource file

EXAMPLES
    gf pack public data.bin
    gf pack public,template data.bin
    gf pack public,template packed/data.go
    gf pack public,template,config packed/data.go
    gf pack public,template,config packed/data.go -n=packed -p=/var/www/my-app
    gf pack /var/www/public packed/data.go -n=packed
`))
}

func Run() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"n,name":   true,
		"p,prefix": true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	srcPath := parser.GetArg(2)
	dstPath := parser.GetArg(3)
	if srcPath == "" {
		mlog.Fatal("SRC path cannot be empty")
	}
	if dstPath == "" {
		mlog.Fatal("DST path cannot be empty")
	}
	if gfile.Exists(dstPath) && gfile.IsDir(dstPath) {
		mlog.Fatalf("DST path '%s' cannot be a directory", dstPath)
	}
	if !gfile.IsEmpty(dstPath) && !allyes.Check() {
		s := gcmd.Scanf("path '%s' is not empty, files might be overwrote, continue? [y/n]: ", dstPath)
		if strings.EqualFold(s, "n") {
			return
		}
	}
	var (
		name   = parser.GetOpt("name")
		prefix = parser.GetOpt("prefix")
	)
	if name == "" && gfile.ExtName(dstPath) == "go" {
		name = gfile.Basename(gfile.Dir(dstPath))
	}
	if name != "" {
		if err := gres.PackToGoFile(srcPath, dstPath, name, prefix); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	} else {
		if err := gres.PackToFile(srcPath, dstPath, prefix); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	}
	mlog.Print("done!")
}
