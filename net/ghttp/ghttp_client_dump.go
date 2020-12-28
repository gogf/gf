// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"github.com/gogf/gf/internal/utils"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gogf/gf/util/gconv"
)

// dumpTextFormat is the format of the dumped raw string
const dumpTextFormat = `+---------------------------------------------+
|                   %s                  |
+---------------------------------------------+
%s
%s
`

// getResponseBody returns the text of the response body.
func getResponseBody(res *http.Response) string {
	if res.Body == nil {
		return ""
	}
	bodyContent, _ := ioutil.ReadAll(res.Body)
	res.Body = utils.NewReadCloser(bodyContent, true)
	return gconv.UnsafeBytesToStr(bodyContent)
}

// RawRequest returns the raw content of the request.
func (r *ClientResponse) RawRequest() string {
	// ClientResponse can be nil.
	if r == nil {
		return ""
	}
	if r.request == nil {
		return ""
	}
	// DumpRequestOut writes more request headers than DumpRequest, such as User-Agent.
	bs, err := httputil.DumpRequestOut(r.request, false)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		dumpTextFormat,
		"REQUEST ",
		gconv.UnsafeBytesToStr(bs),
		r.requestBody,
	)
}

// RawResponse returns the raw content of the response.
func (r *ClientResponse) RawResponse() string {
	// ClientResponse can be nil.
	if r == nil || r.Response == nil {
		return ""
	}
	bs, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		dumpTextFormat,
		"RESPONSE",
		gconv.UnsafeBytesToStr(bs),
		getResponseBody(r.Response),
	)
}

// Raw returns the raw text of the request and the response.
func (r *ClientResponse) Raw() string {
	return fmt.Sprintf("%s\n%s", r.RawRequest(), r.RawResponse())
}

// RawDump outputs the raw text of the request and the response to stdout.
func (r *ClientResponse) RawDump() {
	fmt.Println(r.Raw())
}
