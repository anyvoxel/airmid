package xapp

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/logger"
	"github.com/anyvoxel/airmid/utils"
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

	beanDefinitions := map[string]beans.BeanDefinition{
		"airmidMetricsStartupHandlerConfiguration": beans.MustNewBeanDefinition(
			reflect.TypeOf((*metricsStartupHandlerConfiguration)(nil)),
			beans.WithLazyMode(),
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
	metricsC, err := beans.GetBean[*metricsStartupHandlerConfiguration](app, "airmidMetricsStartupHandlerConfiguration")
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
		utils.AirmidPackageName,
		api.WithInstrumentationVersion(utils.AirmidPackageVersion),
	).Float64UpDownCounter(
		"start_time_seconds",
		api.WithDescription("Start time of the service in unix timestamp"),
	)
	if err != nil {
		logger.FromContext(ctx).ErrorContext(
			ctx, "Cann't initialize start_time_seconds metrics",
			slog.Any("Error", err),
		)
		return nil
	}

	startTime.Add(ctx, float64(m.startupTime.Unix()), api.WithAttributes(convertOptionToAttributes(opt)...))
	return nil
}
