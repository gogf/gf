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
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Get(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodGet, url, data...)
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Put(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPut, url, data...)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Post(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPost, url, data...)
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Delete(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodDelete, url, data...)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Head(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodHead, url, data...)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Patch(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPatch, url, data...)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Connect(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodConnect, url, data...)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Options(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodOptions, url, data...)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Trace(ctx context.Context, url string, data ...any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodTrace, url, data...)
}

// PostForm is different from net/http.PostForm.
// It's a wrapper of Post method, which sets the Content-Type as "multipart/form-data;".
// and It will automatically set boundary characters for the request body and Content-Type.
//
// It's Seem like the following case:
//
// Content-Type: multipart/form-data; boundary=----Boundarye4Ghaog6giyQ9ncN
//
// And form data is like:
// ------Boundarye4Ghaog6giyQ9ncN
// Content-Disposition: form-data; name="checkType"
//
// none
//
// It's used for sending form data.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) PostForm(ctx context.Context, url string, data map[string]string) (resp *Response, err error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range data {
		err := w.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return c.ContentType(w.FormDataContentType()).Post(ctx, url, body)
}

// DoRequest sends request with given HTTP method and data and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
//
// Note that it uses "multipart/form-data" as its Content-Type if it contains file uploading,
// else it uses "application/x-www-form-urlencoded". It also automatically detects the post
// content for JSON format, and for that it automatically sets the Content-Type as
// "application/json".
func (c *Client) DoRequest(
	ctx context.Context, method, url string, data ...any,
) (resp *Response, err error) {
	var requestStartTime = gtime.Now()
	req, err := c.prepareRequest(ctx, method, url, data...)
	if err != nil {
		return nil, err
	}

	// Metrics.
	c.handleMetricsBeforeRequest(req)
	defer c.handleMetricsAfterRequestDone(req, requestStartTime)

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
	if resp != nil && resp.Response != nil {
		req.Response = resp.Response
	}
	return resp, err
}

// prepareRequest verifies request parameters, builds and returns http request.
func (c *Client) prepareRequest(ctx context.Context, method, url string, data ...any) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	url = c.prepareRequestURL(url)
	var (
		params             string
		mediaType          string
		allowFileUploading = true
	)
	if len(data) > 0 {
		mediaType = c.resolveRequestMediaType()
		if params, allowFileUploading, err = c.encodeRequestParams(data[0], mediaType); err != nil {
			return nil, err
		}
	}
	switch {
	case method == http.MethodGet:
		req, err = c.newGetRequest(method, url, params, mediaType)
	case allowFileUploading && strings.Contains(params, httpParamFileHolder):
		req, err = c.newMultipartRequest(method, url, params)
	default:
		req, err = c.newNormalRequest(method, url, params)
	}
	if err != nil {
		return nil, err
	}
	return c.applyRequestOptions(req, ctx), nil
}

// urlSchemeRegex matches urls that already carry a scheme, like "http://", "https://" or "ws://".
// It anchors at the beginning and requires the "://" part, so that neither urls containing
// "http" inside (eg. "myhttpservice.com") nor "host:port" inputs (eg. "localhost:8000")
// are mistakenly treated as schemed ones.
var urlSchemeRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)

// prepareRequestURL returns the url with the client prefix and the protocol completed.
func (c *Client) prepareRequestURL(url string) string {
	if len(c.prefix) > 0 {
		url = c.prefix + gstr.Trim(url)
	}
	if !urlSchemeRegex.MatchString(url) {
		url = httpProtocolName + `://` + url
	}
	return url
}

// resolveRequestMediaType parses and returns the media type from the client custom
// Content-Type header. It falls back to the raw header value if parsing fails.
func (c *Client) resolveRequestMediaType() string {
	mediaType, _, err := mime.ParseMediaType(c.header[httpHeaderContentType])
	if err != nil {
		// Fallback: use the raw header value if parsing fails.
		mediaType = c.header[httpHeaderContentType]
	}
	return mediaType
}

// encodeRequestParams serializes the request data into params string according to
// the given media type. It also returns whether the params allow file uploading.
func (c *Client) encodeRequestParams(data any, mediaType string) (string, bool, error) {
	switch mediaType {
	case httpHeaderContentTypeJson:
		switch data.(type) {
		case string, []byte:
			return gconv.String(data), false, nil
		default:
			b, err := json.Marshal(data)
			if err != nil {
				return "", false, err
			}
			return string(b), false, nil
		}

	case httpHeaderContentTypeXml:
		switch data.(type) {
		case string, []byte:
			return gconv.String(data), false, nil
		default:
			b, err := gjson.New(data).ToXml()
			if err != nil {
				return "", false, err
			}
			return string(b), false, nil
		}

	default:
		return httputil.BuildParams(data, c.noUrlEncode), true, nil
	}
}

// newGetRequest builds and returns a GET request with given params.
// Note that it appends the params to the url if the media type is not json or xml.
func (c *Client) newGetRequest(method, url, params, mediaType string) (req *http.Request, err error) {
	var bodyBuffer *bytes.Buffer
	if params != "" {
		switch mediaType {
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
	return req, nil
}

// newMultipartRequest builds and returns a multipart request for form data which
// contains file uploading items in format "@file:path".
func (c *Client) newMultipartRequest(method, url, params string) (req *http.Request, err error) {
	var (
		buffer          = bytes.NewBuffer(nil)
		writer          = multipart.NewWriter(buffer)
		isFileUploading = false
	)
	for _, item := range strings.Split(params, "&") {
		var isFile bool
		if isFile, err = writeMultipartItem(writer, item); err != nil {
			return nil, err
		}
		if isFile {
			isFileUploading = true
		}
	}
	// Close finishes the multipart message and writes the trailing
	// boundary end line to the output.
	if err = writer.Close(); err != nil {
		return nil, gerror.Wrapf(err, `form writer close failed`)
	}

	if req, err = http.NewRequest(method, url, buffer); err != nil {
		return nil, gerror.Wrapf(
			err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url,
		)
	}
	if isFileUploading {
		req.Header.Set(httpHeaderContentType, writer.FormDataContentType())
	}
	return req, nil
}

// writeMultipartItem writes one "key=value" item into the multipart writer.
// It returns whether this item is a file uploading one.
func writeMultipartItem(writer *multipart.Writer, item string) (isFile bool, err error) {
	array := strings.SplitN(item, "=", 2)
	if len(array) < 2 {
		return false, nil
	}
	if len(array[1]) > 6 && strings.Compare(array[1][0:6], httpParamFileHolder) == 0 {
		path := array[1][6:]
		if !gfile.Exists(path) {
			return false, gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" does not exist`, path)
		}
		var (
			file          io.Writer
			formFileName  = gfile.Basename(path)
			formFieldName = array[0]
		)
		// it sets post content type as `application/octet-stream`
		if file, err = writer.CreateFormFile(formFieldName, formFileName); err != nil {
			return false, gerror.Wrapf(
				err, `CreateFormFile failed with "%s", "%s"`, formFieldName, formFileName,
			)
		}
		var f *os.File
		if f, err = gfile.Open(path); err != nil {
			return false, err
		}
		if _, err = io.Copy(file, f); err != nil {
			_ = f.Close()
			return false, gerror.Wrapf(
				err, `io.Copy failed from "%s" to form "%s"`, path, formFieldName,
			)
		}
		if err = f.Close(); err != nil {
			return false, gerror.Wrapf(err, `close file descriptor failed for "%s"`, path)
		}
		return true, nil
	}
	var (
		fieldName  = array[0]
		fieldValue = array[1]
	)
	// Decode URL-encoded field name and value.
	// If decoding fails, use the original value.
	if v, err := gurl.Decode(fieldName); err == nil {
		fieldName = v
	}
	if v, err := gurl.Decode(fieldValue); err == nil {
		fieldValue = v
	}
	if err = writer.WriteField(fieldName, fieldValue); err != nil {
		return false, gerror.Wrapf(
			err, `write form field failed with "%s", "%s"`, fieldName, fieldValue,
		)
	}
	return false, nil
}

// newNormalRequest builds and returns a request with params as its body. It automatically
// detects and sets the Content-Type in json or form format if the client has no custom one.
func (c *Client) newNormalRequest(method, url, params string) (req *http.Request, err error) {
	paramBytes := []byte(params)
	if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
		err = gerror.Wrapf(err, `http.NewRequest failed for method "%s" and URL "%s"`, method, url)
		return nil, err
	}
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
	return req, nil
}

// applyRequestOptions applies the client custom options to the request:
// context, custom header, host, cookie and basic authentication.
func (c *Client) applyRequestOptions(req *http.Request, ctx context.Context) *http.Request {
	// Context.
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
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
	return req
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
	reqBodyContent, _ := io.ReadAll(req.Body)
	resp.requestBody = reqBodyContent
	for {
		req.Body = utils.NewReadCloser(reqBodyContent, false)
		if resp.Response, err = c.Do(req); err != nil {
			err = gerror.Wrapf(err, `request failed`)
			// The response might not be nil when err != nil.
			if resp.Response != nil {
				_ = resp.Body.Close()
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
