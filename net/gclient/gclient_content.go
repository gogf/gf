// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import "context"

// GetContent is a convenience method for sending GET request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) GetContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodGet, url, data...))
}

// PutContent is a convenience method for sending PUT request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PutContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodPut, url, data...))
}

// PostContent is a convenience method for sending POST request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PostContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodPost, url, data...))
}

// DeleteContent is a convenience method for sending DELETE request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) DeleteContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodDelete, url, data...))
}

// HeadContent is a convenience method for sending HEAD request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) HeadContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodHead, url, data...))
}

// PatchContent is a convenience method for sending PATCH request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PatchContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodPatch, url, data...))
}

// ConnectContent is a convenience method for sending CONNECT request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) ConnectContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodConnect, url, data...))
}

// OptionsContent is a convenience method for sending OPTIONS request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) OptionsContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodOptions, url, data...))
}

// TraceContent is a convenience method for sending TRACE request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) TraceContent(ctx context.Context, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, httpMethodTrace, url, data...))
}

// RequestContent is a convenience method for sending custom http method request, which
// retrieves and returns the result content and automatically closes response object.
func (c *Client) RequestContent(ctx context.Context, method string, url string, data ...interface{}) string {
	return string(c.RequestBytes(ctx, method, url, data...))
}
