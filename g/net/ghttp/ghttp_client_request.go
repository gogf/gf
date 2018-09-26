// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// HTTP客户端请求.

package ghttp

import (
    "time"
    "bytes"
    "strings"
    "net/http"
    "mime/multipart"
    "os"
    "io"
    "gitee.com/johng/gf/g/os/gfile"
    "errors"
    "fmt"
)

// http客户端
type Client struct {
    http.Client                // 底层http client对象
    header   map[string]string // HEADER信息Map
    authUser string            // HTTP基本权限设置：名称
    authPass string            // HTTP基本权限设置：密码
}

// http客户端对象指针
func NewClient() (*Client) {
    return &Client{
        Client : http.Client {
            Transport: &http.Transport {
                DisableKeepAlives: true,
            },
        },
        header : make(map[string]string),
    }
}

// 设置HTTP Headerss
func (c *Client) SetHeader(key, value string) {
    c.header[key] = value
}

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// 设置HTTP访问账号密码
func (c *Client) SetBasicAuth(user, pass string) {
    c.authUser = user
    c.authPass = pass
}

// GET请求
func (c *Client) Get(url string) (*ClientResponse, error) {
    return c.DoRequest("GET", url, []byte(""))
}

// PUT请求
func (c *Client) Put(url, data string) (*ClientResponse, error) {
    return c.DoRequest("PUT", url, []byte(data))
}

// POST请求提交数据，默认使用表单方式提交数据(绝大部分场景下也是如此)。
// 如果服务端对Content-Type有要求，可使用Client对象进行请求，单独设置相关属性。
// 支持文件上传，需要字段格式为：FieldName=@file:
func (c *Client) Post(url, data string) (*ClientResponse, error) {
    var req *http.Request
    if strings.Contains(data, "@file:") {
        buffer := new(bytes.Buffer)
        writer := multipart.NewWriter(buffer)
        for _, item := range strings.Split(data, "&") {
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
        if r, err := http.NewRequest("POST", url, bytes.NewReader([]byte(data))); err != nil {
            return nil, err
        } else {
            req = r
            req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        }
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // HTTP账号密码
    if len(c.authUser) > 0 {
        req.SetBasicAuth(c.authUser, c.authPass)
    }
    // 执行请求
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r, nil
}

// DELETE请求
func (c *Client) Delete(url, data string) (*ClientResponse, error) {
    return c.DoRequest("DELETE", url, []byte(data))
}

func (c *Client) Head(url, data string) (*ClientResponse, error) {
    return c.DoRequest("HEAD", url, []byte(data))
}

func (c *Client) Patch(url, data string) (*ClientResponse, error) {
    return c.DoRequest("PATCH", url, []byte(data))
}

func (c *Client) Connect(url, data string) (*ClientResponse, error) {
    return c.DoRequest("CONNECT", url, []byte(data))
}

func (c *Client) Options(url, data string) (*ClientResponse, error) {
    return c.DoRequest("OPTIONS", url, []byte(data))
}

func (c *Client) Trace(url, data string) (*ClientResponse, error) {
    return c.DoRequest("TRACE", url, []byte(data))
}

// 请求并返回response对象，该方法支持二进制提交数据
func (c *Client) DoRequest(method, url string, data []byte) (*ClientResponse, error) {
    if strings.Compare("POST", strings.ToUpper(method)) == 0 {
        return c.Post(url, string(data))
    }
    req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // 执行请求
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r, nil
}


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
