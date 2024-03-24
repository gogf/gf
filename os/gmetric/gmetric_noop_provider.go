// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "context"

// noopObservableMetric is an implementer for interface Provider with no truly operations.
type noopProvider struct{}

func newNoopProvider() Provider {
	return &noopProvider{}
}

func (*noopProvider) SetAsGlobal() {}

func (*noopProvider) MeterPerformer(option MeterOption) MeterPerformer {
	return newNoopMeterPerformer(option)
}

func (*noopProvider) ForceFlush(ctx context.Context) error {
	return nil
}

func (*noopProvider) Shutdown(ctx context.Context) error {
	return nil
}
