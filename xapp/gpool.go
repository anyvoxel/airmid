package xapp

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
	"go.opentelemetry.io/otel"
	api "go.opentelemetry.io/otel/metric"

	"github.com/anyvoxel/airmid/logger"
	"github.com/anyvoxel/airmid/utils"
)

// GPoolFactory will build the gpool.
type GPoolFactory struct {
	size             int           `airmid:"value:${airmid.gpool.size:=10000}"`
	expiryDuration   time.Duration `airmid:"value:${airmid.gpool.expiry.duration:=10m}"`
	preAlloc         bool          `airmid:"value:${airmid.gpool.prealloc:=false}"`
	maxBlockingTasks int           `airmid:"value:${airmid.gpool.max.blocking.tasks:=0}"`
	nonblocking      bool          `airmid:"value:${airmid.gpool.nonblocking:=true}"`
	disablePurge     bool          `airmid:"value:${airmid.gpool.disable.purge:=false}"`

	panicHandlerProvider gpoolPanicHandlerProvider `airmid:"autowire:?,optional"`
}

type gpoolPanicHandlerProvider interface {
	PoolPanicHandler() func(any)
}

// CreatePool will create the ants.Pool object.
func (f *GPoolFactory) CreatePool() (*ants.Pool, error) {
	opts := []ants.Option{
		ants.WithDisablePurge(f.disablePurge),
		ants.WithExpiryDuration(f.expiryDuration),
		ants.WithMaxBlockingTasks(f.maxBlockingTasks),
		ants.WithNonblocking(f.nonblocking),
		ants.WithPreAlloc(f.preAlloc),
		ants.WithLogger(f),
	}

	if f.panicHandlerProvider != nil {
		opts = append(opts, ants.WithPanicHandler(f.panicHandlerProvider.PoolPanicHandler()))
	}

	gpool, err := ants.NewPool(f.size, opts...)
	if err != nil {
		return nil, err
	}

	err = initGPoolMetrics(gpool)
	if err != nil {
		return nil, err
	}

	return gpool, nil
}

// Printf implement ants.Logger.
func (*GPoolFactory) Printf(format string, args ...any) {
	logger.FromContext(context.TODO()).InfoContext(
		context.TODO(),
		fmt.Sprintf(format, args...),
	)
}

func initGPoolMetrics(p *ants.Pool) error {
	m := otel.Meter(utils.AirmidPackageName, api.WithInstrumentationVersion(utils.AirmidPackageVersion))

	allMetrics := map[string][]api.Float64ObservableGaugeOption{
		"gpool_running_worker": {
			api.WithDescription("It's the number of workers currently running"),
			api.WithFloat64Callback(func(_ context.Context, fo api.Float64Observer) error {
				fo.Observe(float64(p.Running()))
				return nil
			}),
		},
		"gpool_free_worker": {
			api.WithDescription("It's the number of available goroutines to work"),
			api.WithFloat64Callback(func(_ context.Context, fo api.Float64Observer) error {
				fo.Observe(float64(p.Free()))
				return nil
			}),
		},
		"gpool_waiting_worker": {
			api.WithDescription("It's the number of tasks which are waiting be executed"),
			api.WithFloat64Callback(func(_ context.Context, fo api.Float64Observer) error {
				fo.Observe(float64(p.Waiting()))
				return nil
			}),
		},
		"gpool_capacity": {
			api.WithDescription("It's the capacity of pool"),
			api.WithFloat64Callback(func(_ context.Context, fo api.Float64Observer) error {
				fo.Observe(float64(p.Cap()))
				return nil
			}),
		},
	}
	for k, opts := range allMetrics {
		_, err := m.Float64ObservableGauge(
			k,
			opts...,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
