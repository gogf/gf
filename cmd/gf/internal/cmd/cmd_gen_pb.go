package cmd

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)

type (
	cGenPb      struct{}
	cGenPbInput struct {
		g.Meta `name:"pb" brief:"parse proto files and generate protobuf go files"`
		Path   string `name:"path"   short:"p" dc:"protobuf file folder path" d:"manifest/protobuf"`
		Output string `name:"output" short:"o" dc:"output folder path storing generated go files" d:"api"`
	}
	cGenPbOutput struct{}
)

func (c cGenPb) Pb(ctx context.Context, in cGenPbInput) (out *cGenPbOutput, err error) {
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
	outputPath := gfile.RealPath(in.Output)
	if outputPath == "" {
		mlog.Fatalf(`output folder "%s" does not exist`, in.Output)
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
		command.Args = append(command.Args, "--go_out=paths=source_relative:"+outputPath)
		command.Args = append(command.Args, "--go-grpc_out=paths=source_relative:"+outputPath)
		command.Args = append(command.Args, file)
		mlog.Print(command.String())
		if err = command.Run(ctx); err != nil {
			mlog.Fatal(err)
		}
	}
	mlog.Print("done!")
	return
}
