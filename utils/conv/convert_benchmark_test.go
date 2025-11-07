package conv

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func Benchmark_ConvertToPrimitive(b *testing.B) {
	typ := []reflect.Type{
		reflect.TypeOf(time.Duration(0)),
		reflect.TypeOf(true),
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(""),
		reflect.TypeOf([]int{}),
	}
	data := [][]string{
		{"10s"},
		{"true"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10"},
		{"10.1"},
		{"10.1"},
		{"string"},
		{"10"},
	}

	for i := 0; i < b.N; i++ {
		idx := i % len(typ)
		ConvertTo(context.Background(), typ[idx], data[idx])
	}
}

func Benchmark_ConvertToSlice(b *testing.B) {
	typ := []reflect.Type{
		reflect.TypeOf([]int{}),
	}
	data := [][]string{
		{"10"},
	}

	for i := 0; i < b.N; i++ {
		idx := i % len(typ)
		ConvertTo(context.Background(), typ[idx], data[idx])
	}
}
