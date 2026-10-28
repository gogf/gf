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
	noUrlEncode       bool              // No url encoding for request parameters.
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

// New creates and returns a new HTTP client object with a default timeout of 30 seconds.
func New() *Client {
	return NewWithTimeout(30 * time.Second)
}

// NewWithTimeout creates and returns a new HTTP client object with specified timeout.
//
// The transport is cloned from http.DefaultTransport to inherit standard library defaults
// (such as Proxy, HTTP/2 knobs, and future Go defaults), then customized with the project's
// own TLS, keep-alive, and connection pool settings.
func NewWithTimeout(timeout time.Duration) *Client {
	// Clone from http.DefaultTransport to inherit standard library defaults,
	// then override with project-specific settings.
	var transport *http.Transport
	if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = defaultTransport.Clone()
	} else {
		// Fallback to manual construction if DefaultTransport is not *http.Transport
		// (e.g., if the application replaced it with a custom RoundTripper)
		transport = &http.Transport{
			DisableKeepAlives:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
		}
	}
	// No validation for https certification of the server in default.
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = 50
	transport.MaxConnsPerHost = 100
	defaultClient := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	return NewWithHttpClient(&defaultClient)
}

// NewWithHttpClient creates and returns a new Client with given http.Client.
// It panics if client is nil.
func NewWithHttpClient(client *http.Client) *Client {
	if client == nil {
		panic(`gclient: client must not be nil`)
	}
	c := &Client{
		Client:    *client,
		header:    make(map[string]string),
		cookies:   make(map[string]string),
		builder:   gsel.GetBuilder(),
		discovery: nil,
	}
	c.header[httpHeaderUserAgent] = defaultClientAgent
	// It enables OpenTelemetry for client in default.
	c.Use(internalMiddlewareObservability, internalMiddlewareDiscovery)
	return c
}

// Clone deeply clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := New()
	*newClient = *c
	newClient.header = make(map[string]string, len(c.header))
	for k, v := range c.header {
		newClient.header[k] = v
	}
	newClient.cookies = make(map[string]string, len(c.cookies))
	for k, v := range c.cookies {
		newClient.cookies[k] = v
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
