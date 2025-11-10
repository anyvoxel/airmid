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
