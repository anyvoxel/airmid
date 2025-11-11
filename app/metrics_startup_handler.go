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
	"log/slog"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/ioc"
	slogctx "github.com/veqryn/slog-context"
)

type metricsStartupHandlerConfiguration struct {
	initializers []MetricsInitializer `airmid:"autowire:?"`
}

// nolint
// MetricsIntializer will provide the option for initialize default metrics provider
type MetricsInitializer interface {
	GetOptions() []metric.Option
}

type metricsStartupHandler struct {
	startupTime time.Time
}

var (
	_ ApplicationStartupHandler = (*metricsStartupHandler)(nil)
)

func (*metricsStartupHandler) Name() string {
	return "MetricsStartupHandler"
}

func (m *metricsStartupHandler) BeforeLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	m.startupTime = time.Now()

	beanDefinitions := map[string]ioc.BeanDefinition{
		"airmidMetricsStartupHandlerConfiguration": ioc.MustNewBeanDefinition(
			reflect.TypeOf((*metricsStartupHandlerConfiguration)(nil)),
			ioc.WithLazyMode(),
		),
	}

	for name, beanDefinition := range beanDefinitions {
		err := app.RegisterBeanDefinition(name, beanDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*metricsStartupHandler) AfterLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	metricsC, err := ioc.GetBean[*metricsStartupHandlerConfiguration](app, "airmidMetricsStartupHandlerConfiguration")
	if err != nil {
		return err
	}

	opts := make([]metric.Option, 0)
	for _, i := range metricsC.initializers {
		opts = append(opts, i.GetOptions()...)
	}
	provider := metric.NewMeterProvider(
		opts...,
	)
	otel.SetMeterProvider(provider)
	return nil
}

func convertOptionToAttributes(opt *option) []attribute.KeyValue {
	v := make([]attribute.KeyValue, 0)

	if opt == nil {
		return v
	}

	for _, attr := range opt.attrs {
		v = append(v, attribute.Key(attr.Key).String(attr.Value))
	}
	return v
}

func (m *metricsStartupHandler) BeforeStartRunner(ctx context.Context, _ *airmidApplication, opt *option) error {
	startTime, err := otel.Meter(
		anvil.AirmidPackageName,
		api.WithInstrumentationVersion(anvil.AirmidPackageVersion),
	).Float64UpDownCounter(
		"start_time_seconds",
		api.WithDescription("Start time of the service in unix timestamp"),
	)
	if err != nil {
		slogctx.FromCtx(ctx).ErrorContext(
			ctx, "Cann't initialize start_time_seconds metrics",
			slog.Any("Error", err),
		)
		return nil
	}

	startTime.Add(ctx, float64(m.startupTime.Unix()), api.WithAttributes(convertOptionToAttributes(opt)...))
	return nil
}
