// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"crypto/tls"
	"net/http/httptrace"
	"net/textproto"
)

type clientTracer struct {
	*httptrace.ClientTrace
}

// newClientTracer creates and returns object of httptrace.ClientTrace.
func newClientTracer(baseClientTracer *httptrace.ClientTrace) *httptrace.ClientTrace {
	c := &clientTracer{
		ClientTrace: baseClientTracer,
	}
	return &httptrace.ClientTrace{
		GetConn:              c.GetConn,
		GotConn:              c.GotConn,
		PutIdleConn:          c.PutIdleConn,
		GotFirstResponseByte: c.GotFirstResponseByte,
		Got100Continue:       c.Got100Continue,
		Got1xxResponse:       c.Got1xxResponse,
		DNSStart:             c.DNSStart,
		DNSDone:              c.DNSDone,
		ConnectStart:         c.ConnectStart,
		ConnectDone:          c.ConnectDone,
		TLSHandshakeStart:    c.TLSHandshakeStart,
		TLSHandshakeDone:     c.TLSHandshakeDone,
		WroteHeaderField:     c.WroteHeaderField,
		WroteHeaders:         c.WroteHeaders,
		Wait100Continue:      c.Wait100Continue,
		WroteRequest:         c.WroteRequest,
	}
}

// GetConn is called before a connection is created or
// retrieved from an idle pool. The hostPort is the
// "host:port" of the target or proxy. GetConn is called even
// if there's already an idle cached connection available.
func (ct *clientTracer) GetConn(hostPort string) {
	ct.ClientTrace.GetConn(hostPort)
}

// GotConn is called after a successful connection is
// obtained. There is no hook for failure to obtain a
// connection; instead, use the error from
// Transport.RoundTrip.
func (ct *clientTracer) GotConn(info httptrace.GotConnInfo) {
	ct.ClientTrace.GotConn(info)
}

// PutIdleConn is called when the connection is returned to
// the idle pool. If err is nil, the connection was
// successfully returned to the idle pool. If err is non-nil,
// it describes why not. PutIdleConn is not called if
// connection reuse is disabled via Transport.DisableKeepAlives.
// PutIdleConn is called before the caller's Response.Body.Close
// call returns.
// For HTTP/2, this hook is not currently used.
func (ct *clientTracer) PutIdleConn(err error) {
	ct.ClientTrace.PutIdleConn(err)
}

// GotFirstResponseByte is called when the first byte of the response
// headers is available.
func (ct *clientTracer) GotFirstResponseByte() {
	ct.ClientTrace.GotFirstResponseByte()
}

// Got100Continue is called if the server replies with a "100
// Continue" response.
func (ct *clientTracer) Got100Continue() {
	ct.ClientTrace.Got100Continue()
}

// Got1xxResponse is called for each 1xx informational response header
// returned before the final non-1xx response. Got1xxResponse is called
// for "100 Continue" responses, even if Got100Continue is also defined.
// If it returns an error, the client request is aborted with that error value.
func (ct *clientTracer) Got1xxResponse(code int, header textproto.MIMEHeader) error {
	return ct.ClientTrace.Got1xxResponse(code, header)
}

// DNSStart is called when a DNS lookup begins.
func (ct *clientTracer) DNSStart(info httptrace.DNSStartInfo) {
	ct.ClientTrace.DNSStart(info)
}

// DNSDone is called when a DNS lookup ends.
func (ct *clientTracer) DNSDone(info httptrace.DNSDoneInfo) {
	ct.ClientTrace.DNSDone(info)
}

// ConnectStart is called when a new connection's Dial begins.
// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (ct *clientTracer) ConnectStart(network, addr string) {
	ct.ClientTrace.ConnectStart(network, addr)
}

// ConnectDone is called when a new connection's Dial
// completes. The provided err indicates whether the
// connection completed successfully.
// If net.Dialer.DualStack ("Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (ct *clientTracer) ConnectDone(network, addr string, err error) {
	ct.ClientTrace.ConnectDone(network, addr, err)
}

// TLSHandshakeStart is called when the TLS handshake is started. When
// connecting to an HTTPS site via an HTTP proxy, the handshake happens
// after the CONNECT request is processed by the proxy.
func (ct *clientTracer) TLSHandshakeStart() {
	ct.ClientTrace.TLSHandshakeStart()
}

// TLSHandshakeDone is called after the TLS handshake with either the
// successful handshake's connection state, or a non-nil error on handshake
// failure.
func (ct *clientTracer) TLSHandshakeDone(state tls.ConnectionState, err error) {
	ct.ClientTrace.TLSHandshakeDone(state, err)
}

// WroteHeaderField is called after the Transport has written
// each request header. At the time of this call the values
// might be buffered and not yet written to the network.
func (ct *clientTracer) WroteHeaderField(key string, value []string) {
	ct.ClientTrace.WroteHeaderField(key, value)
}

// WroteHeaders is called after the Transport has written
// all request headers.
func (ct *clientTracer) WroteHeaders() {
	ct.ClientTrace.WroteHeaders()
}

// Wait100Continue is called if the Request specified
// "Expect: 100-continue" and the Transport has written the
// request headers but is waiting for "100 Continue" from the
// server before writing the request body.
func (ct *clientTracer) Wait100Continue() {
	ct.ClientTrace.Wait100Continue()
}

// WroteRequest is called with the result of writing the
// request and any body. It may be called multiple times
// in the case of retried requests.
func (ct *clientTracer) WroteRequest(info httptrace.WroteRequestInfo) {
	ct.ClientTrace.WroteRequest(info)
}
