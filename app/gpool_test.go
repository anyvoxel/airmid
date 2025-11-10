// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package app

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestGPoolFactoryCreatePool(t *testing.T) {
	g := NewWithT(t)

	f := &GPoolFactory{
		size:             10000,
		expiryDuration:   10 * time.Minute,
		preAlloc:         false,
		maxBlockingTasks: 0,
		nonblocking:      true,
		disablePurge:     false,
	}

	reader := metric.NewManualReader()
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(reader)))

	p, err := f.CreatePool()
	f.Printf("pool: %p, err: %v", p, err)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(p).ToNot(BeNil())

	metrics := metricdata.ResourceMetrics{}
	err = reader.Collect(context.TODO(), &metrics)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(metrics.ScopeMetrics).To(HaveLen(1))
	g.Expect(metrics.ScopeMetrics[0].Metrics).To(HaveLen(4))
}
