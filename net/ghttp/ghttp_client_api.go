// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/container/gvar"

// Get is a convenience method for sending GET request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Get(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("GET", url, data...)
}

// Put is a convenience method for sending PUT request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Put(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PUT", url, data...)
}

// Post is a convenience method for sending POST request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Post(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("POST", url, data...)
}

// Delete is a convenience method for sending DELETE request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Delete(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("DELETE", url, data...)
}

// Head is a convenience method for sending HEAD request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Head(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("HEAD", url, data...)
}

// Patch is a convenience method for sending PATCH request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Patch(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PATCH", url, data...)
}

// Connect is a convenience method for sending CONNECT request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Connect(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("CONNECT", url, data...)
}

// Options is a convenience method for sending OPTIONS request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Options(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("OPTIONS", url, data...)
}

// Trace is a convenience method for sending TRACE request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func Trace(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("TRACE", url, data...)
}

// DoRequest is a convenience method for sending custom http method request.
// NOTE that remembers CLOSING the response object when it'll never be used.
func DoRequest(method, url string, data ...interface{}) (*ClientResponse, error) {
	return NewClient().DoRequest(method, url, data...)
}

// GetContent is a convenience method for sending GET request, which retrieves and returns
// the result content and automatically closes response object.
func GetContent(url string, data ...interface{}) string {
	return RequestContent("GET", url, data...)
}

// PutContent is a convenience method for sending PUT request, which retrieves and returns
// the result content and automatically closes response object.
func PutContent(url string, data ...interface{}) string {
	return RequestContent("PUT", url, data...)
}

// PostContent is a convenience method for sending POST request, which retrieves and returns
// the result content and automatically closes response object.
func PostContent(url string, data ...interface{}) string {
	return RequestContent("POST", url, data...)
}

// DeleteContent is a convenience method for sending DELETE request, which retrieves and returns
// the result content and automatically closes response object.
func DeleteContent(url string, data ...interface{}) string {
	return RequestContent("DELETE", url, data...)
}

// HeadContent is a convenience method for sending HEAD request, which retrieves and returns
// the result content and automatically closes response object.
func HeadContent(url string, data ...interface{}) string {
	return RequestContent("HEAD", url, data...)
}

// PatchContent is a convenience method for sending PATCH request, which retrieves and returns
// the result content and automatically closes response object.
func PatchContent(url string, data ...interface{}) string {
	return RequestContent("PATCH", url, data...)
}

// ConnectContent is a convenience method for sending CONNECT request, which retrieves and returns
// the result content and automatically closes response object.
func ConnectContent(url string, data ...interface{}) string {
	return RequestContent("CONNECT", url, data...)
}

// OptionsContent is a convenience method for sending OPTIONS request, which retrieves and returns
// the result content and automatically closes response object.
func OptionsContent(url string, data ...interface{}) string {
	return RequestContent("OPTIONS", url, data...)
}

// TraceContent is a convenience method for sending TRACE request, which retrieves and returns
// the result content and automatically closes response object.
func TraceContent(url string, data ...interface{}) string {
	return RequestContent("TRACE", url, data...)
}

// RequestContent is a convenience method for sending custom http method request, which
// retrieves and returns the result content and automatically closes response object.
func RequestContent(method string, url string, data ...interface{}) string {
	return NewClient().RequestContent(method, url, data...)
}

// GetBytes is a convenience method for sending GET request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func GetBytes(url string, data ...interface{}) []byte {
	return RequestBytes("GET", url, data...)
}

// PutBytes is a convenience method for sending PUT request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func PutBytes(url string, data ...interface{}) []byte {
	return RequestBytes("PUT", url, data...)
}

// PostBytes is a convenience method for sending POST request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func PostBytes(url string, data ...interface{}) []byte {
	return RequestBytes("POST", url, data...)
}

// DeleteBytes is a convenience method for sending DELETE request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func DeleteBytes(url string, data ...interface{}) []byte {
	return RequestBytes("DELETE", url, data...)
}

// HeadBytes is a convenience method for sending HEAD request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func HeadBytes(url string, data ...interface{}) []byte {
	return RequestBytes("HEAD", url, data...)
}

// PatchBytes is a convenience method for sending PATCH request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func PatchBytes(url string, data ...interface{}) []byte {
	return RequestBytes("PATCH", url, data...)
}

// ConnectBytes is a convenience method for sending CONNECT request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func ConnectBytes(url string, data ...interface{}) []byte {
	return RequestBytes("CONNECT", url, data...)
}

// OptionsBytes is a convenience method for sending OPTIONS request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func OptionsBytes(url string, data ...interface{}) []byte {
	return RequestBytes("OPTIONS", url, data...)
}

// TraceBytes is a convenience method for sending TRACE request, which retrieves and returns
// the result content as bytes and automatically closes response object.
func TraceBytes(url string, data ...interface{}) []byte {
	return RequestBytes("TRACE", url, data...)
}

// RequestBytes is a convenience method for sending custom http method request, which
// retrieves and returns the result content as bytes and automatically closes response object.
func RequestBytes(method string, url string, data ...interface{}) []byte {
	return NewClient().RequestBytes(method, url, data...)
}

// GetVar sends a GET request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func GetVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("GET", url, data...)
}

// PutVar sends a PUT request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func PutVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("PUT", url, data...)
}

// PostVar sends a POST request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func PostVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("POST", url, data...)
}

// DeleteVar sends a DELETE request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func DeleteVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("DELETE", url, data...)
}

// HeadVar sends a HEAD request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func HeadVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("HEAD", url, data...)
}

// PatchVar sends a PATCH request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func PatchVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("PATCH", url, data...)
}

// ConnectVar sends a CONNECT request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func ConnectVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("CONNECT", url, data...)
}

// OptionsVar sends a OPTIONS request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func OptionsVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("OPTIONS", url, data...)
}

// TraceVar sends a TRACE request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func TraceVar(url string, data ...interface{}) *gvar.Var {
	return RequestVar("TRACE", url, data...)
}

// RequestVar sends request using given HTTP method and data, retrieves converts the result
// to specified pointer. It reads and closes the response object internally automatically.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, et
func RequestVar(method string, url string, data ...interface{}) *gvar.Var {
	response, err := DoRequest(method, url, data...)
	if err != nil {
		return gvar.New(nil)
	}
	defer response.Close()
	return gvar.New(response.ReadAll())
}
