// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"net/http"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

type localMetricManager struct {
	HttpClientOpenConnections      gmetric.UpDownCounter
	HttpClientRequestActive        gmetric.UpDownCounter
	HttpClientRequestTotal         gmetric.Counter
	HttpClientRequestDuration      gmetric.Histogram
	HttpClientRequestDurationTotal gmetric.Counter
	HttpClientConnectionDuration   gmetric.Histogram
	HttpClientRequestBodySize      gmetric.Counter
	HttpClientResponseBodySize     gmetric.Counter
}

const (
	metricAttrKeyServerAddress          gmetric.AttributeKey = "server.address"
	metricAttrKeyServerPort             gmetric.AttributeKey = "server.port"
	metricAttrKeyUrlSchema              gmetric.AttributeKey = "url.schema"
	metricAttrKeyHttpRequestMethod      gmetric.AttributeKey = "http.request.method"
	metricAttrKeyHttpResponseStatusCode gmetric.AttributeKey = "http.response.status_code"
	metricAttrKeyNetworkProtocolVersion gmetric.AttributeKey = "network.protocol.version"
	metricAttrKeHttpConnectionState     gmetric.AttributeKey = "http.connection.state"
)

const (
	connectionStateActive = "active"
	connectionStateIdle   = "idle"
)

var (
	// metricManager for http client metrics.
	metricManager = newMetricManager()
)

func newMetricManager() *localMetricManager {
	durationBuckets := []float64{
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
	}
	mm := &localMetricManager{
		HttpClientRequestDuration: gmetric.MustNewHistogram(gmetric.MetricConfig{
			Name:              "http.client.request.duration",
			Help:              "Measures the duration of client requests.",
			Unit:              "ms",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
			Buckets:           durationBuckets,
		}),
		HttpClientRequestTotal: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.client.request.total",
			Help:              "Total processed request number.",
			Unit:              "",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientRequestActive: gmetric.MustNewUpDownCounter(gmetric.MetricConfig{
			Name:              "http.client.request.active",
			Help:              "Number of active client requests.",
			Unit:              "",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientRequestDurationTotal: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.client.request.duration_total",
			Help:              "Total execution duration of request.",
			Unit:              "ms",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientRequestBodySize: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.client.request.body_size",
			Help:              "Outgoing request bytes total.",
			Unit:              "bytes",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientResponseBodySize: gmetric.MustNewCounter(gmetric.MetricConfig{
			Name:              "http.client.response.body_size",
			Help:              "Response bytes total.",
			Unit:              "bytes",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientOpenConnections: gmetric.MustNewUpDownCounter(gmetric.MetricConfig{
			Name:              "http.client.open_connections",
			Help:              "Number of outbound HTTP connections that are currently active or idle on the client.",
			Unit:              "",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
		}),
		HttpClientConnectionDuration: gmetric.MustNewHistogram(gmetric.MetricConfig{
			Name:              "http.client.connection_duration",
			Help:              "Measures the connection establish duration of client requests.",
			Unit:              "ms",
			Attributes:        gmetric.Attributes{},
			Instrument:        instrumentName,
			InstrumentVersion: gf.VERSION,
			Buckets:           durationBuckets,
		}),
	}
	return mm
}

func (m *localMetricManager) GetMetricOptionForOpenConnectionsByMap(
	state string, attrMap gmetric.AttributeMap,
) gmetric.Option {
	attrMap[metricAttrKeHttpConnectionState] = state
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyUrlSchema,
			metricAttrKeyNetworkProtocolVersion,
			metricAttrKeHttpConnectionState,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForConnectionDuration(r *http.Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyUrlSchema,
			metricAttrKeyNetworkProtocolVersion,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForRequestDurationByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForRequest(r *http.Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return m.GetMetricOptionForRequestByMap(attrMap)
}

func (m *localMetricManager) GetMetricOptionForRequestByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyHttpRequestMethod,
			metricAttrKeyUrlSchema,
			metricAttrKeyNetworkProtocolVersion,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForResponseByMap(attrMap gmetric.AttributeMap) gmetric.Option {
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
			metricAttrKeyHttpRequestMethod,
			metricAttrKeyHttpResponseStatusCode,
			metricAttrKeyUrlSchema,
			metricAttrKeyNetworkProtocolVersion,
		),
	}
}

func (m *localMetricManager) GetMetricAttributeMap(r *http.Request) gmetric.AttributeMap {
	var (
		serverAddress   string
		serverPort      string
		protocolVersion string
		attrMap         = make(gmetric.AttributeMap)
	)
	serverAddress, serverPort = gstr.List2(r.Host, ":")
	if serverPort == "" {
		_, serverPort = gstr.List2(r.RemoteAddr, ":")
	}
	if serverPort == "" {
		serverPort = "80"
		if r.URL.Scheme == "https" {
			serverPort = "443"
		}
	}
	if array := gstr.Split(r.Proto, "/"); len(array) > 1 {
		protocolVersion = array[1]
	}
	attrMap.Sets(gmetric.AttributeMap{
		metricAttrKeyServerAddress:          serverAddress,
		metricAttrKeyServerPort:             serverPort,
		metricAttrKeyUrlSchema:              r.URL.Scheme,
		metricAttrKeyHttpRequestMethod:      r.Method,
		metricAttrKeyNetworkProtocolVersion: protocolVersion,
	})
	if r.Response != nil {
		attrMap.Sets(gmetric.AttributeMap{
			metricAttrKeyHttpResponseStatusCode: r.Response.Status,
		})
	}
	return attrMap
}

func (c *Client) handleMetricsBeforeRequest(r *http.Request) {
	if !gmetric.IsEnabled() {
		return
	}

	var (
		ctx           = r.Context()
		attrMap       = metricManager.GetMetricAttributeMap(r)
		requestOption = metricManager.GetMetricOptionForRequestByMap(attrMap)
	)
	metricManager.HttpClientRequestActive.Inc(
		ctx,
		requestOption,
	)
	metricManager.HttpClientRequestBodySize.Add(
		ctx,
		float64(r.ContentLength),
		requestOption,
	)
}

func (c *Client) handleMetricsAfterRequestDone(r *http.Request, requestStartTime *gtime.Time) {
	if !gmetric.IsEnabled() {
		return
	}

	var (
		ctx             = r.Context()
		attrMap         = metricManager.GetMetricAttributeMap(r)
		duration        = float64(gtime.Now().Sub(requestStartTime).Milliseconds())
		requestOption   = metricManager.GetMetricOptionForRequestByMap(attrMap)
		responseOption  = metricManager.GetMetricOptionForResponseByMap(attrMap)
		histogramOption = metricManager.GetMetricOptionForRequestDurationByMap(attrMap)
	)
	metricManager.HttpClientRequestActive.Dec(
		ctx,
		requestOption,
	)
	metricManager.HttpClientRequestTotal.Inc(
		ctx,
		responseOption,
	)
	metricManager.HttpClientRequestDuration.Record(
		duration,
		histogramOption,
	)
	metricManager.HttpClientRequestDurationTotal.Add(
		ctx,
		duration,
		responseOption,
	)
	if r.Response != nil {
		metricManager.HttpClientResponseBodySize.Add(
			ctx,
			float64(r.Response.ContentLength),
			responseOption,
		)
	}
}
