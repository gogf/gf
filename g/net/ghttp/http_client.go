package ghttp

import (
    "net/http"
    "strings"
    "time"
    "net/url"
    "bytes"
)

// http客户端
type Client struct {
    http.Client
}

// 请求对象
type ClientRequest struct {
    http.Request
    getvals *url.Values // GET参数
}

// 客户端请求结果对象
type ClientResponse struct {
    http.Response
}

// http客户端对象指针
func NewClient() (*Client) {
    return &Client{}
}

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// GET请求
func (c *Client) Get(url string) (*ClientResponse, error) {
    return c.Request("GET", url, []byte(""))
}

// PUT请求
func (c *Client) Put(url, data string) (*ClientResponse, error) {
    return c.Request("PUT", url, []byte(data))
}

// POST请求提交数据
func (c *Client) Post(url, data string) (*ClientResponse, error) {
    resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r, nil
}

// DELETE请求
func (c *Client) Delete(url, data string) (*ClientResponse, error) {
    return c.Request("DELETE", url, []byte(data))
}

func (c *Client) Head(url, data string) (*ClientResponse, error) {
    return c.Request("HEAD", url, []byte(data))
}

func (c *Client) Patch(url, data string) (*ClientResponse, error) {
    return c.Request("PATCH", url, []byte(data))
}

func (c *Client) Connect(url, data string) (*ClientResponse, error) {
    return c.Request("CONNECT", url, []byte(data))
}

func (c *Client) Options(url, data string) (*ClientResponse, error) {
    return c.Request("OPTIONS", url, []byte(data))
}

func (c *Client) Trace(url, data string) (*ClientResponse, error) {
    return c.Request("TRACE", url, []byte(data))
}

// 请求并返回response对象，该方法支持二进制提交数据
func (c *Client) Request(method, url string, data []byte) (*ClientResponse, error) {
    req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r, nil
}


func Get(url string) (*ClientResponse, error) {
    return Request("GET", url, []byte(""))
}

func Put(url, data string) (*ClientResponse, error) {
    return Request("PUT", url, []byte(data))
}

func Post(url, data string) (*ClientResponse, error) {
    return Request("PUT", url, []byte(data))
}

func Delete(url, data string) (*ClientResponse, error) {
    return Request("DELETE", url, []byte(data))
}

func Head(url, data string) (*ClientResponse, error) {
    return Request("HEAD", url, []byte(data))
}

func Patch(url, data string) (*ClientResponse, error) {
    return Request("PATCH", url, []byte(data))
}

func Connect(url, data string) (*ClientResponse, error) {
    return Request("CONNECT", url, []byte(data))
}

func Options(url, data string) (*ClientResponse, error) {
    return Request("OPTIONS", url, []byte(data))
}

func Trace(url, data string) (*ClientResponse, error) {
    return Request("TRACE", url, []byte(data))
}

// 该方法支持二进制提交数据
func Request(method, url string, data []byte) (*ClientResponse, error) {
    return NewClient().Request(method, url, data)
}
