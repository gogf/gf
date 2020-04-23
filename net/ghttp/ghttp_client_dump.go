// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// dumpTextFormat is the format of the dumped raw string
const dumpTextFormat = `+---------------------------------------------+
|               ghttp %s                |
+---------------------------------------------+
%s
%s
`

// ifDumpBody determine whether to output body according to content-type
func ifDumpBody(contentType string) bool {
	// the body should not be output when the body is html or stream.
	if gstr.Contains(contentType, "application/json") ||
		gstr.Contains(contentType, "application/xml") ||
		gstr.Contains(contentType, "multipart/form-data") ||
		gstr.Contains(contentType, "application/x-www-form-urlencoded") ||
		gstr.Contains(contentType, "text/plain") {
		return true
	}
	return false
}

// getRequestBody returns the raw text of the request body.
func getRequestBody(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if !ifDumpBody(contentType) {
		return ""
	}
	// so that the request body can be read again.
	bodyReader, errGetBody := req.GetBody()
	if errGetBody != nil {
		return ""
	}
	bytesBody, errReadBody := ioutil.ReadAll(bodyReader)
	if errReadBody != nil {
		return ""
	}
	return gconv.UnsafeBytesToStr(bytesBody)
}

// getResponseBody returns the text of the response body.
func getResponseBody(resp *http.Response) string {
	contentType := resp.Header.Get("Content-Type")
	if !ifDumpBody(contentType) {
		return ""
	}
	bytesBody, errReadBody := ioutil.ReadAll(resp.Body)
	if errReadBody != nil {
		return ""
	}
	// so that the response body can be read again.
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bytesBody))
	return gconv.UnsafeBytesToStr(bytesBody)
}

// getRequest returns the request related to the response.
// will return the copy of request when the request failed.
func (r *ClientResponse) getRequest() *http.Request {
	if r.Response != nil && r.Request != nil {
		return r.Request
	}
	// r.req is the copy of request when the http request failed.
	if r.req != nil {
		return r.req
	}
	return nil
}

// RawRequest returns the raw text of the request.
func (r *ClientResponse) RawRequest() string {
	// ClientResponse can be nil.
	if r == nil {
		return ""
	}
	req := r.getRequest()
	if req == nil {
		return ""
	}
	// DumpRequestOut writes more request headers than DumpRequest, such as User-Agent.
	bs, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(dumpTextFormat, "REQUEST ", gconv.UnsafeBytesToStr(bs), getRequestBody(req))
}

// RawResponse returns the raw text of the response.
func (r *ClientResponse) RawResponse() string {
	// ClientResponse can be nil.
	if r == nil || r.Response == nil {
		return ""
	}
	bs, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(dumpTextFormat, "RESPONSE", gconv.UnsafeBytesToStr(bs), getResponseBody(r.Response))
}

// Raw returns the raw text of the request and the response.
func (r *ClientResponse) Raw() string {
	return fmt.Sprintf("%s\n%s", r.RawRequest(), r.RawResponse())
}
