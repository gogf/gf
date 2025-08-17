// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package httpclient provides http client used for SDK.
package httpclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gurl"
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
	Handler
}

// New creates and returns a http client for SDK.
func New(config Config) *Client {
	client := config.Client
	if client == nil {
		client = gclient.New()
	}
	handler := config.Handler
	if handler == nil {
		handler = NewDefaultHandler(config.Logger, config.RawDump)
	}
	if !gstr.HasPrefix(config.URL, "http") {
		config.URL = fmt.Sprintf("http://%s", config.URL)
	}
	return &Client{
		Client:  client.Prefix(config.URL),
		Handler: handler,
	}
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
		return c.HandleResponse(ctx, result, res)
	}
}

// Get sends a request using GET method.
func (c *Client) Get(ctx context.Context, path string, in, out interface{}) error {
	// TODO: Path params will also be built in urlParams, not graceful now.
	if urlParams := ghttp.BuildParams(in); urlParams != "" && urlParams != "{}" {
		path += "?" + urlParams
	}
	res, err := c.ContentJson().Get(ctx, c.handlePath(path, in))
	if err != nil {
		return gerror.Wrap(err, `http request failed`)
	}
	return c.HandleResponse(ctx, res, out)
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
