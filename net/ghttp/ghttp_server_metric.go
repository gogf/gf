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
	HttpServerRequestActive   gmetric.Counter
	HttpServerRequestTotal    gmetric.Counter
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
		HttpServerRequestActive: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.server.request.active",
			Help:              "Number of active server requests.",
			Unit:              "",
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

func (m *metricManager) GetMetricOptionForTotalByMap(attrMap gmetric.AttributeMap) gmetric.Option {
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
