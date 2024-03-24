// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
)

// clientTracerTracing is used for implementing httptrace.ClientTrace.
type clientTracerTracing struct {
	context.Context
	span        trace.Span
	request     *http.Request
	requestBody []byte
	headers     map[string]interface{}
	mtx         sync.Mutex
}

// newClientTracerTracing creates and returns object of httptrace.ClientTrace.
func newClientTracerTracing(
	ctx context.Context,
	span trace.Span,
	request *http.Request,
) *httptrace.ClientTrace {
	ct := &clientTracerTracing{
		Context: ctx,
		span:    span,
		request: request,
		headers: make(map[string]interface{}),
	}

	reqBodyContent, _ := io.ReadAll(ct.request.Body)
	ct.requestBody = reqBodyContent
	ct.request.Body = utils.NewReadCloser(reqBodyContent, false)

	return &httptrace.ClientTrace{
		GetConn:              ct.GetConn,
		GotConn:              ct.GotConn,
		PutIdleConn:          ct.PutIdleConn,
		GotFirstResponseByte: ct.GotFirstResponseByte,
		Got100Continue:       ct.Got100Continue,
		Got1xxResponse:       ct.Got1xxResponse,
		DNSStart:             ct.DNSStart,
		DNSDone:              ct.DNSDone,
		ConnectStart:         ct.ConnectStart,
		ConnectDone:          ct.ConnectDone,
		TLSHandshakeStart:    ct.TLSHandshakeStart,
		TLSHandshakeDone:     ct.TLSHandshakeDone,
		WroteHeaderField:     ct.WroteHeaderField,
		WroteHeaders:         ct.WroteHeaders,
		Wait100Continue:      ct.Wait100Continue,
		WroteRequest:         ct.WroteRequest,
	}
}

// GetConn is called before a connection is created or
// retrieved from an idle pool. The hostPort is the
// "host:port" of the target or proxy. GetConn is called even
// if there's already an idle cached connection available.
func (ct *clientTracerTracing) GetConn(host string) {}

// GotConn is called after a successful connection is
// obtained. There is no hook for failure to obtain a
// connection; instead, use the error from
// Transport.RoundTrip.
func (ct *clientTracerTracing) GotConn(info httptrace.GotConnInfo) {
	remoteAddr := ""
	if info.Conn.RemoteAddr() != nil {
		remoteAddr = info.Conn.RemoteAddr().String()
	}
	localAddr := ""
	if info.Conn.LocalAddr() != nil {
		localAddr = info.Conn.LocalAddr().String()
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpAddressRemote, remoteAddr),
		attribute.String(tracingAttrHttpAddressLocal, localAddr),
	)
}

// PutIdleConn is called when the connection is returned to
// the idle pool. If err is nil, the connection was
// successfully returned to the idle pool. If err is non-nil,
// it describes why not. PutIdleConn is not called if
// connection reuse is disabled via Transport.DisableKeepAlives.
// PutIdleConn is called before the caller's Response.Body.Close
// call returns.
// For HTTP/2, this hook is not currently used.
func (ct *clientTracerTracing) PutIdleConn(err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

// GotFirstResponseByte is called when the first byte of the response
// headers is available.
func (ct *clientTracerTracing) GotFirstResponseByte() {}

// Got100Continue is called if the server replies with a "100
// Continue" response.
func (ct *clientTracerTracing) Got100Continue() {}

// Got1xxResponse is called for each 1xx informational response header
// returned before the final non-1xx response. Got1xxResponse is called
// for "100 Continue" responses, even if Got100Continue is also defined.
// If it returns an error, the client request is aborted with that error value.
func (ct *clientTracerTracing) Got1xxResponse(code int, header textproto.MIMEHeader) error {
	return nil
}

// DNSStart is called when a DNS lookup begins.
func (ct *clientTracerTracing) DNSStart(info httptrace.DNSStartInfo) {
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpDnsStart, info.Host),
	)
}

// DNSDone is called when a DNS lookup ends.
func (ct *clientTracerTracing) DNSDone(info httptrace.DNSDoneInfo) {
	var buffer strings.Builder
	for _, v := range info.Addrs {
		if buffer.Len() != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v.String())
	}
	if info.Err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, info.Err))
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpDnsDone, buffer.String()),
	)
}

// ConnectStart is called when a new connection's Dial begins.
// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (ct *clientTracerTracing) ConnectStart(network, addr string) {
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpConnectStart, network+"@"+addr),
	)
}

// ConnectDone is called when a new connection's Dial
// completes. The provided err indicates whether the
// connection completed successfully.
// If net.Dialer.DualStack ("Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (ct *clientTracerTracing) ConnectDone(network, addr string, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpConnectDone, network+"@"+addr),
	)
}

// TLSHandshakeStart is called when the TLS handshake is started. When
// connecting to an HTTPS site via an HTTP proxy, the handshake happens
// after the CONNECT request is processed by the proxy.
func (ct *clientTracerTracing) TLSHandshakeStart() {}

// TLSHandshakeDone is called after the TLS handshake with either the
// successful handshake's connection state, or a non-nil error on handshake
// failure.
func (ct *clientTracerTracing) TLSHandshakeDone(_ tls.ConnectionState, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

// WroteHeaderField is called after the Transport has written
// each request header. At the time of this call the values
// might be buffered and not yet written to the network.
func (ct *clientTracerTracing) WroteHeaderField(k string, v []string) {
	if len(v) > 1 {
		ct.headers[k] = v
	} else if len(v) == 1 {
		ct.headers[k] = v[0]
	}
}

// WroteHeaders is called after the Transport has written
// all request headers.
func (ct *clientTracerTracing) WroteHeaders() {}

// Wait100Continue is called if the Request specified
// "Expect: 100-continue" and the Transport has written the
// request headers but is waiting for "100 Continue" from the
// server before writing the request body.
func (ct *clientTracerTracing) Wait100Continue() {}

// WroteRequest is called with the result of writing the
// request and any body. It may be called multiple times
// in the case of retried requests.
func (ct *clientTracerTracing) WroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, info.Err))
	}

	reqBodyContent, err := gtrace.SafeContentForHttp(ct.requestBody, ct.request.Header)
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`converting safe content failed: %s`, err.Error()))
	}

	ct.span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
		attribute.String(tracingEventHttpRequestHeaders, gconv.String(ct.headers)),
		attribute.String(tracingEventHttpRequestBaggage, gtrace.GetBaggageMap(ct.Context).String()),
		attribute.String(tracingEventHttpRequestBody, reqBodyContent),
	))
}
