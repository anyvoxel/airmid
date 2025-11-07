package xtime

import (
	"testing"
)

func Benchmark_SystemTimer_Mills(b *testing.B) {
	tt := NewSystemTimer()
	for i := 0; i < b.N; i++ {
		tt.CurrentTimeMills()
	}
}

func Benchmark_SystemTimer_Mills_Parallel(b *testing.B) {
	tt := NewSystemTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.CurrentTimeMills()
		}
	})
}
