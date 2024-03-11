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
	HttpServerRequestDuration gmetric.Histogram
	HttpServerRequestTotal    gmetric.Counter
	HttpServerRequestActive   gmetric.Gauge
}

const (
	metricAttrKeyNetServiceName         gmetric.AttributeKey = "net.service.name"
	metricAttrKeyNetHostAddress         gmetric.AttributeKey = "net.host.address"
	metricAttrKeyNetHostPort            gmetric.AttributeKey = "net.host.port"
	metricAttrKeyHttpRequestRoute       gmetric.AttributeKey = "http.request.route"
	metricAttrKeyHttpRequestSchema      gmetric.AttributeKey = "http.request.schema"
	metricAttrKeyHttpRequestVersion     gmetric.AttributeKey = "http.request.version"
	metricAttrKeyHttpRequestMethod      gmetric.AttributeKey = "http.request.method"
	metricAttrKeyHttpResponseErrorCode  gmetric.AttributeKey = "http.response.error_code"
	metricAttrKeyHttpResponseStatusCode gmetric.AttributeKey = "http.response.status_code"
)

func newMetricManager() *metricManager {
	mm := &metricManager{
		HttpServerRequestDuration: gmetric.MustNewHistogram(gmetric.HistogramConfig{
			MetricConfig: gmetric.MetricConfig{
				Name:              "http.server.request.duration",
				Help:              "Measures the duration of inbound request.",
				Unit:              "ms",
				Attributes:        gmetric.Attributes{},
				Instrument:        instrumentName,
				InstrumentVersion: gf.VERSION,
			},
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
			},
		}),
		HttpServerRequestTotal: gmetric.MustNewCounter(gmetric.CounterConfig{
			MetricConfig: gmetric.MetricConfig{
				Name:              "http.server.request.total",
				Help:              "Total processed request number.",
				Unit:              "",
				Attributes:        gmetric.Attributes{},
				Instrument:        instrumentName,
				InstrumentVersion: gf.VERSION,
			},
		}),
		HttpServerRequestActive: gmetric.MustNewGauge(gmetric.GaugeConfig{
			MetricConfig: gmetric.MetricConfig{
				Name:              "http.server.request.active",
				Help:              "Number of active server requests.",
				Unit:              "",
				Attributes:        gmetric.Attributes{},
				Instrument:        instrumentName,
				InstrumentVersion: gf.VERSION,
			},
		}),
	}
	return mm
}

func (m *metricManager) getMetricOptionForDurationByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyNetServiceName,
			metricAttrKeyNetHostAddress,
			metricAttrKeyNetHostPort,
		),
	}
}

func (m *metricManager) GetMetricOptionForActiveByRequest(r *Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return gmetric.Option{
		Attributes: attrMap.PickEx(
			metricAttrKeyHttpResponseErrorCode,
			metricAttrKeyHttpResponseStatusCode,
		),
	}
}

func (m *metricManager) GetMetricOptionForActiveByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.PickEx(
			metricAttrKeyHttpResponseErrorCode,
			metricAttrKeyHttpResponseStatusCode,
		),
	}
}

func (m *metricManager) GetMetricOptionForTotalByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.PickEx(),
	}
}

func (m *metricManager) GetMetricAttributeMap(r *Request) gmetric.AttributeMap {
	var (
		hostAddress string
		hostPort    string
		reqRoute    string
		reqVersion  string
		localAddr   = r.Context().Value(http.LocalAddrContextKey)
		attrMap     = make(gmetric.AttributeMap)
	)
	if localAddr != nil {
		addr := localAddr.(net.Addr)
		hostAddress, hostPort = gstr.List2(addr.String(), ":")
	}
	if r.Router != nil {
		reqRoute = r.Router.Uri
	}
	if array := gstr.Split(r.Proto, "/"); len(array) > 1 {
		reqVersion = array[1]
	}
	attrMap.Sets(gmetric.AttributeMap{
		metricAttrKeyNetServiceName:     r.Server.GetName(),
		metricAttrKeyNetHostAddress:     hostAddress,
		metricAttrKeyNetHostPort:        hostPort,
		metricAttrKeyHttpRequestRoute:   reqRoute,
		metricAttrKeyHttpRequestSchema:  r.GetSchema(),
		metricAttrKeyHttpRequestVersion: reqVersion,
		metricAttrKeyHttpRequestMethod:  r.Method,
	})
	if r.LeaveTime != 0 {
		var errCode int
		if err := r.GetError(); err != nil {
			errCode = gerror.Code(err).Code()
		}
		attrMap.Sets(gmetric.AttributeMap{
			metricAttrKeyHttpResponseErrorCode:  errCode,
			metricAttrKeyHttpResponseStatusCode: r.Response.Status,
		})
	}
	return attrMap
}
