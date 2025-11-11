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
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	api "go.opentelemetry.io/otel/metric"

	"github.com/anyvoxel/airmid/anvil"
)

// NOTE: how can we compare the benchmark:
// 1. `go test -benchmem -run=^$ -bench ^Benchmark_Global_Meter$ github.com/anyvoxel/airmid/app  -count=10 | tee meter.txt`
// 2. run the command for Benchmark_Global_Meter_Once, output to meter_once.txt
// 3. modify the meter_once.txt, change `Benchmark_Global_Meter_Once` to `Benchmark_Global_Meter`
//    3.1 NOTE: use this to ensure output's function have same name
// 4. benchstat meter.txt meter_once.txt
//    4.1 go install golang.org/x/perf/cmd/benchstat@latest
// According to the benchstat, Benchmark_Global_Meter_Once is 99% faster than Benchmark_Global_Meter

func Benchmark_Global_Meter(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			startTime, err := otel.Meter(anvil.AirmidPackageName, api.WithInstrumentationVersion(anvil.AirmidPackageVersion)).Float64UpDownCounter(
				"start_time_seconds",
				api.WithDescription("Start time of the service in unix timestamp"),
			)
			if err != nil {
				panic(err)
			}

			startTime.Add(context.Background(), 1)
		}
	})
}

var (
	startTimeOnce sync.Once
	startTime     api.Float64UpDownCounter
)

func emitStartTime() {
	startTimeOnce.Do(func() {
		v, err := otel.Meter(anvil.AirmidPackageName, api.WithInstrumentationVersion(anvil.AirmidPackageVersion)).Float64UpDownCounter(
			"start_time_seconds",
			api.WithDescription("Start time of the service in unix timestamp"),
		)
		if err != nil {
			panic(err)
		}

		startTime = v
	})

	startTime.Add(context.Background(), 1)
}

func Benchmark_Global_Meter_Once(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			emitStartTime()
		}
	})
}
