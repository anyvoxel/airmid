package xapp

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
