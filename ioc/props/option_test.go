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

package props

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidate(t *testing.T) {
	type testCase struct {
		desp string
		o    *getOption
		err  string
	}
	testCases := []testCase{
		{
			desp: "default option test",
			o:    defaultGetOption(),
			err:  "",
		},
		{
			desp: "nil typ",
			o: &getOption{
				Target:  reflect.Value{},
				Typ:     nil,
				Default: nil,
			},
			err: "typ cannot be nil",
		},
		{
			desp: "target and typ mismatch",
			o: &getOption{
				Target: reflect.ValueOf(int(1)),
				Typ:    reflect.TypeOf(""),
			},
			err: `target\('int'\) doesn't match typ\('string'\)`,
		},
		{
			desp: "target and typ match",
			o: &getOption{
				Target: reflect.ValueOf(int(1)),
				Typ:    reflect.TypeOf(int(1)),
			},
			err: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.o.Validate()
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestComplete(t *testing.T) {
	type testCase struct {
		desp string
		o    *getOption
		err  string
	}
	testCases := []testCase{
		{
			desp: "default option test",
			o:    defaultGetOption(),
			err:  "",
		},
		{
			desp: "validate nil type",
			o: &getOption{
				Target:  reflect.Value{},
				Typ:     nil,
				Default: nil,
			},
			err: "typ cannot be nil",
		},
		{
			desp: "auto target to int",
			o: &getOption{
				Target:  reflect.Value{},
				Typ:     reflect.TypeOf(int(0)),
				Default: nil,
			},
			err: "",
		},
		{
			desp: "auto target to *string",
			o: &getOption{
				Target:  reflect.Value{},
				Typ:     reflect.TypeOf((*int)(nil)),
				Default: nil,
			},
			err: "",
		},
		{
			desp: "not auto target",
			o: &getOption{
				Target: func() reflect.Value {
					i := int(1)
					return reflect.ValueOf(&i)
				}(),
				Typ:     reflect.TypeOf((*int)(nil)),
				Default: nil,
			},
			err: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.o.Complete()

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.o.Validate()).ToNot(HaveOccurred())
			g.Expect(tc.o.Target.Type()).To(Equal(tc.o.Typ))
		})
	}
}

func TestWithDefault(t *testing.T) {
	g := NewWithT(t)

	o := defaultGetOption()
	g.Expect(o.Default).To(Equal((*string)(nil)))

	WithDefault("1").Apply(o)
	g.Expect(*o.Default).To(Equal("1"))
}

func TestWithTarget(t *testing.T) {
	g := NewWithT(t)

	o := defaultGetOption()
	g.Expect(o.IsTargetValid()).To(BeFalse())

	i := "1"
	WithTarget(&i).Apply(o)
	g.Expect(o.Target.Elem().String()).To(Equal("1"))
	g.Expect(o.Typ).To(Equal(reflect.TypeOf((*string)(nil))))

	o.Target.Elem().Set(reflect.ValueOf("2"))
	g.Expect(o.Target.Elem().String()).To(Equal("2"))
	g.Expect(i).To(Equal("2"))
}

func TestWithType(t *testing.T) {
	g := NewWithT(t)

	o := defaultGetOption()
	g.Expect(o.Typ).To(Equal(reflect.TypeOf("")))

	WithType(reflect.TypeOf(int(1))).Apply(o)
	g.Expect(o.Typ).To(Equal(reflect.TypeOf(int(1))))
}
