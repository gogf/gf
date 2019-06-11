// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// HTTP客户端请求.

package ghttp

func Get(url string) (*ClientResponse, error) {
    return DoRequest("GET", url)
}

func Put(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("PUT", url, data...)
}

func Post(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("POST", url, data...)
}

func Delete(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("DELETE", url, data...)
}

func Head(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("HEAD", url, data...)
}

func Patch(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("PATCH", url, data...)
}

func Connect(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("CONNECT", url, data...)
}

func Options(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("OPTIONS", url, data...)
}

func Trace(url string, data...interface{}) (*ClientResponse, error) {
    return DoRequest("TRACE", url, data...)
}

// 该方法支持二进制提交数据
func DoRequest(method, url string, data...interface{}) (*ClientResponse, error) {
    return NewClient().DoRequest(method, url, data...)
}

// GET请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func GetContent(url string, data...interface{}) string {
    return RequestContent("GET", url, data...)
}

// PUT请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func PutContent(url string, data...interface{}) string {
    return RequestContent("PUT", url, data...)
}

// POST请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func PostContent(url string, data...interface{}) string {
    return RequestContent("POST", url, data...)
}

// DELETE请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func DeleteContent(url string, data...interface{}) string {
    return RequestContent("DELETE", url, data...)
}

func HeadContent(url string, data...interface{}) string {
    return RequestContent("HEAD", url, data...)
}

func PatchContent(url string, data...interface{}) string {
    return RequestContent("PATCH", url, data...)
}

func ConnectContent(url string, data...interface{}) string {
    return RequestContent("CONNECT", url, data...)
}

func OptionsContent(url string, data...interface{}) string {
    return RequestContent("OPTIONS", url, data...)
}

func TraceContent(url string, data...interface{}) string {
    return RequestContent("TRACE", url, data...)
}

// 请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func RequestContent(method string, url string, data...interface{}) string {
    return NewClient().DoRequestContent(method, url, data...)
}

