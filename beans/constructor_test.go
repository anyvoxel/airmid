package beans

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/xerrors"
)

func TestReflectConstructor(t *testing.T) {
	g := NewWithT(t)
	actual, err := (&ReflectConstructor{typ: reflect.TypeOf(t)}).NewObject(nil)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(actual.Interface()).To(Equal(&testing.T{}))
}

type testMethodConstructorStruct struct{}

func (t *testMethodConstructorStruct) F1() int {
	return 1
}

func (t *testMethodConstructorStruct) F2WithError() (int, error) {
	return 2, xerrors.ErrContinue
}

func (t *testMethodConstructorStruct) F2WithNilError() (int, error) {
	return 3, nil
}

func (t *testMethodConstructorStruct) F3WithArgs(x int, y int) int {
	return x + y
}

type testFnConstructorArgumentResolver struct {
	fn func(args []ConstructorArgument) ([]reflect.Value, error)
}

func (t *testFnConstructorArgumentResolver) Resolve(args []ConstructorArgument) ([]reflect.Value, error) {
	return t.fn(args)
}

func TestMethodConstructor(t *testing.T) {
	type testCase struct {
		desp     string
		typ      reflect.Type
		method   reflect.Method
		resolver ConstructorArgumentResolver

		err    string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "output 1",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(0),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "",
			expect: 1,
		},
		{
			desp:   "output 2 with error",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(1),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "Continue",
			expect: 2,
		},
		{
			desp:   "output 2 with nil error",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(2),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return nil, nil
				},
			},
			err:    "",
			expect: 3,
		},
		{
			desp:   "output 5 with args",
			typ:    reflect.TypeOf((*testMethodConstructorStruct)(nil)),
			method: reflect.TypeOf((*testMethodConstructorStruct)(nil)).Method(3),
			resolver: &testFnConstructorArgumentResolver{
				fn: func(args []ConstructorArgument) ([]reflect.Value, error) {
					return []reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)}, nil
				},
			},
			err:    "",
			expect: 5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			m := &MethodConstructor{
				typ:    tc.typ,
				method: tc.method,
			}
			actual, err := m.NewObject(tc.resolver)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
			} else {
				g.Expect(err).ToNot(HaveOccurred())
			}

			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}
