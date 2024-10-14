# OpenTelemetry Protocol (stable)

Since v1.35, the Jaeger backend can receive trace data from the OpenTelemetry SDKs in their native OpenTelemetry Protocol (OTLP)  . It is no longer necessary to configure the OpenTelemetry SDKs with Jaeger exporters, nor deploy the OpenTelemetry Collector between the OpenTelemetry SDKs and the Jaeger backend.

The OTLP data is accepted in these formats: (1) binary gRPC, (2) Protobuf over HTTP, (3) JSON over HTTP. For more details on the OTLP receiver see the official documentation  . Note that not all configuration options are supported in jaeger-collector (see --collector.otlp.* CLI Flags ), and only tracing data is accepted, since Jaeger does not store other telemetry types.

|Port	|Protocol	|Endpoint	|Format|
|---|---|---|---|
|4317	|gRPC	|n/a	|Protobuf|
|4318	|HTTP	|/v1/traces	|Protobuf or JSON|

Two ports, 4317 and 4318, have also been added to the export list, which are used by OTLP receivers to listen for gRPC and HTTP connections.

## All in One
all-in-one is an executable designed for quick local testing. It includes the Jaeger UI, jaeger-collector, jaeger-query, and jaeger-agent, with an in memory storage component.

The simplest way to start the all-in-one is to use the pre-built image published to DockerHub (a single command line).

```shell
docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.55

```

Or run the `jaeger-all-in-one(.exe)` executable from the [binary distribution archives][download]:

```bash
jaeger-all-in-one --collector.zipkin.host-port=:9411
```

You can then navigate to `http://localhost:16686` to access the Jaeger UI.

The container exposes the following ports:

Port  | Protocol | Component | Function
----- | -------  | --------- | ---
6831  | UDP      | agent     | accept `jaeger.thrift` over Thrift-compact protocol (used by most SDKs)
6832  | UDP      | agent     | accept `jaeger.thrift` over Thrift-binary protocol (used by Node.js SDK)
5775  | UDP      | agent     | (deprecated) accept `zipkin.thrift` over compact Thrift protocol (used by legacy clients only)
5778  | HTTP     | agent     | serve configs (sampling, etc.)
16686 | HTTP     | query     | serve frontend
4317  | HTTP     | collector | accept OpenTelemetry Protocol (OTLP) over gRPC
4318  | HTTP     | collector | accept OpenTelemetry Protocol (OTLP) over HTTP
14268 | HTTP     | collector | accept `jaeger.thrift` directly from clients
14250 | HTTP     | collector | accept `model.proto`
9411  | HTTP     | collector | Zipkin compatible endpoint (optional)

https://www.jaegertracing.io/docs/1.55/apis/#opentelemetry-protocol-stable

Read more details:

Introducing native support for OpenTelemetry in Jaeger

https://medium.com/jaegertracing/introducing-native-support-for-opentelemetry-in-jaeger-eb661be8183c

Better alignment with OpenTelemetry by focusing on OTLP

https://medium.com/jaegertracing/better-alignment-with-opentelemetry-by-focusing-on-otlp-f3688939073f



# Jaeger

Open Source, End-to-End Distributed Tracing. Hosted by Cloud Native Computing Foundation.

https://medium.com/jaegertracing


## For more information on distributed tracing initialization, please refer to the content in the [provider](./provider) directory