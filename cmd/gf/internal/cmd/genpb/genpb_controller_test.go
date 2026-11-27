// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests controller generation for unary and streaming gRPC services.

package genpb

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// TestControllerStreaming covers both protoc stream formats and ignores legacy stream interfaces.
func TestControllerStreaming(t *testing.T) {
	for _, tc := range []struct {
		name    string
		generic bool
		unary   bool
	}{
		{name: "legacy", unary: true},
		{name: "generic", generic: true, unary: true},
		{name: "legacy_stream_only"},
		{name: "generic_stream_only", generic: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					output = t.TempDir()
					api    = gfile.Join(output, "api", "chat", "v1")
					ctrl   = gfile.Join(output, "controller")
					source = `package v1
import (
    context "context"
    grpc "google.golang.org/grpc"
    emptypb "google.golang.org/protobuf/types/known/emptypb"
    v1 "example.com/demo/common/v1"
)
type ChatServer interface {
    Unary(context.Context, *Request) (*Reply, error)
    Watch(*Request, Chat_WatchServer) error
    Upload(Chat_UploadServer) error
    Chat(Chat_ChatServer) error
    External(*emptypb.Empty, Chat_ExternalServer) error
    Shared(*v1.Shared, Chat_SharedServer) error
    mustEmbedUnimplementedChatServer()
}
type Chat_SharedServer interface {
    Send(*Reply) error
    grpc.ServerStream
}
type Chat_WatchServer interface {
    Send(*Reply) error
    grpc.ServerStream
}
type Chat_UploadServer interface {
    SendAndClose(*Reply) error
    Recv() (*Request, error)
    grpc.ServerStream
}
type Chat_ChatServer interface {
    Send(*Reply) error
    Recv() (*Request, error)
    grpc.ServerStream
}
type Chat_ExternalServer interface {
    Send(*emptypb.Empty) error
    grpc.ServerStream
}
`
					expected = []string{
						"Unary(ctx context.Context, req *v1.Request) (res *v1.Reply, err error)",
						"Watch(*v1.Request, v1.Chat_WatchServer) error",
						"Upload(v1.Chat_UploadServer) error",
						"Chat(v1.Chat_ChatServer) error",
						"External(*emptypb.Empty, v1.Chat_ExternalServer) error",
						"Shared(*v11.Shared, v1.Chat_SharedServer) error",
					}
				)
				if tc.generic {
					source = gstr.ReplaceByMap(source, map[string]string{
						"Chat_WatchServer)":    "grpc.ServerStreamingServer[Reply])",
						"Chat_UploadServer)":   "grpc.ClientStreamingServer[Request, Reply])",
						"Chat_ChatServer)":     "grpc.BidiStreamingServer[Request, Reply])",
						"Chat_SharedServer)":   "grpc.ServerStreamingServer[Reply])",
						"Chat_ExternalServer)": "grpc.ServerStreamingServer[emptypb.Empty])",
					})
					// Generic protoc output uses aliases rather than the legacy stream interfaces.
					source, _, _ = strings.Cut(source, "type Chat_SharedServer interface")
					expected[1] = "Watch(*v1.Request, grpc.ServerStreamingServer[v1.Reply]) error"
					expected[2] = "Upload(grpc.ClientStreamingServer[v1.Request, v1.Reply]) error"
					expected[3] = "Chat(grpc.BidiStreamingServer[v1.Request, v1.Reply]) error"
					expected[4] = "External(*emptypb.Empty, grpc.ServerStreamingServer[emptypb.Empty]) error"
					expected[5] = "Shared(*v11.Shared, grpc.ServerStreamingServer[v1.Reply]) error"
				}
				if !tc.unary {
					source = gstr.Replace(source, "Unary(context.Context, *Request) (*Reply, error)", "")
					expected = expected[1:]
				}
				t.AssertNil(gfile.PutContents(gfile.Join(output, "go.mod"), "module example.com/demo\n\ngo 1.25.0\n"))
				t.AssertNil(gfile.PutContents(gfile.Join(api, "chat_grpc.pb.go"), source))
				controllers, err := (CGenPb{}).parseControllers(gfile.Join(api, "chat_grpc.pb.go"))
				t.AssertNil(err)
				t.Assert(len(controllers), 1)
				t.Assert(controllers[0].Name, "Chat")
				t.Assert(len(controllers[0].Methods), len(expected))
				for i, method := range controllers[0].Methods {
					t.Assert(method.Definition, expected[i])
				}
				in := generateControllerInput{OutputApiPath: gfile.Dir(gfile.Dir(api)), OutputCtrlPath: ctrl}
				t.AssertNil((CGenPb{}).generateController(context.Background(), in))
				filePath := gfile.Join(ctrl, "chat", "chat.go")
				content := gfile.GetContents(filePath)
				file, err := parser.ParseFile(token.NewFileSet(), filePath, content, 0)
				t.AssertNil(err)
				imports := make(map[string]bool)
				for _, spec := range file.Imports {
					importPath, unquoteErr := strconv.Unquote(spec.Path.Value)
					t.AssertNil(unquoteErr)
					imports[importPath] = true
				}
				t.Assert(imports["example.com/demo/api/chat/v1"], true)
				t.Assert(imports["example.com/demo/common/v1"], true)
				t.Assert(imports["google.golang.org/grpc"], tc.generic)
				t.Assert(imports["context"], tc.unary)
				t.Assert(imports["google.golang.org/protobuf/types/known/emptypb"], true)
				methods := make(map[string]int)
				for _, declaration := range file.Decls {
					method, ok := declaration.(*ast.FuncDecl)
					if !ok || method.Recv == nil {
						continue
					}
					statement, ok := method.Body.List[0].(*ast.ReturnStmt)
					t.Assert(ok, true)
					methods[method.Name.Name] = len(statement.Results)
				}
				expectedReturns := map[string]int{"Watch": 1, "Upload": 1, "Chat": 1, "External": 1, "Shared": 1}
				if tc.unary {
					expectedReturns["Unary"] = 2
				}
				t.Assert(methods, expectedReturns)
				// Re-running must not duplicate methods/imports or overwrite an existing implementation.
				custom := strings.Replace(content, "return gerror.NewCode(gcode.CodeNotImplemented)", "return nil", 1)
				t.AssertNil(gfile.PutContents(filePath, custom))
				t.AssertNil((CGenPb{}).generateController(context.Background(), in))
				t.Assert(gfile.GetContents(filePath), custom)
			})
		})
	}
}

// TestControllerStreamingProtoc compiles controllers generated from real generic and legacy protobuf output.
func TestControllerStreamingProtoc(t *testing.T) {
	for _, binary := range []string{"protoc", "protoc-gen-go", "protoc-gen-go-grpc"} {
		if _, err := exec.LookPath(binary); err != nil {
			t.Skipf("%s is required for the protobuf integration test: %v", binary, err)
		}
	}
	for _, generic := range []bool{false, true} {
		t.Run(strconv.FormatBool(generic), func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					output    = t.TempDir()
					api       = gfile.Join(output, "api", "chat", "v1")
					ctrl      = gfile.Join(output, "controller")
					protoPath = gfile.Join(output, "chat.proto")
					root      = gfile.RealPath("../../../../..")
				)
				t.AssertNE(root, "")
				module := fmt.Sprintf(`module example.com/demo

go 1.25.0

require (
    github.com/gogf/gf/v2 v2.10.5
    github.com/gogf/gf/contrib/rpc/grpcx/v2 v2.10.5
    google.golang.org/protobuf v1.36.12
)

replace github.com/gogf/gf/v2 => %s
replace github.com/gogf/gf/contrib/rpc/grpcx/v2 => %s
replace github.com/gogf/gf/contrib/registry/file/v2 => %s
`, strconv.Quote(filepath.ToSlash(root)),
					strconv.Quote(filepath.ToSlash(gfile.Join(root, "contrib/rpc/grpcx"))),
					strconv.Quote(filepath.ToSlash(gfile.Join(root, "contrib/registry/file"))))
				t.AssertNil(gfile.PutContents(gfile.Join(output, "go.mod"), module))
				t.AssertNil(gfile.PutContents(gfile.Join(output, "common", "v1", "shared.proto"), `syntax = "proto3";
package common.v1;
option go_package = "example.com/demo/common/v1;v1";
message Shared {}
`))
				t.AssertNil(gfile.PutContents(protoPath, `syntax = "proto3";
package chat.v1;
option go_package = "example.com/demo/api/chat/v1;v1";
import "google/protobuf/empty.proto";
import "common/v1/shared.proto";
message Request {}
message Reply {}
service Chat {
    rpc Unary(Request) returns (Reply);
    rpc Watch(Request) returns (stream Reply);
    rpc Upload(stream Request) returns (Reply);
    rpc Chat(stream Request) returns (stream Reply);
    rpc External(google.protobuf.Empty) returns (stream google.protobuf.Empty);
    rpc Shared(common.v1.Shared) returns (stream Reply);
}
`))
				t.AssertNil(gfile.Mkdir(api))
				common := gfile.Join(output, "common", "v1")
				command := exec.Command("protoc", "--proto_path="+output,
					"--go_out=paths=source_relative:"+output, gfile.Join(common, "shared.proto"))
				result, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("%v: %s", err, result)
				}
				command = exec.Command("protoc", "--proto_path="+output,
					"--go_out=paths=source_relative:"+api,
					"--go-grpc_out=paths=source_relative:"+api,
					"--go-grpc_opt=use_generic_streams_experimental="+strconv.FormatBool(generic), protoPath)
				result, err = command.CombinedOutput()
				if err != nil {
					t.Fatalf("%v: %s", err, result)
				}
				t.AssertNil((CGenPb{}).generateController(context.Background(), generateControllerInput{
					OutputApiPath: gfile.Join(output, "api"), OutputCtrlPath: ctrl,
				}))
				// Compile the actual generated API and controller against the local framework modules.
				command = exec.Command("go", "test", "-mod=mod", "./...")
				command.Dir = output
				command.Env = append(os.Environ(), "GOWORK=off")
				result, err = command.CombinedOutput()
				if err != nil {
					t.Fatalf("%v: %s", err, result)
				}
			})
		})
	}
}
