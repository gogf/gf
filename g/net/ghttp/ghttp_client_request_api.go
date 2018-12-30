// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// HTTP客户端请求.

package ghttp

func Get(url string) (*ClientResponse, error) {
    return DoRequest("GET", url, []byte(""))
}

func Put(url, data string) (*ClientResponse, error) {
    return DoRequest("PUT", url, []byte(data))
}

func Post(url, data string) (*ClientResponse, error) {
    return DoRequest("POST", url, []byte(data))
}

func Delete(url, data string) (*ClientResponse, error) {
    return DoRequest("DELETE", url, []byte(data))
}

func Head(url, data string) (*ClientResponse, error) {
    return DoRequest("HEAD", url, []byte(data))
}

func Patch(url, data string) (*ClientResponse, error) {
    return DoRequest("PATCH", url, []byte(data))
}

func Connect(url, data string) (*ClientResponse, error) {
    return DoRequest("CONNECT", url, []byte(data))
}

func Options(url, data string) (*ClientResponse, error) {
    return DoRequest("OPTIONS", url, []byte(data))
}

func Trace(url, data string) (*ClientResponse, error) {
    return DoRequest("TRACE", url, []byte(data))
}

// 该方法支持二进制提交数据
func DoRequest(method, url string, data []byte) (*ClientResponse, error) {
    return NewClient().DoRequest(method, url, data)
}

// GET请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func GetContent(url string, data...string) string {
    return RequestContent("GET", url, data...)
}

// PUT请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func PutContent(url string, data...string) string {
    return RequestContent("PUT", url, data...)
}

// POST请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func PostContent(url string, data...string) string {
    return RequestContent("POST", url, data...)
}

// DELETE请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func DeleteContent(url string, data...string) string {
    return RequestContent("DELETE", url, data...)
}

func HeadContent(url string, data...string) string {
    return RequestContent("HEAD", url, data...)
}

func PatchContent(url string, data...string) string {
    return RequestContent("PATCH", url, data...)
}

func ConnectContent(url string, data...string) string {
    return RequestContent("CONNECT", url, data...)
}

func OptionsContent(url string, data...string) string {
    return RequestContent("OPTIONS", url, data...)
}

func TraceContent(url string, data...string) string {
    return RequestContent("TRACE", url, data...)
}

// 请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func RequestContent(method string, url string, data...string) string {
    return NewClient().DoRequestContent(method, url, data...)
}

