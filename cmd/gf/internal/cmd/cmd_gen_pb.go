package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)

type (
	cGenPbInput struct {
		g.Meta `name:"pb" brief:"parse proto files and generate protobuf go files"`
	}
	cGenPbOutput struct{}
)

func (c cGen) Pb(ctx context.Context, in cGenPbInput) (out *cGenPbOutput, err error) {
	// Necessary check.
	if gproc.SearchBinary("protoc") == "" {
		mlog.Fatalf(`command "protoc" not found in your environment, please install protoc first to proceed this command`)
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
		goPathSrc   = gfile.RealPath(gfile.Join(genv.Get("GOPATH").String(), "src"))
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
	// pbFolder := "protobuf"
	// _, _ = gfile.ScanDirFileFunc(pbFolder, "*.go", true, func(path string) string {
	//	content := gfile.GetContents(path)
	//	content = gstr.ReplaceByArray(content, g.SliceStr{
	//		`gtime "gtime"`, `gtime "github.com/gogf/gf/v2/os/gtime"`,
	//	})
	//	_ = gfile.PutContents(path, content)
	//	utils.GoFmt(path)
	//	return path
	// })
	mlog.Print("done!")
	return
}
