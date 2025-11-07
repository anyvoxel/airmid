package xapp

import (
	"testing"
)

func Benchmark_Shutdown_Parallel(b *testing.B) {
	msg := ""
	m := NewSignalShutdownManager([]ShutdownHandler{
		func(s string) {
			msg += s
		},
	})
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.Shutdown("1")
		}
	})
}
