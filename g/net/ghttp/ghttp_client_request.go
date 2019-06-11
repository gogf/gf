// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// HTTP客户端请求.

package ghttp

import (
    "bytes"
	"crypto/tls"
	"encoding/json"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/os/gfile"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "strings"
    "time"
)

// http客户端
type Client struct {
    http.Client                     // 底层http client对象
    header        map[string]string // HEADER信息Map
    cookies       map[string]string // 自定义COOKIE
    prefix        string            // 设置请求的URL前缀
    authUser      string            // HTTP基本权限设置：名称
    authPass      string            // HTTP基本权限设置：密码
    browserMode   bool              // 是否模拟浏览器模式(自动保存提交COOKIE)
    retryCount    int               // 失败重试次数(网络失败情况下)
    retryInterval int               // 失败重试间隔
}

// http客户端对象指针
func NewClient() *Client {
    return &Client{
        Client : http.Client {
            Transport: &http.Transport {
            	// 默认不校验HTTPS证书有效性
	            TLSClientConfig : &tls.Config{
	            	InsecureSkipVerify: true,
				},
	            // 默认关闭KeepAlive功能
                DisableKeepAlives: true,
            },
        },
        header : make(map[string]string),
        cookies: make(map[string]string),
    }
}

// 克隆当前客户端对象，复制属性。
func (c *Client) Clone() *Client {
    newClient := NewClient()
    *newClient = *c
    newClient.header  = make(map[string]string)
    newClient.cookies = make(map[string]string)
    for k, v := range c.header {
        newClient.header[k] = v
    }
    for k, v := range c.cookies {
        newClient.cookies[k] = v
    }
    return newClient
}

// GET请求
func (c *Client) Get(url string) (*ClientResponse, error) {
    return c.DoRequest("GET", url)
}

// PUT请求
func (c *Client) Put(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("PUT", url, data...)
}

// POST请求提交数据，默认使用表单方式提交数据(绝大部分场景下也是如此)。
// 如果服务端对Content-Type有要求，可使用Client对象进行请求，单独设置相关属性。
// 支持文件上传，需要字段格式为：FieldName=@file:
func (c *Client) Post(url string, data...interface{}) (*ClientResponse, error) {
    if len(c.prefix) > 0 {
        url = c.prefix + url
    }
    param := ""
    if len(data) > 0 {
        param = BuildParams(data[0])
    }
    req := (*http.Request)(nil)
    if strings.Contains(param, "@file:") {
        // 文件上传
        buffer := new(bytes.Buffer)
        writer := multipart.NewWriter(buffer)
        for _, item := range strings.Split(param, "&") {
            array := strings.Split(item, "=")
            if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
                path := array[1][6:]
                if !gfile.Exists(path) {
                    return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
                }
                if file, err := writer.CreateFormFile(array[0], path); err == nil {
                    if f, err := os.Open(path); err == nil {
                        defer f.Close()
                        if _, err = io.Copy(file, f); err != nil {
                            return nil, err
                        }
                    } else {
                        return nil, err
                    }
                } else {
                    return nil, err
                }
            } else {
                writer.WriteField(array[0], array[1])
            }
        }
        writer.Close()
        if r, err := http.NewRequest("POST", url, buffer); err != nil {
            return nil, err
        } else {
            req = r
            req.Header.Set("Content-Type", writer.FormDataContentType())
        }
    } else {
        // 识别提交数据格式
        paramBytes := []byte(param)
        if r, err := http.NewRequest("POST", url, bytes.NewReader(paramBytes)); err != nil {
            return nil, err
        } else {
            req = r
            if json.Valid(paramBytes) {
                req.Header.Set("Content-Type", "application/json")
            } else {
                req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
            }
        }
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // COOKIE
    if len(c.cookies) > 0 {
        headerCookie := ""
        for k, v := range c.cookies {
            if len(headerCookie) > 0 {
                headerCookie += ";"
            }
            headerCookie += k + "=" + v
        }
        if len(headerCookie) > 0 {
            req.Header.Set("Cookie", headerCookie)
        }
    }
    // HTTP账号密码
    if len(c.authUser) > 0 {
        req.SetBasicAuth(c.authUser, c.authPass)
    }
    // 执行请求
    resp := (*http.Response)(nil)
    for {
        if r, err := c.Do(req); err != nil {
            if c.retryCount > 0 {
                c.retryCount--
            } else {
                return nil, err
            }
        } else {
            resp = r
            break
        }
    }
    r := &ClientResponse{
        cookies : make(map[string]string),
    }
    r.Response = resp
	// 浏览器模式
	if c.browserMode {
		now := time.Now()
		for _, v := range r.Cookies() {
			if v.Expires.UnixNano() < now.UnixNano() {
				delete(c.cookies, v.Name)
			} else {
				c.cookies[v.Name] = v.Value
			}
		}
	}
    return r, nil
}

// DELETE请求
func (c *Client) Delete(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("DELETE", url, data...)
}

func (c *Client) Head(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("HEAD", url, data...)
}

func (c *Client) Patch(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("PATCH", url, data...)
}

func (c *Client) Connect(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("CONNECT", url, data...)
}

func (c *Client) Options(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("OPTIONS", url, data...)
}

func (c *Client) Trace(url string, data...interface{}) (*ClientResponse, error) {
    return c.DoRequest("TRACE", url, data...)
}

// GET请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) GetContent(url string, data...interface{}) string {
    return c.DoRequestContent("GET", url, data...)
}

// PUT请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) PutContent(url string, data...interface{}) string {
    return c.DoRequestContent("PUT", url, data...)
}

// POST请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) PostContent(url string, data...interface{}) string {
    return c.DoRequestContent("POST", url, data...)
}

// DELETE请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) DeleteContent(url string, data...interface{}) string {
    return c.DoRequestContent("DELETE", url, data...)
}

func (c *Client) HeadContent(url string, data...interface{}) string {
    return c.DoRequestContent("HEAD", url, data...)
}

func (c *Client) PatchContent(url string, data...interface{}) string {
    return c.DoRequestContent("PATCH", url, data...)
}

func (c *Client) ConnectContent(url string, data...interface{}) string {
    return c.DoRequestContent("CONNECT", url, data...)
}

func (c *Client) OptionsContent(url string, data...interface{}) string {
    return c.DoRequestContent("OPTIONS", url, data...)
}

func (c *Client) TraceContent(url string, data...interface{}) string {
    return c.DoRequestContent("TRACE", url, data...)
}

// 请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) DoRequestContent(method string, url string, data...interface{}) string {
    response, err := c.DoRequest(method, url, data...)
    if err != nil {
        return ""
    }
    defer response.Close()
    return string(response.ReadAll())
}

// 请求并返回response对象
func (c *Client) DoRequest(method, url string, data...interface{}) (*ClientResponse, error) {
    if strings.EqualFold("POST", method) {
        return c.Post(url, data...)
    }
    if len(c.prefix) > 0 {
        url = c.prefix + url
    }
    param := ""
    if len(data) > 0 {
        param = BuildParams(data[0])
    }
    req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader([]byte(param)))
    if err != nil {
        return nil, err
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // COOKIE
    if len(c.cookies) > 0 {
        headerCookie := ""
        for k, v := range c.cookies {
            if len(headerCookie) > 0 {
                headerCookie += ";"
            }
            headerCookie += k + "=" + v
        }
        if len(headerCookie) > 0 {
            req.Header.Set("Cookie", headerCookie)
        }
    }
    // 执行请求
    resp := (*http.Response)(nil)
    for {
        if r, err := c.Do(req); err != nil {
            if c.retryCount > 0 {
                c.retryCount--
            } else {
                return nil, err
            }
        } else {
            resp = r
            break
        }
    }
    r := &ClientResponse{
        cookies : make(map[string]string),
    }
    r.Response = resp
    // 浏览器模式
    if c.browserMode {
        now := time.Now()
        for _, v := range r.Cookies() {
            if v.Expires.UnixNano() < now.UnixNano() {
                delete(c.cookies, v.Name)
            } else {
                c.cookies[v.Name] = v.Value
            }
        }
    }
    //fmt.Println(url, c.cookies)
    return r, nil
}




