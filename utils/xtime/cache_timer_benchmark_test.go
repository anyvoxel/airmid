package xtime

import (
	"testing"
)

func Benchmark_CacheTimer_Mills(b *testing.B) {
	tt := NewCacheTimer()
	for i := 0; i < b.N; i++ {
		tt.CurrentTimeMills()
	}
}

func Benchmark_CacheTimer_Mills_Parallel(b *testing.B) {
	tt := NewCacheTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.CurrentTimeMills()
		}
	})
}
