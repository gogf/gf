// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genpb

import (
	"context"
	"fmt"
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
	Methods       []generateCtrlMethod
}

type generateCtrlMethod struct {
	Name       string
	Definition string
}

const (
	// controllerContextImportPath is the import path of the `context` package used by the
	// generated controller methods.
	controllerContextImportPath = "context"
	// controllerGrpcxImportPath is the import path of the grpcx package used by the generated
	// controller registration function.
	controllerGrpcxImportPath = "github.com/gogf/gf/v2/net/grpcx"
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
	)
	_, err := gregex.ReplaceStringFuncMatch(
		`type (\w+)Server interface {([\s\S]+?)}`,
		content,
		func(match []string) string {
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
				ctrl.Methods = append(ctrl.Methods, generateCtrlMethod{
					Name:       gstr.Split(line, "(")[0],
					Definition: line,
				})
			}
			if len(ctrl.Methods) > 0 {
				controllers = append(controllers, ctrl)
			}
			return match[0]
		},
	)
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
		generatedContent += gstr.ReplaceByMap(controllerMethodTemplate, g.MapStrStr{
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
	if _, err := utils.EnsureImports(filePath, specs...); err != nil {
		return gerror.Wrapf(err, `ensuring imports of controller file "%s" failed`, filePath)
	}
	return nil
}
