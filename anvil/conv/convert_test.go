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

	. "github.com/onsi/gomega"
)

func TestConvertToBool(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"true"},
			expect: true,
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type bool",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToBool(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToUint(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: uint(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type uint",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToUint(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToUint8(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: uint8(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type uint8",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToUint8(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToUint16(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: uint16(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type uint16",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToUint16(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToUint32(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: uint32(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type uint32",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToUint32(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToUint64(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: uint64(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type uint64",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToUint64(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToInt(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: int(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type int",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToInt(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToInt8(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: int8(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type int8",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToInt8(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToInt16(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: int16(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type int16",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToInt16(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToInt32(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: int32(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type int32",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToInt32(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToInt64(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1"},
			expect: int64(1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type int64",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToInt64(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToFloat32(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1.1"},
			expect: float32(1.1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type float32",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToFloat32(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToFloat64(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1.1"},
			expect: float64(1.1),
			err:    "",
		},
		{
			desp:   "convert failed",
			data:   []string{"true1"},
			expect: true,
			err:    "Cann't convert true1 to type float64",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToFloat64(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToString(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1.1"},
			expect: "1.1",
			err:    "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToString(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertToTimeDuration(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal test",
			data:   []string{"1s"},
			expect: time.Duration(1) * time.Second,
			err:    "",
		},
		{
			desp:   "miss unit",
			data:   []string{"1"},
			expect: nil,
			err:    "missing unit in duration \"1\"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertToTimeDuration(context.Background(), tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestConvertTo(t *testing.T) {
	type testCase struct {
		desp   string
		data   []string
		typ    reflect.Type
		expect any
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal bool",
			data:   []string{"true"},
			typ:    reflect.TypeOf(true),
			expect: true,
			err:    "",
		},
		{
			desp:   "normal int",
			data:   []string{"1"},
			typ:    reflect.TypeOf(int(0)),
			expect: int(1),
			err:    "",
		},
		{
			desp:   "normal int8",
			data:   []string{"1"},
			typ:    reflect.TypeOf(int8(0)),
			expect: int8(1),
			err:    "",
		},
		{
			desp:   "normal int16",
			data:   []string{"1"},
			typ:    reflect.TypeOf(int16(0)),
			expect: int16(1),
			err:    "",
		},
		{
			desp:   "normal int32",
			data:   []string{"1"},
			typ:    reflect.TypeOf(int32(0)),
			expect: int32(1),
			err:    "",
		},
		{
			desp:   "normal int64",
			data:   []string{"1"},
			typ:    reflect.TypeOf(int64(0)),
			expect: int64(1),
			err:    "",
		},
		{
			desp:   "normal uint",
			data:   []string{"1"},
			typ:    reflect.TypeOf(uint(0)),
			expect: uint(1),
			err:    "",
		},
		{
			desp:   "normal uint8",
			data:   []string{"1"},
			typ:    reflect.TypeOf(uint8(0)),
			expect: uint8(1),
			err:    "",
		},
		{
			desp:   "normal uint16",
			data:   []string{"1"},
			typ:    reflect.TypeOf(uint16(0)),
			expect: uint16(1),
			err:    "",
		},
		{
			desp:   "normal uint32",
			data:   []string{"1"},
			typ:    reflect.TypeOf(uint32(0)),
			expect: uint32(1),
			err:    "",
		},
		{
			desp:   "normal uint64",
			data:   []string{"1"},
			typ:    reflect.TypeOf(uint64(0)),
			expect: uint64(1),
			err:    "",
		},
		{
			desp:   "normal float32",
			data:   []string{"1.1"},
			typ:    reflect.TypeOf(float32(0)),
			expect: float32(1.1),
			err:    "",
		},
		{
			desp:   "normal float64",
			data:   []string{"1.1"},
			typ:    reflect.TypeOf(float64(0)),
			expect: float64(1.1),
			err:    "",
		},
		{
			desp:   "normal string",
			data:   []string{"1.1"},
			typ:    reflect.TypeOf(""),
			expect: "1.1",
			err:    "",
		},
		{
			desp:   "normal time duration",
			data:   []string{"1s"},
			typ:    reflect.TypeOf(time.Duration(0)),
			expect: time.Duration(1) * time.Second,
			err:    "",
		},
		{
			desp:   "normal []int",
			data:   []string{"1"},
			typ:    reflect.TypeOf([]int{}),
			expect: []int{1},
			err:    "",
		},
		{
			desp:   "normal []int failed",
			data:   []string{"true"},
			typ:    reflect.TypeOf([]int{}),
			expect: []int{},
			err:    "Cann't convert true to type int",
		},
		{
			desp:   "unsupport type",
			data:   []string{"true"},
			typ:    reflect.TypeOf(testing.T{}),
			expect: []int{},
			err:    "Unsupport target type testing.T",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ConvertTo(context.Background(), tc.typ, tc.data)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestToString(t *testing.T) {
	type testCase struct {
		desp   string
		i      any
		expect string
		err    string
	}
	testCases := []testCase{
		{
			desp:   "normal string",
			i:      "123i",
			expect: "123i",
			err:    "",
		},
		{
			desp:   "normal bool",
			i:      true,
			expect: "true",
			err:    "",
		},
		{
			desp:   "normal float32",
			i:      float32(1.1),
			expect: "1.1",
			err:    "",
		},
		{
			desp:   "normal float64",
			i:      float64(1.1),
			expect: "1.1",
			err:    "",
		},
		{
			desp:   "normal int",
			i:      int(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal int8",
			i:      int8(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal int16",
			i:      int16(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal int32",
			i:      int32(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal int64",
			i:      int64(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal uint",
			i:      uint(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal uint8",
			i:      uint8(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal uint16",
			i:      uint16(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal uint32",
			i:      uint32(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal uint64",
			i:      uint64(1),
			expect: "1",
			err:    "",
		},
		{
			desp:   "normal time duration",
			i:      time.Duration(1) * time.Second,
			expect: "1s",
			err:    "",
		},
		{
			desp:   "unsupport type",
			i:      testing.T{},
			expect: "",
			err:    "Unsupport target type 'testing.T'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ToString(tc.i)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
