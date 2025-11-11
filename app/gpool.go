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
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
	slogctx "github.com/veqryn/slog-context"
	"go.opentelemetry.io/otel"
	api "go.opentelemetry.io/otel/metric"

	"github.com/anyvoxel/airmid/anvil"
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
	slogctx.FromCtx(context.TODO()).InfoContext(
		context.TODO(),
		fmt.Sprintf(format, args...),
	)
}

func initGPoolMetrics(p *ants.Pool) error {
	m := otel.Meter(anvil.AirmidPackageName, api.WithInstrumentationVersion(anvil.AirmidPackageVersion))

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
