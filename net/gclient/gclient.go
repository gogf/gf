// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gclient provides convenient http client functionalities.
package gclient

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
)

// Client is the HTTP client for HTTP request management.
type Client struct {
	http.Client                         // Underlying HTTP Client.
	header            map[string]string // Custom header map.
	cookies           map[string]string // Custom cookie map.
	prefix            string            // Prefix for request.
	authUser          string            // HTTP basic authentication: user.
	authPass          string            // HTTP basic authentication: pass.
	retryCount        int               // Retry count when request fails.
	retryInterval     time.Duration     // Retry interval when request fails.
	middlewareHandler []HandlerFunc     // Interceptor handlers
	discovery         gsvc.Discovery    // Discovery for service.
	builder           gsel.Builder      // Builder for request balance.
}

const (
	httpProtocolName          = `http`
	httpParamFileHolder       = `@file:`
	httpRegexParamJson        = `^[\w\[\]]+=.+`
	httpRegexHeaderRaw        = `^([\w\-]+):\s*(.+)`
	httpHeaderHost            = `Host`
	httpHeaderCookie          = `Cookie`
	httpHeaderUserAgent       = `User-Agent`
	httpHeaderContentType     = `Content-Type`
	httpHeaderContentTypeJson = `application/json`
	httpHeaderContentTypeXml  = `application/xml`
	httpHeaderContentTypeForm = `application/x-www-form-urlencoded`
)

var (
	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`GClient %s at %s`, gf.VERSION, hostname)
)

// New creates and returns a new HTTP client object.
func New() *Client {
	c := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header:    make(map[string]string),
		cookies:   make(map[string]string),
		builder:   gsel.GetBuilder(),
		discovery: gsvc.GetRegistry(),
	}
	c.header[httpHeaderUserAgent] = defaultClientAgent
	// It enables OpenTelemetry for client in default.
	c.Use(internalMiddlewareTracing, internalMiddlewareDiscovery)
	return c
}

// Clone deeply clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := New()
	*newClient = *c
	if len(c.header) > 0 {
		newClient.header = make(map[string]string)
		for k, v := range c.header {
			newClient.header[k] = v
		}
	}
	if len(c.cookies) > 0 {
		newClient.cookies = make(map[string]string)
		for k, v := range c.cookies {
			newClient.cookies[k] = v
		}
	}
	return newClient
}

// LoadKeyCrt creates and returns a TLS configuration object with given certificate and key files.
func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error) {
	crtPath, err := gfile.Search(crtFile)
	if err != nil {
		return nil, err
	}
	keyPath, err := gfile.Search(keyFile)
	if err != nil {
		return nil, err
	}
	crt, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		err = gerror.Wrapf(err, `tls.LoadX509KeyPair failed for certFile "%s", keyFile "%s"`, crtPath, keyPath)
		return nil, err
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}
