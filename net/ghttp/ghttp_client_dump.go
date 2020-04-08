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

const dumpTextFormat = `+---------------------------------------------+
|               ghttp %s                |
+---------------------------------------------+
%s
%s
`

func ifDumpBody(contentType string) bool {
	// only dump body when contentType in.
	if gstr.Contains(contentType, "application/json") ||
		gstr.Contains(contentType, "application/xml") ||
		gstr.Contains(contentType, "multipart/form-data") ||
		gstr.Contains(contentType, "application/x-www-form-urlencoded") ||
		gstr.Contains(contentType, "text/plain") {
		return true
	}
	return false
}

// getRequestBody return request body string.
// if error occurs, return empty string
func getRequestBody(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if !ifDumpBody(contentType) {
		return ""
	}
	// must use this method for reading the body more than once
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

// getResponseBody return response body string.
// if error occurs, return empty string
func getResponseBody(resp *http.Response) string {
	contentType := resp.Header.Get("Content-Type")
	if !ifDumpBody(contentType) {
		return ""
	}
	bytesBody, errReadBody := ioutil.ReadAll(resp.Body)
	if errReadBody != nil {
		return ""
	}
	// for reading the body more than once
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bytesBody))
	return gconv.UnsafeBytesToStr(bytesBody)
}

func (r *ClientResponse) getRequest() *http.Request {
	if r.Response != nil && r.Request != nil {
		return r.Request
	}
	if r.req != nil {
		return r.req
	}
	return nil
}

// RawRequest dump request to raw string
func (r *ClientResponse) RawRequest() string {
	// this can be nil
	if r == nil {
		return ""
	}
	req := r.getRequest()
	if req == nil {
		return ""
	}
	// DumpRequestOut will write more header than DumpRequest, such as User-Agent.
	// read body using getRequestBody method.
	bs, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(dumpTextFormat, "REQUEST ", gconv.UnsafeBytesToStr(bs), getRequestBody(req))
}

// RawResponse dump response to raw string
func (r *ClientResponse) RawResponse() string {
	// this can be nil
	if r == nil || r.Response == nil {
		return ""
	}
	// read body using getResponseBody method.
	bs, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(dumpTextFormat, "RESPONSE", gconv.UnsafeBytesToStr(bs), getResponseBody(r.Response))
}

// Raw dump request and response string
func (r *ClientResponse) Raw() string {
	return fmt.Sprintf("%s\n%s", r.RawRequest(), r.RawResponse())
}
