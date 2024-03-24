// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

type noopMeterPerformer struct{}

func newNoopMeterPerformer(option MeterOption) MeterPerformer {
	return &noopMeterPerformer{}
}

func (*noopMeterPerformer) CounterPerformer(name string, option MetricOption) (CounterPerformer, error) {
	return newNoopCounterPerformer(), nil
}

func (*noopMeterPerformer) UpDownCounterPerformer(name string, option MetricOption) (UpDownCounterPerformer, error) {
	return newNoopUpDownCounterPerformer(), nil
}

func (*noopMeterPerformer) HistogramPerformer(name string, option MetricOption) (HistogramPerformer, error) {
	return newNoopHistogramPerformer(), nil
}

func (*noopMeterPerformer) ObservableCounterPerformer(name string, option MetricOption) (ObservableCounterPerformer, error) {
	return newNoopObservableCounterPerformer(), nil
}

func (*noopMeterPerformer) ObservableUpDownCounterPerformer(name string, option MetricOption) (ObservableUpDownCounterPerformer, error) {
	return newNoopObservableUpDownCounterPerformer(), nil
}

func (*noopMeterPerformer) ObservableGaugePerformer(name string, option MetricOption) (ObservableGaugePerformer, error) {
	return newNoopObservableGaugePerformer(), nil
}

// RegisterCallback registers callback on certain metrics.
// A callback is bound to certain component and version, it is called when the associated metrics are read.
// Multiple callbacks on the same component and version will be called by their registered sequence.
func (*noopMeterPerformer) RegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric) error {
	return nil
}
