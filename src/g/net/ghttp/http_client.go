package ghttp

import (
    "net/http"
    "strings"
    "time"
    "g/os/glog"
)

// http客户端对象指针
func NewClient() (*Client) {
    return &Client{}
}

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// GET请求
func (c *Client) Get(url string) *ClientResponse {
    return c.Request("GET", url, "")
}

// PUT请求
func (c *Client) Put(url, data string) *ClientResponse {
    return c.Request("PUT", url, data)
}

// POST请求提交数据
func (c *Client) Post(url, data string) *ClientResponse {
    resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
    if err != nil {
        //glog.Println(err)
        return nil
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r
}

// DELETE请求
func (c *Client) Delete(url, data string) *ClientResponse {
    return c.Request("DELETE", url, data)
}

func (c *Client) Head(url, data string) *ClientResponse {
    return c.Request("HEAD", url, data)
}

func (c *Client) Patch(url, data string) *ClientResponse  {
    return c.Request("PATCH", url, data)
}

func (c *Client) Connect(url, data string) *ClientResponse{
    return c.Request("CONNECT", url, data)
}

func (c *Client) Options(url, data string) *ClientResponse{
    return c.Request("OPTIONS", url, data)
}

func (c *Client) Trace(url, data string) *ClientResponse  {
    return c.Request("TRACE", url, data)
}

// 请求并返回response对象
func (c *Client) Request(method, url, data string) *ClientResponse {
    req, err := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(data))
    if err != nil {
        glog.Println("creating request failed: " + err.Error())
        return nil
    }
    resp, err := c.Do(req)
    if err != nil {
        glog.Println("sending request failed: " + err.Error())
        return nil
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r
}


func Get(url string) *ClientResponse {
    return Request("GET", url, "")
}

func Put(url, data string) *ClientResponse {
    return Request("PUT", url, data)
}

func Post(url, data string) *ClientResponse {
    return Request("PUT", url, data)
}

func Delete(url, data string) *ClientResponse {
    return Request("DELETE", url, data)
}

func Head(url, data string) *ClientResponse {
    return Request("HEAD", url, data)
}

func Patch(url, data string) *ClientResponse  {
    return Request("PATCH", url, data)
}

func Connect(url, data string) *ClientResponse{
    return Request("CONNECT", url, data)
}

func Options(url, data string) *ClientResponse{
    return Request("OPTIONS", url, data)
}

func Trace(url, data string) *ClientResponse  {
    return Request("TRACE", url, data)
}

func Request(method, url, data string) *ClientResponse {
    return NewClient().Request(method, url, data)
}
