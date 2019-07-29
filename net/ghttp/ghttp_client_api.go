// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func Get(url string) (*ClientResponse, error) {
	return DoRequest("GET", url)
}

func Put(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PUT", url, data...)
}

func Post(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("POST", url, data...)
}

func Delete(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("DELETE", url, data...)
}

func Head(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("HEAD", url, data...)
}

func Patch(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("PATCH", url, data...)
}

func Connect(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("CONNECT", url, data...)
}

func Options(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("OPTIONS", url, data...)
}

func Trace(url string, data ...interface{}) (*ClientResponse, error) {
	return DoRequest("TRACE", url, data...)
}

func DoRequest(method, url string, data ...interface{}) (*ClientResponse, error) {
	return NewClient().DoRequest(method, url, data...)
}

func GetContent(url string, data ...interface{}) string {
	return RequestContent("GET", url, data...)
}

func PutContent(url string, data ...interface{}) string {
	return RequestContent("PUT", url, data...)
}

func PostContent(url string, data ...interface{}) string {
	return RequestContent("POST", url, data...)
}

func DeleteContent(url string, data ...interface{}) string {
	return RequestContent("DELETE", url, data...)
}

func HeadContent(url string, data ...interface{}) string {
	return RequestContent("HEAD", url, data...)
}

func PatchContent(url string, data ...interface{}) string {
	return RequestContent("PATCH", url, data...)
}

func ConnectContent(url string, data ...interface{}) string {
	return RequestContent("CONNECT", url, data...)
}

func OptionsContent(url string, data ...interface{}) string {
	return RequestContent("OPTIONS", url, data...)
}

func TraceContent(url string, data ...interface{}) string {
	return RequestContent("TRACE", url, data...)
}

func RequestContent(method string, url string, data ...interface{}) string {
	return NewClient().RequestContent(method, url, data...)
}

func GetBytes(url string, data ...interface{}) []byte {
	return RequestBytes("GET", url, data...)
}

func PutBytes(url string, data ...interface{}) []byte {
	return RequestBytes("PUT", url, data...)
}

func PostBytes(url string, data ...interface{}) []byte {
	return RequestBytes("POST", url, data...)
}

func DeleteBytes(url string, data ...interface{}) []byte {
	return RequestBytes("DELETE", url, data...)
}

func HeadBytes(url string, data ...interface{}) []byte {
	return RequestBytes("HEAD", url, data...)
}

func PatchBytes(url string, data ...interface{}) []byte {
	return RequestBytes("PATCH", url, data...)
}

func ConnectBytes(url string, data ...interface{}) []byte {
	return RequestBytes("CONNECT", url, data...)
}

func OptionsBytes(url string, data ...interface{}) []byte {
	return RequestBytes("OPTIONS", url, data...)
}

func TraceBytes(url string, data ...interface{}) []byte {
	return RequestBytes("TRACE", url, data...)
}

func RequestBytes(method string, url string, data ...interface{}) []byte {
	return NewClient().RequestBytes(method, url, data...)
}
