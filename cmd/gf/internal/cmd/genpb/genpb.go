package genpb

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)

type (
	CGenPb      struct{}
	CGenPbInput struct {
		g.Meta     `name:"pb" config:"gfcli.gen.pb" brief:"parse proto files and generate protobuf go files"`
		Path       string `name:"path"       short:"p"  dc:"protobuf file folder path" d:"manifest/protobuf"`
		OutputApi  string `name:"outputApi"  short:"oa" dc:"output folder path storing generated go files of api" d:"api"`
		OutputCtrl string `name:"outputCtrl" short:"oc" dc:"output folder path storing generated go files of controller" d:"internal/controller"`
	}
	CGenPbOutput struct{}
)

func (c CGenPb) Pb(ctx context.Context, in CGenPbInput) (out *CGenPbOutput, err error) {
	// Necessary check.
	protoc := gproc.SearchBinary("protoc")
	if protoc == "" {
		mlog.Fatalf(`command "protoc" not found in your environment, please install protoc first to proceed this command`)
	}

	// protocol fold checks.
	protoPath := gfile.RealPath(in.Path)
	if protoPath == "" {
		mlog.Fatalf(`proto files folder "%s" does not exist`, in.Path)
	}
	// output path checks.
	outputApiPath := gfile.RealPath(in.OutputApi)
	if outputApiPath == "" {
		mlog.Fatalf(`output api folder "%s" does not exist`, in.OutputApi)
	}
	outputCtrlPath := gfile.RealPath(in.OutputCtrl)
	if outputCtrlPath == "" {
		mlog.Fatalf(`output controller folder "%s" does not exist`, in.OutputCtrl)
	}

	// folder scanning.
	files, err := gfile.ScanDirFile(protoPath, "*.proto", true)
	if err != nil {
		mlog.Fatal(err)
	}
	if len(files) == 0 {
		mlog.Fatalf(`no proto files found in folder "%s"`, in.Path)
	}

	if err = gfile.Chdir(protoPath); err != nil {
		mlog.Fatal(err)
	}
	for _, file := range files {
		var command = gproc.NewProcess(protoc, nil)
		command.Args = append(command.Args, "--proto_path="+gfile.Pwd())
		command.Args = append(command.Args, "--go_out=paths=source_relative:"+outputApiPath)
		command.Args = append(command.Args, "--go-grpc_out=paths=source_relative:"+outputApiPath)
		command.Args = append(command.Args, file)
		mlog.Print(command.String())
		if err = command.Run(ctx); err != nil {
			mlog.Fatal(err)
		}
	}
	// Generate struct tag according comment rules.
	err = c.generateStructTag(ctx, generateStructTagInput{OutputApiPath: outputApiPath})
	if err != nil {
		return
	}
	// Generate controllers according comment rules.
	err = c.generateController(ctx, generateControllerInput{
		OutputApiPath:  outputApiPath,
		OutputCtrlPath: outputCtrlPath,
	})
	if err != nil {
		return
	}
	mlog.Print("done!")
	return
}
