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

type metricManager struct {
	HttpServerRequestActive        gmetric.UpDownCounter
	HttpServerRequestTotal         gmetric.Counter
	HttpServerRequestDuration      gmetric.Histogram
	HttpServerRequestDurationTotal gmetric.Counter
	HttpServerRequestBodySize      gmetric.Counter
	HttpServerResponseBodySize     gmetric.Counter
}

const (
	metricAttrKeyServerName             gmetric.AttributeKey = "server.name"
	metricAttrKeyServerAddress          gmetric.AttributeKey = "server.address"
	metricAttrKeyServerPort             gmetric.AttributeKey = "server.port"
	metricAttrKeyHttpRoute              gmetric.AttributeKey = "http.route"
	metricAttrKeyUrlSchema              gmetric.AttributeKey = "url.schema"
	metricAttrKeyHttpRequestMethod      gmetric.AttributeKey = "http.request.method"
	metricAttrKeyErrorCode              gmetric.AttributeKey = "error.code"
	metricAttrKeyHttpResponseStatusCode gmetric.AttributeKey = "http.response.status_code"
	metricAttrKeyNetworkProtocolVersion gmetric.AttributeKey = "network.protocol.version"
)

func newMetricManager() *metricManager {
	mm := &metricManager{
		HttpServerRequestDuration: gmetric.MustNewHistogram(gmetric.MetricConfig{
			Name:              "http.server.request.duration",
			Help:              "Measures the duration of inbound request.",
			Unit:              "ms",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
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
		}),
		HttpServerRequestTotal: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.server.request.total",
			Help:              "Total processed request number.",
			Unit:              "",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpServerRequestActive: gmetric.MustNewUpDownCounter(gmetric.MetricConfig{
			Name:              "http.server.request.active",
			Help:              "Number of active server requests.",
			Unit:              "",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpServerRequestDurationTotal: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.server.request.duration_total",
			Help:              "Total execution duration of request.",
			Unit:              "ms",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpServerRequestBodySize: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.server.request.body_size",
			Help:              "Incoming request bytes total.",
			Unit:              "bytes",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpServerResponseBodySize: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.server.response.body_size",
			Help:              "Response bytes total.",
			Unit:              "bytes",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
	}
	return mm
}

func (m *metricManager) GetMetricOptionForDurationByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerName,
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
		),
	}
}

func (m *metricManager) GetMetricOptionForActiveByRequest(r *Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return gmetric.Option{
		Attributes: attrMap.PickEx(
			metricAttrKeyErrorCode,
			metricAttrKeyHttpResponseStatusCode,
		),
	}
}

func (m *metricManager) GetMetricOptionForActiveByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.PickEx(
			metricAttrKeyErrorCode,
			metricAttrKeyHttpResponseStatusCode,
		),
	}
}

func (m *metricManager) GetMetricOptionForResponseByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.PickEx(),
	}
}

func (m *metricManager) GetMetricAttributeMap(r *Request) gmetric.AttributeMap {
	var (
		serverAddress   string
		serverPort      string
		httpRoute       string
		protocolVersion string
		handler         = r.GetServeHandler()
		localAddr       = r.Context().Value(http.LocalAddrContextKey)
		attrMap         = make(gmetric.AttributeMap)
	)
	if localAddr != nil {
		addr := localAddr.(net.Addr)
		serverAddress, serverPort = gstr.List2(addr.String(), ":")
	}
	if handler.Handler.Router != nil {
		httpRoute = handler.Handler.Router.Uri
	} else {
		httpRoute = r.URL.Path
	}
	if array := gstr.Split(r.Proto, "/"); len(array) > 1 {
		protocolVersion = array[1]
	}
	attrMap.Sets(gmetric.AttributeMap{
		metricAttrKeyServerName:             r.Server.GetName(),
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
	s.metricManager.HttpServerRequestActive.Inc(
		r.Context(),
		s.metricManager.GetMetricOptionForActiveByRequest(r),
	)
}

func (s *Server) handleMetricsAfterRequestDone(r *Request) {
	if !gmetric.IsEnabled() {
		return
	}
	var (
		mm             = s.metricManager
		ctx            = r.Context()
		attrMap        = mm.GetMetricAttributeMap(r)
		durationMilli  = float64(r.LeaveTime.Sub(r.EnterTime).Milliseconds())
		responseOption = mm.GetMetricOptionForResponseByMap(attrMap)
	)
	mm.HttpServerRequestActive.Dec(ctx, mm.GetMetricOptionForActiveByMap(attrMap))
	mm.HttpServerRequestTotal.Inc(ctx, responseOption)
	mm.HttpServerRequestBodySize.Add(
		ctx,
		float64(r.ContentLength),
		responseOption,
	)
	mm.HttpServerResponseBodySize.Add(
		ctx,
		float64(r.Response.BytesWritten()),
		responseOption,
	)
	mm.HttpServerRequestDurationTotal.Add(
		ctx,
		durationMilli,
		responseOption,
	)
	mm.HttpServerRequestDuration.Record(durationMilli, mm.GetMetricOptionForDurationByMap(attrMap))
}
