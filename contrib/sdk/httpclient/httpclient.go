// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package httpclient provides http client used for SDK.
package httpclient

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

// Client is a http client for SDK.
type Client struct {
	*gclient.Client
	config Config
}

// New creates and returns a http client for SDK.
func New(config Config) *Client {
	return &Client{
		Client: config.Client,
		config: config,
	}
}

func (c *Client) handleResponse(ctx context.Context, res *gclient.Response, out interface{}) error {
	if c.config.RawDump {
		c.config.Logger.Debugf(ctx, "raw request&response:\n%s", res.Raw())
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

// Request sends request to service by struct object `req`, and receives response to struct object `res`.
func (c *Client) Request(ctx context.Context, req, res interface{}) error {
	var (
		method = gmeta.Get(req, gtag.Method).String()
		path   = gmeta.Get(req, gtag.Path).String()
	)
	switch gstr.ToUpper(method) {
	case http.MethodGet:
		return c.Get(ctx, path, req, res)

	default:
		result, err := c.ContentJson().DoRequest(ctx, method, c.handlePath(path, req), req)
		if err != nil {
			return err
		}
		return c.handleResponse(ctx, result, res)
	}
}

// Get sends a request using GET method.
func (c *Client) Get(ctx context.Context, path string, in, out interface{}) error {
	urlParams := ghttp.BuildParams(in)
	if urlParams != "" {
		path += "?" + ghttp.BuildParams(in)
	}
	res, err := c.ContentJson().Get(ctx, c.handlePath(path, in))
	if err != nil {
		return gerror.Wrap(err, `http request failed`)
	}
	return c.handleResponse(ctx, res, out)
}

func (c *Client) handlePath(path string, in interface{}) string {
	if gstr.Contains(path, "{") {
		data := gconv.MapStrStr(in)
		path, _ = gregex.ReplaceStringFuncMatch(`\{(\w+)\}`, path, func(match []string) string {
			if v, ok := data[match[1]]; ok {
				return gurl.Encode(v)
			}
			return match[1]
		})
	}
	return path
}
