// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genpb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type generateControllerInput struct {
	OutputApiPath  string
	OutputCtrlPath string
}

type generateCtrl struct {
	Name    string
	Package string
	Version string
	Methods []generateCtrlMethod
}

type generateCtrlMethod struct {
	Name       string
	Definition string
}

const (
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
				Name:    match[1],
				Package: gfile.Basename(gfile.Dir(gfile.Dir(filePath))),
				Version: gfile.Basename(gfile.Dir(filePath)),
				Methods: make([]generateCtrlMethod, 0),
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
		utils.GoFmt(filePath)
	}
	return nil
}
