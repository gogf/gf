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
	HttpClientRequestActive        gmetric.UpDownCounter
	HttpClientRequestTotal         gmetric.Counter
	HttpClientRequestDuration      gmetric.Histogram
	HttpClientRequestDurationTotal gmetric.Counter
	HttpClientConnectionDuration   gmetric.Histogram
	HttpClientRequestBodySize      gmetric.Counter
	HttpClientResponseBodySize     gmetric.Counter
}

const (
	metricAttrKeyServerAddress          = "server.address"
	metricAttrKeyServerPort             = "server.port"
	metricAttrKeyUrlSchema              = "url.schema"
	metricAttrKeyHttpRequestMethod      = "http.request.method"
	metricAttrKeyHttpResponseStatusCode = "http.response.status_code"
	metricAttrKeyNetworkProtocolVersion = "network.protocol.version"
)

var (
	// metricManager for http client metrics.
	metricManager = newMetricManager()
)

func newMetricManager() *localMetricManager {
	meter := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        instrumentName,
		InstrumentVersion: gf.VERSION,
	})
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
		HttpClientRequestDuration: meter.MustHistogram(
			"http.client.request.duration",
			gmetric.MetricOption{
				Help:       "Measures the duration of client requests.",
				Unit:       "ms",
				Attributes: gmetric.Attributes{},
				Buckets:    durationBuckets,
			},
		),
		HttpClientRequestTotal: meter.MustCounter(
			"http.client.request.total",
			gmetric.MetricOption{
				Help:       "Total processed request number.",
				Unit:       "",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpClientRequestActive: meter.MustUpDownCounter(
			"http.client.request.active",
			gmetric.MetricOption{
				Help:       "Number of active client requests.",
				Unit:       "",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpClientRequestDurationTotal: meter.MustCounter(
			"http.client.request.duration_total",
			gmetric.MetricOption{
				Help:       "Total execution duration of request.",
				Unit:       "ms",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpClientRequestBodySize: meter.MustCounter(
			"http.client.request.body_size",
			gmetric.MetricOption{
				Help:       "Outgoing request bytes total.",
				Unit:       "bytes",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpClientResponseBodySize: meter.MustCounter(
			"http.client.response.body_size",
			gmetric.MetricOption{
				Help:       "Response bytes total.",
				Unit:       "bytes",
				Attributes: gmetric.Attributes{},
			},
		),
		HttpClientConnectionDuration: meter.MustHistogram(
			"http.client.connection_duration",
			gmetric.MetricOption{
				Help:       "Measures the connection establish duration of client requests.",
				Unit:       "ms",
				Attributes: gmetric.Attributes{},
				Buckets:    durationBuckets,
			},
		),
	}
	return mm
}

func (m *localMetricManager) GetMetricOptionForHistogram(r *http.Request) gmetric.Option {
	attrMap := m.GetMetricAttributeMap(r)
	return gmetric.Option{
		Attributes: attrMap.Pick(
			metricAttrKeyServerAddress,
			metricAttrKeyServerPort,
		),
	}
}

func (m *localMetricManager) GetMetricOptionForHistogramByMap(attrMap gmetric.AttributeMap) gmetric.Option {
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
		ctx             = r.Context()
		attrMap         = metricManager.GetMetricAttributeMap(r)
		requestOption   = metricManager.GetMetricOptionForRequestByMap(attrMap)
		requestBodySize = float64(r.ContentLength)
	)
	metricManager.HttpClientRequestActive.Inc(
		ctx,
		requestOption,
	)
	if requestBodySize > 0 {
		metricManager.HttpClientRequestBodySize.Add(
			ctx,
			requestBodySize,
			requestOption,
		)
	}
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
		histogramOption = metricManager.GetMetricOptionForHistogramByMap(attrMap)
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
		var responseBodySize = float64(r.Response.ContentLength)
		if responseBodySize > 0 {
			metricManager.HttpClientResponseBodySize.Add(
				ctx,
				responseBodySize,
				responseOption,
			)
		}
	}
}
