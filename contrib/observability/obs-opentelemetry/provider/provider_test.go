// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func Test_newResource(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config
	}
	tests := []struct {
		name              string
		args              args
		wantResources     []attribute.KeyValue
		unwantedResources []attribute.KeyValue
	}{
		{
			name: "with conflict schema version",
			args: args{
				ctx: context.Background(),
				cfg: &config{
					resourceAttributes: []attribute.KeyValue{
						semconv.ServiceNameKey.String("test-semconv-resource"),
					},
				},
			},
			wantResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("test-semconv-resource"),
			},
			unwantedResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("unknown_service:___Test_newResource_in_github_com_hertz_contrib_obs_opentelemetry_provider.test"),
			},
		},
		{
			name: "resource override",
			args: args{
				ctx: context.Background(),
				cfg: &config{
					resource: resource.Default(),
					resourceAttributes: []attribute.KeyValue{
						semconv.ServiceNameKey.String("test-resource"),
					},
				},
			},
			wantResources: nil,
			unwantedResources: []attribute.KeyValue{
				semconv.ServiceNameKey.String("test-resource"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newResource(tt.args.ctx, tt.args.cfg)
			for _, res := range tt.wantResources {
				assert.Contains(t, got.Attributes(), res)
			}
			for _, unwantedResource := range tt.unwantedResources {
				assert.NotContains(t, got.Attributes(), unwantedResource)
			}
		})
	}
}
