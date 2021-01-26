// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/internal/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"
)

type clientTracer struct {
	context.Context
	span        trace.Span
	request     *http.Request
	requestBody []byte
	headers     map[string]interface{}
	mtx         sync.Mutex
}

func newClientTrace(ctx context.Context, span trace.Span, request *http.Request) *httptrace.ClientTrace {
	ct := &clientTracer{
		Context: ctx,
		span:    span,
		request: request,
		headers: make(map[string]interface{}),
	}
	if ct.request.ContentLength <= tracingMaxContentLogSize {
		reqBodyContent, _ := ioutil.ReadAll(ct.request.Body)
		ct.requestBody = reqBodyContent
		ct.request.Body = utils.NewReadCloser(reqBodyContent, false)
	}
	return &httptrace.ClientTrace{
		GetConn:              ct.getConn,
		GotConn:              ct.gotConn,
		PutIdleConn:          ct.putIdleConn,
		GotFirstResponseByte: ct.gotFirstResponseByte,
		Got100Continue:       ct.got100Continue,
		Got1xxResponse:       ct.got1xxResponse,
		DNSStart:             ct.dnsStart,
		DNSDone:              ct.dnsDone,
		ConnectStart:         ct.connectStart,
		ConnectDone:          ct.connectDone,
		TLSHandshakeStart:    ct.tlsHandshakeStart,
		TLSHandshakeDone:     ct.tlsHandshakeDone,
		WroteHeaderField:     ct.wroteHeaderField,
		WroteHeaders:         ct.wroteHeaders,
		Wait100Continue:      ct.wait100Continue,
		WroteRequest:         ct.wroteRequest,
	}
}

func (ct *clientTracer) getConn(host string) {

}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	ct.span.SetAttributes(
		label.String("http.connection.remote", info.Conn.RemoteAddr().String()),
		label.String("http.connection.local", info.Conn.LocalAddr().String()),
	)
}

func (ct *clientTracer) putIdleConn(err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

func (ct *clientTracer) dnsStart(info httptrace.DNSStartInfo) {
	ct.span.SetAttributes(
		label.String("http.dns.start", info.Host),
	)
}

func (ct *clientTracer) dnsDone(info httptrace.DNSDoneInfo) {
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
		label.String("http.dns.done", buffer.String()),
	)
}

func (ct *clientTracer) connectStart(network, addr string) {
	ct.span.SetAttributes(
		label.String("http.connect.start", network+"@"+addr),
	)
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	ct.span.SetAttributes(
		label.String("http.connect.done", network+"@"+addr),
	)
}

func (ct *clientTracer) tlsHandshakeStart() {

}

func (ct *clientTracer) tlsHandshakeDone(_ tls.ConnectionState, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if len(v) > 1 {
		ct.headers[k] = v
	} else {
		ct.headers[k] = v[0]
	}
}

func (ct *clientTracer) wroteHeaders() {

}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, info.Err))
	}
	var bodyContent string
	if ct.request.ContentLength <= tracingMaxContentLogSize {
		bodyContent = string(ct.requestBody)
	} else {
		bodyContent = fmt.Sprintf("[Request Body Too Large For Logging, Max: %d bytes]", tracingMaxContentLogSize)
	}
	ct.span.AddEvent("http.request", trace.WithAttributes(
		label.Any(`http.request.headers`, ct.headers),
		label.String(`http.request.body`, bodyContent),
	))
}

func (ct *clientTracer) got100Continue() {

}

func (ct *clientTracer) wait100Continue() {

}

func (ct *clientTracer) gotFirstResponseByte() {

}

func (ct *clientTracer) got1xxResponse(code int, header textproto.MIMEHeader) error {
	return nil
}
