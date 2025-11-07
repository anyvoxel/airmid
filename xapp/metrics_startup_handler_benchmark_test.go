package xapp

import (
	"context"
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	api "go.opentelemetry.io/otel/metric"

	"github.com/anyvoxel/airmid/utils"
)

// NOTE: how can we compare the benchmark:
// 1. `go test -benchmem -run=^$ -bench ^Benchmark_Global_Meter$ github.com/anyvoxel/airmid/xapp  -count=10 | tee meter.txt`
// 2. run the command for Benchmark_Global_Meter_Once, output to meter_once.txt
// 3. modify the meter_once.txt, change `Benchmark_Global_Meter_Once` to `Benchmark_Global_Meter`
//    3.1 NOTE: use this to ensure output's function have same name
// 4. benchstat meter.txt meter_once.txt
//    4.1 go install golang.org/x/perf/cmd/benchstat@latest
// According to the benchstat, Benchmark_Global_Meter_Once is 99% faster than Benchmark_Global_Meter

func Benchmark_Global_Meter(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			startTime, err := otel.Meter(utils.AirmidPackageName, api.WithInstrumentationVersion(utils.AirmidPackageVersion)).Float64UpDownCounter(
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
		v, err := otel.Meter(utils.AirmidPackageName, api.WithInstrumentationVersion(utils.AirmidPackageVersion)).Float64UpDownCounter(
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
