// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Get(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodGet, url, data...)
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Put(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPut, url, data...)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Post(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPost, url, data...)
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Delete(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodDelete, url, data...)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Head(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodHead, url, data...)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Patch(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPatch, url, data...)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Connect(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodConnect, url, data...)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Options(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodOptions, url, data...)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Trace(ctx context.Context, url string, data ...interface{}) (*Response, error) {
	return c.DoRequest(ctx, http.MethodTrace, url, data...)
}

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and Client.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
//
// To make a request with a specified context.Context, use NewRequestWithContext
// and Client.Do.
// Deprecated, use Post instead.
func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error) {
	return nil, gerror.NewCode(
		gcode.CodeNotSupported,
		`PostForm is not supported, please use Post instead`,
	)
}

// DoRequest sends request with given HTTP method and data and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
//
// Note that it uses "multipart/form-data" as its Content-Type if it contains file uploading,
// else it uses "application/x-www-form-urlencoded". It also automatically detects the post
// content for JSON format, and for that it automatically sets the Content-Type as
// "application/json".
func (c *Client) DoRequest(ctx context.Context, method, url string, data ...interface{}) (resp *Response, err error) {
	req, err := c.prepareRequest(ctx, method, url, data...)
	if err != nil {
		return nil, err
	}

	// Client middleware.
	if len(c.middlewareHandler) > 0 {
		mdlHandlers := make([]HandlerFunc, 0, len(c.middlewareHandler)+1)
		mdlHandlers = append(mdlHandlers, c.middlewareHandler...)
		mdlHandlers = append(mdlHandlers, func(cli *Client, r *http.Request) (*Response, error) {
			return cli.callRequest(r)
		})
		ctx = context.WithValue(req.Context(), clientMiddlewareKey, &clientMiddleware{
			client:       c,
			handlers:     mdlHandlers,
			handlerIndex: -1,
		})
		req = req.WithContext(ctx)
		resp, err = c.Next(req)
	} else {
		resp, err = c.callRequest(req)
	}
	return resp, err
}

// prepareRequest verifies request parameters, builds and returns http request.
func (c *Client) prepareRequest(ctx context.Context, method, url string, data ...interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	if len(c.prefix) > 0 {
		url = c.prefix + gstr.Trim(url)
	}
	if !gstr.ContainsI(url, httpProtocolName) {
		url = httpProtocolName + `://` + url
	}
	var params string
	if len(data) > 0 {
		switch c.header[httpHeaderContentType] {
		case httpHeaderContentTypeJson:
			switch data[0].(type) {
			case string, []byte:
				params = gconv.String(data[0])
			default:
				if b, err := json.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}

		case httpHeaderContentTypeXml:
			switch data[0].(type) {
			case string, []byte:
				params = gconv.String(data[0])
			default:
				if b, err := gjson.New(data[0]).ToXml(); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			}
		default:
			params = httputil.BuildParams(data[0])
		}
	}
	if method == http.MethodGet {
		var bodyBuffer *bytes.Buffer
		if params != "" {
			switch c.header[httpHeaderContentType] {
			case
				httpHeaderContentTypeJson,
				httpHeaderContentTypeXml:
				bodyBuffer = bytes.NewBuffer([]byte(params))
			default:
				// It appends the parameters to the url
				// if http method is GET and Content-Type is not specified.
				if gstr.Contains(url, "?") {
					url = url + "&" + params
				} else {
					url = url + "?" + params
				}
				bodyBuffer = bytes.NewBuffer(nil)
			}
		} else {
			bodyBuffer = bytes.NewBuffer(nil)
		}
		if req, err = http.NewRequest(method, url, bodyBuffer); err != nil {
			err = gerror.Wrapf(err, `http.NewRequest failed with method "%s" and URL "%s"`, method, url)
			return nil, err
		}
	} else {
		if strings.Contains(params, httpParamFileHolder) {
			// File uploading request.
			var (
				buffer = bytes.NewBuffer(nil)
				writer = multipart.NewWriter(buffer)
			)
			for _, item := range strings.Split(params, "&") {
				array := strings.Split(item, "=")
				if len(array[1]) > 6 && strings.Compare(array[1][0:6], httpParamFileHolder) == 0 {
					path := array[1][6:]
					if !gfile.Exists(path) {
						return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" does not exist`, path)
					}
					var (
						file          io.Writer
						formFileName  = gfile.Basename(path)
						formFieldName = array[0]
					)
					if file, err = writer.CreateFormFile(formFieldName, formFileName); err != nil {
						err = gerror.Wrapf(err, `CreateFormFile failed with "%s", "%s"`, formFieldName, formFileName)
						return nil, err
					} else {
						var f *os.File
						if f, err = gfile.Open(path); err != nil {
							return nil, err
						}
						if _, err = io.Copy(file, f); err != nil {
							err = gerror.Wrapf(err, `io.Copy failed from "%s" to form "%s"`, path, formFieldName)
							_ = f.Close()
							return nil, err
						}
						_ = f.Close()
					}
				} else {
					var (
						fieldName  = array[0]
						fieldValue = array[1]
					)
					if err = writer.WriteField(fieldName, fieldValue); err != nil {
						err = gerror.Wrapf(err, `write form field failed with "%s", "%s"`, fieldName, fieldValue)
						return nil, err
					}
				}
			}
			// Close finishes the multipart message and writes the trailing
			// boundary end line to the output.
			if err = writer.Close(); err != nil {
				err = gerror.Wrapf(err, `form writer close failed`)
				return nil, err
			}

			if req, err = http.NewRequest(method, url, buffer); err != nil {
				err = gerror.Wrapf(err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url)
				return nil, err
			} else {
				req.Header.Set(httpHeaderContentType, writer.FormDataContentType())
			}
		} else {
			// Normal request.
			paramBytes := []byte(params)
			if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
				err = gerror.Wrapf(err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url)
				return nil, err
			} else {
				if v, ok := c.header[httpHeaderContentType]; ok {
					// Custom Content-Type.
					req.Header.Set(httpHeaderContentType, v)
				} else if len(paramBytes) > 0 {
					if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
						// Auto-detecting and setting the post content format: JSON.
						req.Header.Set(httpHeaderContentType, httpHeaderContentTypeJson)
					} else if gregex.IsMatchString(httpRegexParamJson, params) {
						// If the parameters passed like "name=value", it then uses form type.
						req.Header.Set(httpHeaderContentType, httpHeaderContentTypeForm)
					}
				}
			}
		}
	}

	// Context.
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			if len(req.Header.Get(httpHeaderContentType)) < 0 {
				req.Header.Set(k, v)
			}
		}
	}
	// It's necessary set the req.Host if you want to custom the host value of the request.
	// It uses the "Host" value from header if it's not empty.
	if reqHeaderHost := req.Header.Get(httpHeaderHost); reqHeaderHost != "" {
		req.Host = reqHeaderHost
	}
	// Custom Cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set(httpHeaderCookie, headerCookie)
		}
	}
	// HTTP basic authentication.
	if len(c.authUser) > 0 {
		req.SetBasicAuth(c.authUser, c.authPass)
	}
	return req, nil
}

// callRequest sends request with give http.Request, and returns the responses object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) callRequest(req *http.Request) (resp *Response, err error) {
	resp = &Response{
		request: req,
	}
	// Dump feature.
	// The request body can be reused for dumping
	// raw HTTP request-response procedure.
	reqBodyContent, _ := ioutil.ReadAll(req.Body)
	resp.requestBody = reqBodyContent
	req.Body = utils.NewReadCloser(reqBodyContent, false)
	for {
		if resp.Response, err = c.Do(req); err != nil {
			err = gerror.Wrapf(err, `request failed`)
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				_ = resp.Response.Body.Close()
			}
			if c.retryCount > 0 {
				c.retryCount--
				time.Sleep(c.retryInterval)
			} else {
				// return resp, err
				break
			}
		} else {
			break
		}
	}
	return resp, err
}
