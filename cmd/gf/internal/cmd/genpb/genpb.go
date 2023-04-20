// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genpb

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gtag"
)

type (
	CGenPb      struct{}
	CGenPbInput struct {
		g.Meta     `name:"pb" config:"{CGenPbConfig}" brief:"{CGenPbBrief}" eg:"{CGenPbEg}"`
		Path       string `name:"path" short:"p"  dc:"protobuf file folder path" d:"manifest/protobuf"`
		OutputApi  string `name:"api"  short:"a"  dc:"output folder path storing generated go files of api" d:"api"`
		OutputCtrl string `name:"ctrl" short:"c"  dc:"output folder path storing generated go files of controller" d:"internal/controller"`
	}
	CGenPbOutput struct{}
)

const (
	CGenPbConfig = `gfcli.gen.pb`
	CGenPbBrief  = `parse proto files and generate protobuf go files`
	CGenPbEg     = `
gf gen pb
gf gen pb -p . -a . -p .
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenPbEg`:     CGenPbEg,
		`CGenPbBrief`:  CGenPbBrief,
		`CGenPbConfig`: CGenPbConfig,
	})
}

func (c CGenPb) Pb(ctx context.Context, in CGenPbInput) (out *CGenPbOutput, err error) {
	// Necessary check.
	protoc := gproc.SearchBinary("protoc")
	if protoc == "" {
		mlog.Fatalf(`command "protoc" not found in your environment, please install protoc first: https://grpc.io/docs/languages/go/quickstart/`)
	}

	// protocol fold checks.
	var (
		protoPath    = gfile.RealPath(in.Path)
		isParsingPWD bool
	)
	if protoPath == "" {
		// Use current working directory as protoPath if there are proto files under.
		currentPath := gfile.Pwd()
		currentFiles, _ := gfile.ScanDirFile(currentPath, "*.proto")
		if len(currentFiles) > 0 {
			protoPath = currentPath
			isParsingPWD = true
		} else {
			mlog.Fatalf(`proto files folder "%s" does not exist`, in.Path)
		}
	}
	// output path checks.
	outputApiPath := gfile.RealPath(in.OutputApi)
	if outputApiPath == "" {
		if isParsingPWD {
			outputApiPath = protoPath
		} else {
			mlog.Fatalf(`output api folder "%s" does not exist`, in.OutputApi)
		}
	}
	outputCtrlPath := gfile.RealPath(in.OutputCtrl)
	if outputCtrlPath == "" {
		if isParsingPWD {
			outputCtrlPath = ""
		} else {
			mlog.Fatalf(`output controller folder "%s" does not exist`, in.OutputCtrl)
		}
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
	if outputCtrlPath != "" {
		err = c.generateController(ctx, generateControllerInput{
			OutputApiPath:  outputApiPath,
			OutputCtrlPath: outputCtrlPath,
		})
		if err != nil {
			return
		}
	}
	mlog.Print("done!")
	return
}
