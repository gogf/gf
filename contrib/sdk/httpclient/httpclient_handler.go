// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package httpclient

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

// Handler is the interface for http response handling.
type Handler interface {
	// HandleResponse handles the http response and transforms its body to the specified object.
	// The parameter `out` specifies the object that the response body is transformed to.
	HandleResponse(ctx context.Context, res *gclient.Response, out interface{}) error
}

// DefaultHandler handle ghttp.DefaultHandlerResponse of json format.
type DefaultHandler struct {
	Logger  *glog.Logger
	RawDump bool
}

func NewDefaultHandler(logger *glog.Logger, rawRump bool) *DefaultHandler {
	if rawRump && logger == nil {
		logger = g.Log()
	}
	return &DefaultHandler{
		Logger:  logger,
		RawDump: rawRump,
	}
}

func (h DefaultHandler) HandleResponse(ctx context.Context, res *gclient.Response, out interface{}) error {
	defer res.Close()
	if h.RawDump {
		h.Logger.Debugf(ctx, "raw request&response:\n%s", res.Raw())
	}
	var (
		responseBytes = res.ReadAll()
		result        = ghttp.DefaultHandlerResponse{
			Data: out,
		}
	)
	if !json.Valid(responseBytes) {
		return gerror.Newf(`invalid response content: %s`, responseBytes)
	}
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return gerror.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	if result.Code != gcode.CodeOK.Code() {
		return gerror.NewCode(
			gcode.New(result.Code, result.Message, nil),
			result.Message,
		)
	}
	return nil
}
