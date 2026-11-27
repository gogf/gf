// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file generates controllers from protobuf gRPC server interfaces.

package genpb

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

type generateControllerInput struct {
	OutputApiPath  string
	OutputCtrlPath string
}

type generateCtrl struct {
	Name    string
	Package string
	Version string
	// ApiImportPath is the import path of the api package of the controller, which is
	// resolved from the proto generated go files. It is empty when the output folder is
	// not inside a go module.
	ApiImportPath string
	// Imports contains packages referenced by streaming method signatures.
	Imports []utils.ImportSpec
	Methods []generateCtrlMethod
}

type generateCtrlMethod struct {
	Name       string
	Definition string
	// Streaming indicates a method that returns only an error.
	Streaming bool
}

const (
	// controllerContextImportPath is the import path of the `context` package used by the
	// generated controller methods.
	controllerContextImportPath = "context"
	// controllerGrpcxImportPath is the import path of the grpcx package used by the generated
	// controller registration function. The grpcx package lives in the contrib module, and it
	// is not part of the core framework module.
	controllerGrpcxImportPath = "github.com/gogf/gf/contrib/rpc/grpcx/v2"
	// controllerGcodeImportPath is the import path of the gcode package used by the generated
	// controller methods.
	controllerGcodeImportPath = "github.com/gogf/gf/v2/errors/gcode"
	// controllerGerrorImportPath is the import path of the gerror package used by the
	// generated controller methods.
	controllerGerrorImportPath = "github.com/gogf/gf/v2/errors/gerror"

	controllerTemplate = `
package {Package}

type Controller struct {
	{Version}.Unimplemented{Name}Server
}

func Register(s *grpcx.GrpcServer) {
	{Version}.Register{Name}Server(s.Server, &Controller{})
}
`
	controllerMethodTemplate = `
func (*Controller) {Definition} {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
`
	// controllerStreamMethodTemplate returns the single error required by streaming RPCs.
	controllerStreamMethodTemplate = `
func (*Controller) {Definition} {
	return gerror.NewCode(gcode.CodeNotImplemented)
}
`
)

func (c CGenPb) generateController(ctx context.Context, in generateControllerInput) (err error) {
	files, err := gfile.ScanDirFile(in.OutputApiPath, "*_grpc.pb.go", true)
	if err != nil {
		return err
	}
	var controllers []generateCtrl
	for _, file := range files {
		fileControllers, err := c.parseControllers(file)
		if err != nil {
			return err
		}
		controllers = append(controllers, fileControllers...)
	}
	if len(controllers) == 0 {
		return nil
	}
	// Generate controller files.
	err = c.doGenerateControllers(in, controllers)
	return
}

func (c CGenPb) parseControllers(filePath string) ([]generateCtrl, error) {
	var (
		controllers []generateCtrl
		content     = gfile.GetContents(filePath)
		parseErr    error
	)
	file, err := parser.ParseFile(token.NewFileSet(), filePath, content, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	// Reuse the generated protobuf file's aliases for external message and gRPC types.
	var (
		imports = make(map[string]utils.ImportSpec)
		aliases = make(map[string]string)
		// The controller also imports its API and framework packages, which may collide with protobuf aliases.
		usedNames = map[string]string{
			gfile.Basename(gfile.Dir(filePath)): utils.GetImportPath(gfile.Dir(filePath)),
			"context":                           controllerContextImportPath,
			"grpcx":                             controllerGrpcxImportPath,
			"gcode":                             controllerGcodeImportPath,
			"gerror":                            controllerGerrorImportPath,
		}
	)
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		name := path.Base(importPath)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		alias := name
		for suffix := 1; ; suffix++ {
			if existing, ok := usedNames[alias]; !ok || existing == importPath {
				break
			}
			alias = name + strconv.Itoa(suffix)
		}
		usedNames[alias] = importPath
		aliases[name] = alias
		imports[name] = utils.ImportSpec{Name: alias, Path: importPath}
	}
	_, err = gregex.ReplaceStringFuncMatch(
		`type (\w+)Server interface {([\s\S]+?)}`,
		content,
		func(match []string) string {
			// Legacy per-RPC stream interfaces are not services and have no registration function.
			if parseErr != nil || gregex.IsMatchString(`(?m)^\s*grpc\.ServerStream\s*$`, match[2]) {
				return match[0]
			}
			ctrl := generateCtrl{
				Name:          match[1],
				Package:       strings.ReplaceAll(gfile.Basename(gfile.Dir(gfile.Dir(filePath))), "-", "_"),
				Version:       gfile.Basename(gfile.Dir(filePath)),
				ApiImportPath: utils.GetImportPath(gfile.Dir(filePath)),
				Methods:       make([]generateCtrlMethod, 0),
			}
			lines := gstr.Split(match[2], "\n")
			for _, line := range lines {
				line = gstr.Trim(line)
				if line == "" || !gstr.IsLetterUpper(line[0]) {
					continue
				}
				// Comment.
				if gregex.IsMatchString(`^//.+`, line) {
					continue
				}
				line, _ = gregex.ReplaceStringFuncMatch(
					`^(\w+)\(context\.Context, \*(\w+)\) \(\*(\w+), error\)$`,
					line,
					func(match []string) string {
						return fmt.Sprintf(
							`%s(ctx context.Context, req *%s.%s) (res *%s.%s, err error)`,
							match[1], ctrl.Version, match[2], ctrl.Version, match[3],
						)
					},
				)
				streaming := strings.HasSuffix(line, ") error")
				if streaming {
					var packages []string
					line, packages, parseErr = controllerStreamingDefinition(line, ctrl.Version, aliases)
					if parseErr != nil {
						return match[0]
					}
					for _, name := range packages {
						if spec, ok := imports[name]; ok {
							ctrl.Imports = append(ctrl.Imports, spec)
						}
					}
				}
				ctrl.Methods = append(ctrl.Methods, generateCtrlMethod{
					Name:       gstr.Split(line, "(")[0],
					Definition: line,
					Streaming:  streaming,
				})
			}
			if len(ctrl.Methods) > 0 {
				controllers = append(controllers, ctrl)
			}
			return match[0]
		},
	)
	if parseErr != nil {
		return nil, parseErr
	}
	return controllers, err
}

func (c CGenPb) doGenerateControllers(in generateControllerInput, controllers []generateCtrl) (err error) {
	for _, controller := range controllers {
		err = c.doGenerateController(in, controller)
		if err != nil {
			return err
		}
	}
	err = utils.ReplaceGeneratedContentGFV2(in.OutputCtrlPath)
	return nil
}

func (c CGenPb) doGenerateController(in generateControllerInput, controller generateCtrl) (err error) {
	var (
		folderPath = gfile.Join(in.OutputCtrlPath, controller.Package)
		filePath   = gfile.Join(folderPath, controller.Package+".go")
		isDirty    bool
	)
	if !gfile.Exists(folderPath) {
		if err = gfile.Mkdir(folderPath); err != nil {
			return err
		}
	}
	if !gfile.Exists(filePath) {
		templateContent := gstr.ReplaceByMap(controllerTemplate, g.MapStrStr{
			"{Name}":    controller.Name,
			"{Version}": controller.Version,
			"{Package}": controller.Package,
		})
		if err = gfile.PutContents(filePath, templateContent); err != nil {
			return err
		}
		isDirty = true
	}
	// Exist controller content.
	var ctrlContent string
	files, err := gfile.ScanDirFile(folderPath, "*.go", false)
	if err != nil {
		return err
	}
	for _, file := range files {
		if ctrlContent != "" {
			ctrlContent += "\n"
		}
		ctrlContent += gfile.GetContents(file)
	}
	// Generate method content.
	var generatedContent string
	for _, method := range controller.Methods {
		if gstr.Contains(ctrlContent, fmt.Sprintf(`%s(`, method.Name)) {
			continue
		}
		if generatedContent != "" {
			generatedContent += "\n"
		}
		methodTemplate := controllerMethodTemplate
		if method.Streaming {
			methodTemplate = controllerStreamMethodTemplate
		}
		generatedContent += gstr.ReplaceByMap(methodTemplate, g.MapStrStr{
			"{Definition}": method.Definition,
		})
	}
	if generatedContent != "" {
		err = gfile.PutContentsAppend(filePath, generatedContent)
		if err != nil {
			return err
		}
		isDirty = true
	}
	if isDirty {
		if err = ensureControllerImports(filePath, controller); err != nil {
			return err
		}
		utils.GoFmt(filePath)
	}
	return nil
}

// ensureControllerImports ensures the imports required by the controller template and by the
// controller method template, as GoFmt only removes the unused imports and never adds them.
//
// The imports of the method template are ensured as well even for a controller without any
// method, and the following GoFmt removes the ones that the generated content does not use,
// which keeps the same result as the imports.Process used to give.
func ensureControllerImports(filePath string, controller generateCtrl) error {
	var specs = make([]utils.ImportSpec, 0, 5)
	if controller.ApiImportPath != "" {
		// The alias is given explicitly, as the package name of the api package may differ
		// from the version segment of its import path.
		specs = append(specs, utils.ImportSpec{
			Name: controller.Version,
			Path: controller.ApiImportPath,
		})
	} else {
		// The output folder is outside of any go module, for example a temporary folder, so
		// the import path of the api package cannot be resolved: leave it imported by the
		// user instead of writing an invalid import path.
		mlog.Printf(
			`cannot resolve the import path of the api package for controller "%s", `+
				`the "%s" import should be added manually`,
			filePath, controller.Version,
		)
	}
	specs = append(specs,
		utils.ImportSpec{Path: controllerContextImportPath},
		utils.ImportSpec{Path: controllerGrpcxImportPath},
		utils.ImportSpec{Path: controllerGcodeImportPath},
		utils.ImportSpec{Path: controllerGerrorImportPath},
	)
	specs = append(specs, controller.Imports...)
	if _, err := utils.EnsureImports(filePath, specs...); err != nil {
		return gerror.Wrapf(err, `ensuring imports of controller file "%s" failed`, filePath)
	}
	return nil
}

// controllerStreamingDefinition qualifies local protobuf types and preserves imported selectors.
func controllerStreamingDefinition(line, version string, aliases map[string]string) (string, []string, error) {
	expression, err := parser.ParseExpr("interface { " + line + " }")
	if err != nil {
		return "", nil, err
	}
	var (
		method    = expression.(*ast.InterfaceType).Methods.List[0]
		signature = method.Type.(*ast.FuncType)
		packages  []string
		buffer    bytes.Buffer
	)
	for _, parameter := range signature.Params.List {
		// Imported message types and grpc's generic stream types keep their original qualifiers.
		ast.Inspect(parameter.Type, func(node ast.Node) bool {
			if selector, ok := node.(*ast.SelectorExpr); ok {
				if name, ok := selector.X.(*ast.Ident); ok {
					packages = append(packages, name.Name)
					if alias, ok := aliases[name.Name]; ok {
						name.Name = alias
					}
				}
			}
			return true
		})
		parameter.Type = qualifyControllerStreamType(parameter.Type, version)
	}
	if err = format.Node(&buffer, token.NewFileSet(), signature); err != nil {
		return "", nil, err
	}
	return method.Names[0].Name + strings.TrimPrefix(buffer.String(), "func"), packages, nil
}

// qualifyControllerStreamType adds the API alias to local messages and legacy stream interfaces.
func qualifyControllerStreamType(expression ast.Expr, version string) ast.Expr {
	switch expression := expression.(type) {
	case *ast.Ident:
		return &ast.SelectorExpr{X: ast.NewIdent(version), Sel: expression}
	case *ast.StarExpr:
		expression.X = qualifyControllerStreamType(expression.X, version)
	case *ast.IndexExpr:
		expression.Index = qualifyControllerStreamType(expression.Index, version)
	case *ast.IndexListExpr:
		for i, index := range expression.Indices {
			expression.Indices[i] = qualifyControllerStreamType(index, version)
		}
	}
	return expression
}
