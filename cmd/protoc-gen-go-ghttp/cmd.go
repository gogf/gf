package main

import (
	"bytes"
	"fmt"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/http"
	"strings"
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
	goModPath                = getGoModImportName()
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
	gen.P("var _ = ", contextPackage.Ident("Background"), "()")
	gen.P("var _ = ", goframePackage.Ident("Meta"), "{}")
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
	// 获取路由和方法类型
	uri, apiMethod := getOptionMethodUri(rule)
	svcStructBuffer := bytes.NewBuffer(nil)
	err := unimplementedTpl.Execute(svcStructBuffer, map[string]interface{}{
		"svc_name":       string(svc.Desc.Name()),
		"method_name":    string(method.Desc.Name()),
		"in_name":        string(method.Input.Desc.Name()),
		"out_name":       string(method.Output.Desc.Name()),
		"method_comment": scanMethodComment(method),
		"http_pattern":   uri,
		"http_method":    apiMethod,
		"original_name":  method.Desc.Name(),
		"svr_name":       svc.Desc.FullName(),
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
			"fields":         processField(message, rule, needGenGMeta, g),
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

func processField(message *protogen.Message, rule *annotations.HttpRule, needGenGMeta bool, gen *protogen.GeneratedFile) []string {
	result := []string{}
	if needGenGMeta {
		uri, apiMethod := getOptionMethodUri(rule)
		result = append(result, fmt.Sprintf("g.Meta     `path:\"%s\" method:\"%s\"`", uri, apiMethod))
	}
	for _, item := range message.Fields {
		field := fmt.Sprintf("%s %s", item.GoName, processFieldType(item, gen))
		goComment := processFieldComment(item)
		if goComment != "" {
			field += " " + goComment
		}
		result = append(result, field)
	}
	return result
}

func processFieldType(field *protogen.Field, gen *protogen.GeneratedFile) string {
	if field.Desc.IsWeak() {
		return "struct{}"
	}
	goType := field.Desc.Kind().String()
	if field.Desc.Kind() == protoreflect.MessageKind {
		goType = field.Message.GoIdent.GoName
		if field.GoIdent.GoImportPath != field.Message.GoIdent.GoImportPath {
			goType = gen.QualifiedGoIdent(getFullImportPath(field.Message.GoIdent.GoImportPath).Ident(field.Message.GoIdent.GoName))
		}
		goType = "*" + goType
	} else if field.Desc.Kind() == protoreflect.EnumKind {
		goType = gen.QualifiedGoIdent(field.Enum.GoIdent)
	} else if field.Desc.HasPresence() && field.Desc.Kind() != protoreflect.MessageKind && field.Desc.Kind() != protoreflect.BytesKind {
		goType = "*" + goType
	}
	if field.Desc.IsList() {
		return "[]" + goType
	}
	if field.Desc.IsMap() {
		keyType := processFieldType(field.Message.Fields[0], gen)
		valType := processFieldType(field.Message.Fields[1], gen)
		return fmt.Sprintf("map[%v]%v", keyType, valType)
	}
	return goType
}

func getFullImportPath(path protogen.GoImportPath) protogen.GoImportPath {
	if fullImportPath == nil || !*fullImportPath || filePath == nil || *filePath == "*" {
		return path
	}
	pathStr := strings.Trim(path.String(), "\"")
	filePathStr := *filePath
	filePathStr = strings.TrimLeft(filePathStr, ".")
	filePathStr = strings.TrimRight(filePathStr, "/")
	return protogen.GoImportPath(goModPath + filePathStr + pathStr)
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
