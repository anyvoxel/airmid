package xreflect

import (
	"reflect"
	"testing"
	"unsafe"

	. "github.com/onsi/gomega"
)

func TestIndirectToInterface(t *testing.T) {
	type testCase struct {
		desp   string
		v      any
		expect any
	}
	testCases := []testCase{
		{
			desp:   "normal value",
			v:      int(100),
			expect: int(100),
		},
		{
			desp:   "indirect reflect.Value",
			v:      reflect.ValueOf(int(100)),
			expect: int(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual := IndirectToInterface(tc.v)
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestIndirectToValue(t *testing.T) {
	type testCase struct {
		desp   string
		v      any
		expect any
	}
	testCases := []testCase{
		{
			desp:   "normal value",
			v:      int(100),
			expect: int(100),
		},
		{
			desp:   "indirect reflect.Value",
			v:      reflect.ValueOf(int(100)),
			expect: int(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual := IndirectToValue(tc.v)
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

func TestIndirectToSetableValue(t *testing.T) {
	type testStruct struct {
		V *int
	}

	type testCase struct {
		desp   string
		v      any
		err    string
		expect any
	}
	testCases := []testCase{
		{
			desp: "normal value",
			v: func() any {
				var i int
				return &i
			}(),
			expect: func() any {
				var i int
				return i
			}(),
		},
		{
			desp: "nil pointer",
			v: func() any {
				st := &testStruct{}
				return reflect.ValueOf(st).Elem().Field(0)
			}(),
			expect: func() any {
				var i int
				return i
			}(),
		},
		{
			desp: "cannot set",
			v:    int(0),
			err:  "cannot been set, it must setable",
		},
		{
			desp: "cannot set pointer",
			v:    (*int)(nil),
			err:  "cannot been set, it must setable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := IndirectToSetableValue(tc.v)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

func TestNewValue(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		expect any
		isErr  bool
	}

	vint := int(0)
	vtestCase := testCase{}
	pvint := &vint
	pvtestCase := &vtestCase
	testCases := []testCase{
		{
			desp:   "normal primitive type",
			typ:    reflect.TypeOf(int(1)),
			expect: int(0),
			isErr:  false,
		},
		{
			desp:   "normal struct type",
			typ:    reflect.TypeOf(testCase{}),
			expect: testCase{},
			isErr:  false,
		},
		{
			desp:   "normal ptr to primitive type",
			typ:    reflect.TypeOf(&vint),
			expect: &vint,
			isErr:  false,
		},
		{
			desp:   "normal ptr to struct type",
			typ:    reflect.TypeOf(&vtestCase),
			expect: &vtestCase,
			isErr:  false,
		},
		{
			desp:   "err ptr to primitive type",
			typ:    reflect.TypeOf(&pvint),
			expect: nil,
			isErr:  true,
		},
		{
			desp:   "normal ptr to struct type",
			typ:    reflect.TypeOf(&pvtestCase),
			expect: nil,
			isErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := NewValue(tc.typ)
			if tc.isErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

func TestIsPrimitiveType(t *testing.T) {
	type testCase struct {
		desp string
		typ  reflect.Type
		err  string
	}
	testCases := []testCase{
		{
			desp: "bool is primitive type",
			typ:  reflect.TypeOf(bool(true)),
			err:  "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := IsPrimitiveType(tc.typ)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestReflect(t *testing.T) {
	var p int
	g := NewWithT(t)

	// test for CanSet & IsNil for pointer to value
	v := reflect.ValueOf(&p)
	g.Expect(v.CanSet()).To(BeFalse(), "Pointer to int CanSet")
	g.Expect(v.IsNil()).To(BeFalse(), "Pointer to int IsNil")

	// test for CanSet & IsNil for nil pointer
	v = reflect.ValueOf((*int)(nil))
	g.Expect(v.CanSet()).To(BeFalse(), "Nil pointer to int CanSet")
	g.Expect(v.IsNil()).To(BeTrue(), "Nil pointer to int IsNil")

	type testCase struct {
		V  *int
		V1 int
	}
	st := &testCase{}
	v = reflect.ValueOf(st).Elem().Field(0)
	g.Expect(v.CanSet()).To(BeTrue(), "Struct field pointer CanSet")
	g.Expect(v.IsNil()).To(BeTrue(), "Struct field pointer IsNil")

	// nil pointer cannot set elem
	g.Expect(func() {
		v.Elem().Set(reflect.ValueOf(int(1)))
	}).To(PanicWith(&reflect.ValueError{
		Method: "reflect.Value.Set",
		Kind:   0,
	}))

	// non-nil pointer can set elem
	v.Set(reflect.New(v.Type().Elem()))
	g.Expect(func() {
		v.Elem().Set(reflect.ValueOf(int(1)))
	}).ToNot(PanicWith(&reflect.ValueError{
		Method: "reflect.Value.Set",
		Kind:   0,
	}))
	g.Expect(*st.V).To(Equal(1))

	type testCase2 struct {
		V any
	}
	st2 := &testCase2{}
	v = reflect.ValueOf(st2).Elem().Field(0)
	g.Expect(v.CanSet()).To(BeTrue(), "Struct field interface CanSet")
	g.Expect(v.IsNil()).To(BeTrue(), "Struct field interface IsNil")
	v.Set(reflect.ValueOf(int(1)))
	g.Expect(st2.V).To(Equal(int(1)))

	var pp = 10
	v = reflect.ValueOf(&pp)
	g.Expect(v.Elem().CanAddr()).To(BeTrue())

	v = reflect.NewAt(reflect.TypeOf(int(0)), unsafe.Pointer(v.Elem().UnsafeAddr()))
	v.Elem().SetInt(10)
	g.Expect(pp).To(Equal(10))

	v = reflect.New(reflect.TypeOf(int(0)))
	v.Elem().SetInt(10)

	v = reflect.New(reflect.TypeOf(testCase{}))
	v1 := reflect.NewAt(v.Elem().Type(), unsafe.Pointer(v.Elem().UnsafeAddr()))
	v2 := reflect.NewAt(v.Elem().Type(), unsafe.Pointer(v.Elem().UnsafeAddr()))

	v1.Elem().Field(1).SetInt(10)
	g.Expect(v2.Elem().Field(1).Interface()).To(Equal(int(10)))
}
