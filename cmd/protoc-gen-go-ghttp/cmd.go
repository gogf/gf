package main

import (
	"bytes"
	"fmt"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/http"
	"text/template"
)

var (
	// 现在的Import是先直接写死，到时候需要根据这些来进行写入
	contextPackage           = protogen.GoImportPath("context")
	goframePackage           = protogen.GoImportPath("github.com/gogf/gf/v2/frame/g")
	svcStructTpl, _          = template.New("templateSvcStruct").Parse(templateSvcStruct)
	unimplementedTpl, _      = template.New("unimplemented").Parse(templateRouterFunc)
	templateImplStructTpl, _ = template.New("templateImplStruct").Parse(templateImplStruct)
	versionCommentTpl, _     = template.New("versionComment").Parse(versionComment)
	methodMessageTpl, _      = template.New("methodMessageStruct").Parse(methodMessageStruct)
)

func process(genFile *protogen.Plugin, file *protogen.File) {
	gen := genFile.NewGeneratedFile(file.GeneratedFilenamePrefix+".ghttp.go", file.GoImportPath)
	processCopyrightAndVersion(gen, file, genFile)
}

func processCopyrightAndVersion(gen *protogen.GeneratedFile, file *protogen.File, genFile *protogen.Plugin) {
	versionBuffer := bytes.NewBuffer(nil)
	err := versionCommentTpl.Execute(versionBuffer, map[string]string{
		"protoc_version": getProtocVersion(genFile),
		"ghttp_version":  httpGenVersion,
		"source_path":    file.Desc.Path(),
	})
	if err != nil {
		info("gf-gen-go-http: Execute template error: %s\n", err.Error())
		panic(err.Error())
	}
	gen.P(versionBuffer.String())
	gen.P()
	gen.P("package ", file.GoPackageName)
	gen.P()
	gen.P(templateImport)
	gen.P()

	processContent(gen, file)
}

func processContent(gen *protogen.GeneratedFile, file *protogen.File) {
	for _, svc := range file.Services {
		processSvcStruct(gen, svc)
		serviceHttpInterfaces := []map[string]interface{}{}
		for _, method := range svc.Methods {
			if method == nil || method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
				continue
			}
			methodItem := processMethod(gen, method, svc)
			if len(methodItem) != 0 {
				serviceHttpInterfaces = append(serviceHttpInterfaces, methodItem)
			}
		}
		processSvcInterface(gen, serviceHttpInterfaces, string(svc.Desc.Name()))
	}
}

func processSvcStruct(gen *protogen.GeneratedFile, svc *protogen.Service) {
	svcStructBuffer := bytes.NewBuffer(nil)
	err := svcStructTpl.Execute(svcStructBuffer, map[string]string{"svc_name": string(svc.Desc.Name())})
	if err != nil {
		info("gf-gen-go-http: Execute template error: %s\n", err.Error())
		panic(err)
	}
	gen.P(svcStructBuffer.String())
	gen.P()
}

func processSvcInterface(gen *protogen.GeneratedFile, data []map[string]interface{}, svcName string) {
	svcStructBuffer := bytes.NewBuffer(nil)
	err := templateImplStructTpl.Execute(svcStructBuffer, map[string]interface{}{
		"svc_name": svcName,
		"svc_list": data,
	})
	if err != nil {
		info("gf-gen-go-http: Execute template error: %s\n", err.Error())
		panic(err.Error())
	}
	gen.P(svcStructBuffer.String())
	gen.P()
}

func processMethod(g *protogen.GeneratedFile, method *protogen.Method, svc *protogen.Service) map[string]interface{} {
	rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule == nil || !ok {
		return nil
	}
	// 生成 input / output
	processMessage(g, method, rule)
	// 生成路由注册方法
	svcStructBuffer := bytes.NewBuffer(nil)
	err := unimplementedTpl.Execute(svcStructBuffer, map[string]interface{}{
		"svc_name":       string(svc.Desc.Name()),
		"method_name":    string(method.Desc.Name()),
		"in_name":        string(method.Input.Desc.Name()),
		"out_name":       string(method.Output.Desc.Name()),
		"method_comment": scanMethodComment(method),
	})
	if err != nil {
		info("gf-gen-go-http: Execute template error: %s\n", err.Error())
		panic(err.Error())
	}
	g.P(svcStructBuffer.String())
	g.P()
	return map[string]interface{}{
		"method_name": string(method.Desc.Name()),
		"in_name":     string(method.Input.Desc.Name()),
		"out_name":    string(method.Output.Desc.Name()),
	}
}

func processMessage(g *protogen.GeneratedFile, method *protogen.Method, rule *annotations.HttpRule) {
	processMessageFunc := func(message *protogen.Message, needGenGMeta bool) {
		methodMessageTplBuffer := bytes.NewBuffer(nil)
		err := methodMessageTpl.Execute(methodMessageTplBuffer, map[string]interface{}{
			"method_comment": scanMessageComment(message),
			"method_name":    string(method.Desc.Name()),
			"message_name":   string(message.Desc.Name()),
			"fields":         processField(message, rule, needGenGMeta),
		})
		if err != nil {
			info("gf-gen-go-http: Execute template error: %s\n", err.Error())
			panic(err.Error())
		}
		g.P(methodMessageTplBuffer.String())
		g.P()
	}
	processMessageFunc(method.Input, true)
	processMessageFunc(method.Output, false)
}

func processField(message *protogen.Message, rule *annotations.HttpRule, needGenGMeta bool) []string {
	result := []string{}
	if needGenGMeta {
		apiMethod, uri := getOptionMethodUri(rule)
		result = append(result, fmt.Sprintf("g.Meta     `path:\"%s\" method:\"%s\"`", uri, apiMethod))
	}
	for _, item := range message.Fields {
		goType := item.Desc.Kind().String()
		if item.Desc.Kind() == protoreflect.MessageKind {
			goType = "*" + string(item.Message.Desc.Name())
			if item.Desc.IsList() {
				goType = "[]" + goType
			}
			if item.Desc.IsMap() {
				goType = fmt.Sprintf("map[%s]*%s", item.Desc.MapKey().Kind().String(), string(item.Desc.MapValue().Message().Name()))
			}
		}
		field := fmt.Sprintf("%s %s", item.GoName, goType)
		goComment := processFieldComment(item)
		if goComment != "" {
			field += " " + goComment
		}
		result = append(result, field)
	}
	return result
}

func getOptionMethodUri(rule *annotations.HttpRule) (string, string) {
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		return pattern.Get, http.MethodGet
	case *annotations.HttpRule_Put:
		return pattern.Put, http.MethodPut
	case *annotations.HttpRule_Post:
		return pattern.Post, http.MethodPost
	case *annotations.HttpRule_Delete:
		return pattern.Delete, http.MethodDelete
	case *annotations.HttpRule_Patch:
		return pattern.Patch, http.MethodPatch
	case *annotations.HttpRule_Custom:
		return pattern.Custom.Path, pattern.Custom.Kind
	}
	return "(?)", "(?)"
}

func getProtocVersion(genFile *protogen.Plugin) string {
	ver := genFile.Request.GetCompilerVersion()
	return fmt.Sprintf("v%d.%d.%d", ver.GetMajor(), ver.GetMinor(), ver.GetPatch())
}
