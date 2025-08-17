// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net"
	"net/http"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/text/gstr"
)

type localMetricManager struct {
	HttpServerRequestActive        gmetric.UpDownCounter
	HttpServerRequestTotal         gmetric.Counter
	HttpServerRequestDuration      gmetric.Histogram
	HttpServerRequestDurationTotal gmetric.Counter
	HttpServerRequestBodySize      gmetric.Counter
	HttpServerResponseBodySize     gmetric.Counter
}

const (
	metricAttrKeyServerAddress          = "server.address"
	metricAttrKeyServerPort             = "server.port"
	metricAttrKeyHttpRoute              = "http.route"
	metricAttrKeyUrlSchema              = "url.schema"
	metricAttrKeyHttpRequestMethod      = "http.request.method"
	metricAttrKeyErrorCode              = "error.code"
	metricAttrKeyHttpResponseStatusCode = "http.response.status_code"
	metricAttrKeyNetworkProtocolVersion = "network.protocol.version"
)

var (
	// metricManager for http server metrics.
	metricManager = newMetricManager()
)

func newMetricManager() *localMetricManager {
	meter := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        instrumentName,
		InstrumentVersion: gf.VERSION,
	})
	mm := &localMetricManager{
		HttpServerRequestDuration: meter.MustHistogram(
			"http.server.request.duration",
			gmetric.MetricOption{
				Help:       "Measures the duration of inbound request.",
				Unit:       "ms",
				Attributes: gmetric.Attributes{},
				Buckets: []float64{
					1,
					5,
					10,
					25,
					50,
					75,
					100,
					250,
					500,
					750,
					1000,
					2500,
					5000,
					7500,
					10000,
					30000,
					60000,
				},
			},
		),
		HttpServerRequestTotal: meter.MustCounter(
			"http.server.request.total",
			gmetric.MetricOption{
				Help:       "Total processed request number.",
				Unit:       "",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpServerRequestActive: meter.MustUpDownCounter(
			"http.server.request.active",
			gmetric.MetricOption{
				Help:       "Number of active server requests.",
				Unit:       "",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpServerRequestDurationTotal: meter.MustCounter(
			"http.server.request.duration_total",
			gmetric.MetricOption{
				Help:       "Total execution duration of request.",
				Unit:       "ms",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpServerRequestBodySize: meter.MustCounter(
			"http.server.request.body_size",
			gmetric.MetricOption{
				Help:       "Incoming request bytes total.",
				Unit:       "bytes",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpServerResponseBodySize: meter.MustCounter(
			"http.server.response.body_size",
			gmetric.MetricOption{
				Help:       "Response bytes total.",
				Unit:       "bytes",
				Attributes: gmetric.Attributes{},
			},
		),
	}
	return mm
}

func (m *localMetricManager) GetMetricOptionForRequestDurationByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForRequest(r *Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return m.GetMetricOptionForRequestByMap(attrMap)
}

func (m *localMetricManager) GetMetricOptionForRequestByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyHttpRoute,
			metricAttrKeyUrlSchema,
			metricAttrKeyHttpRequestMethod,
			metricAttrKeyNetworkProtocolVersion,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForResponseByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyHttpRoute,
			metricAttrKeyUrlSchema,
			metricAttrKeyHttpRequestMethod,
			metricAttrKeyNetworkProtocolVersion,
			metricAttrKeyErrorCode,
			metricAttrKeyHttpResponseStatusCode,
		),
	}
}

func (m *localMetricManager) GetMetricAttributeMap(r *Request) gmetric.AttributeMap {
	var (
		serverAddress   string
		serverPort      string
		httpRoute       string
		protocolVersion string
		handler         = r.GetServeHandler()
		localAddr       = r.Context().Value(http.LocalAddrContextKey)
		attrMap         = make(gmetric.AttributeMap)
	)
	serverAddress, serverPort = gstr.List2(r.Host, ":")
	if localAddr != nil {
		_, serverPort = gstr.List2(localAddr.(net.Addr).String(), ":")
	}
	if handler != nil && handler.Handler.Router != nil {
		httpRoute = handler.Handler.Router.Uri
	} else {
		httpRoute = r.URL.Path
	}
	if array := gstr.Split(r.Proto, "/"); len(array) > 1 {
		protocolVersion = array[1]
	}
	attrMap.Sets(gmetric.AttributeMap{
		metricAttrKeyServerAddress:          serverAddress,
		metricAttrKeyServerPort:             serverPort,
		metricAttrKeyHttpRoute:              httpRoute,
		metricAttrKeyUrlSchema:              r.GetSchema(),
		metricAttrKeyHttpRequestMethod:      r.Method,
		metricAttrKeyNetworkProtocolVersion: protocolVersion,
	})
	if r.LeaveTime != nil {
		var errCode int
		if err := r.GetError(); err != nil {
			errCode = gerror.Code(err).Code()
		}
		attrMap.Sets(gmetric.AttributeMap{
			metricAttrKeyErrorCode:              errCode,
			metricAttrKeyHttpResponseStatusCode: r.Response.Status,
		})
	}
	return attrMap
}

func (s *Server) handleMetricsBeforeRequest(r *Request) {
	if !gmetric.IsEnabled() {
		return
	}
	var (
		ctx           = r.Context()
		attrMap       = metricManager.GetMetricAttributeMap(r)
		requestOption = metricManager.GetMetricOptionForRequestByMap(attrMap)
	)
	metricManager.HttpServerRequestActive.Inc(
		ctx,
		requestOption,
	)
	metricManager.HttpServerRequestBodySize.Add(
		ctx,
		float64(r.ContentLength),
		requestOption,
	)
}

func (s *Server) handleMetricsAfterRequestDone(r *Request) {
	if !gmetric.IsEnabled() {
		return
	}
	var (
		ctx             = r.Context()
		attrMap         = metricManager.GetMetricAttributeMap(r)
		durationMilli   = float64(r.LeaveTime.Sub(r.EnterTime).Milliseconds())
		responseOption  = metricManager.GetMetricOptionForResponseByMap(attrMap)
		histogramOption = metricManager.GetMetricOptionForRequestDurationByMap(attrMap)
	)
	metricManager.HttpServerRequestTotal.Inc(ctx, responseOption)
	metricManager.HttpServerRequestActive.Dec(
		ctx,
		metricManager.GetMetricOptionForRequestByMap(attrMap),
	)
	metricManager.HttpServerResponseBodySize.Add(
		ctx,
		float64(r.Response.BytesWritten()),
		responseOption,
	)
	metricManager.HttpServerRequestDurationTotal.Add(
		ctx,
		durationMilli,
		responseOption,
	)
	metricManager.HttpServerRequestDuration.Record(
		durationMilli,
		histogramOption,
	)
}
