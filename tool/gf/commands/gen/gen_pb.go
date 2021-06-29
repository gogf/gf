package gen

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
)

func HelpPb() {
	mlog.Print(gstr.TrimLeft(`
USAGE 
    gf gen pb 

`))
}

// doGenPb parses current `proto` files in folder `protocol` and generates `pb` files to `protobuf`.
func doGenPb() {
	// protoc search.
	protocBinPath := gproc.SearchBinary("protoc")
	if protocBinPath == "" {
		mlog.Fatal(`"protoc" command not found, install it first to proceed proto files parsing`)
	}
	// protocol fold checks.
	protoFolder := "protocol"
	if !gfile.Exists(protoFolder) {
		mlog.Fatalf(`proto files folder "%s" does not exist`, protoFolder)
	}
	// folder scanning.
	files, err := gfile.ScanDirFile(protoFolder, "*.proto", true)
	if err != nil {
		mlog.Fatal(err)
	}
	if len(files) == 0 {
		mlog.Fatalf(`no proto files found in folder "%s"`, protoFolder)
	}
	dirSet := gset.NewStrSet()
	for _, file := range files {
		dirSet.Add(gfile.Dir(file))
	}
	var (
		servicePath = gfile.RealPath(".")
		goPathSrc   = gfile.RealPath(gfile.Join(genv.Get("GOPATH"), "src"))
	)
	dirSet.Iterator(func(protoDirPath string) bool {
		parsingCommand := fmt.Sprintf(
			"protoc --gofast_out=plugins=grpc:. %s/*.proto -I%s",
			protoDirPath,
			servicePath,
		)
		if goPathSrc != "" {
			parsingCommand += " -I" + goPathSrc
		}
		mlog.Print(parsingCommand)
		if output, err := gproc.ShellExec(parsingCommand); err != nil {
			mlog.Print(output)
			mlog.Fatal(err)
		}
		return true
	})
	// Custom replacement.
	//pbFolder := "protobuf"
	//_, _ = gfile.ScanDirFileFunc(pbFolder, "*.go", true, func(path string) string {
	//	content := gfile.GetContents(path)
	//	content = gstr.ReplaceByArray(content, g.SliceStr{
	//		`gtime "gtime"`, `gtime "github.com/gogf/gf/os/gtime"`,
	//	})
	//	_ = gfile.PutContents(path, content)
	//	utils.GoFmt(path)
	//	return path
	//})
	mlog.Print("done!")
}
