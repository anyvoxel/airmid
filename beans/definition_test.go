package beans

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/utils/pointer"
)

type testAutowireBean struct {
}

// nolint
// testStruct0 is the valid struct with most struct tag
type testStruct0 struct {
	f1 int               // no tag, doesn't generate fd
	f2 int               `airmid:"value:${f2v}"`       // value from property
	F3 string            `airmid:"value:${f3v:=f3vd}"` // value from property with default
	f4 *testAutowireBean `airmid:"autowire:?"`         // autowire by type
	F5 *testAutowireBean `airmid:"autowire:f5b"`       // autowire by bean name
}

// nolint
// testStruct1 for empty airmid tag
type testStruct1 struct {
	f1 int `airmid:""`
}

func TestNewBeanDefinition(t *testing.T) {
	type testCase struct {
		desp   string
		typ    reflect.Type
		opts   []BeanDefinitionOption
		expect BeanDefinition
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			typ:  reflect.TypeOf((*testStruct0)(nil)),
			opts: nil,
			expect: &beanDefinitionHolder{
				Typ:   reflect.TypeOf((*testStruct0)(nil)),
				name:  "testStruct0",
				scope: ScopeSingleton,
				fieldDescriptors: []FieldDescriptor{
					{
						FieldIndex: 1,
						Name:       "f2",
						Typ:        reflect.TypeOf(int(0)),
						Unexported: true,
						Property: &PropertyFieldDescriptor{
							Name:    "f2v",
							Default: nil,
						},
						Bean: nil,
					},
					{
						FieldIndex: 2,
						Name:       "F3",
						Typ:        reflect.TypeOf(""),
						Unexported: false,
						Property: &PropertyFieldDescriptor{
							Name:    "f3v",
							Default: pointer.StringPtr("f3vd"),
						},
						Bean: nil,
					},
					{
						FieldIndex: 3,
						Name:       "f4",
						Typ:        reflect.TypeOf((*testAutowireBean)(nil)),
						Unexported: true,
						Property:   nil,
						Bean: &BeanFieldDescriptor{
							Name: "?",
						},
					},
					{
						FieldIndex: 4,
						Name:       "F5",
						Typ:        reflect.TypeOf((*testAutowireBean)(nil)),
						Unexported: false,
						Property:   nil,
						Bean: &BeanFieldDescriptor{
							Name: "f5b",
						},
					},
				},
			},
			err: "",
		},
		{
			desp: "explicit bean name",
			typ:  reflect.TypeOf((*testStruct0)(nil)),
			opts: []BeanDefinitionOption{
				WithBeanName("Test0"),
			},
			expect: &beanDefinitionHolder{
				Typ:   reflect.TypeOf((*testStruct0)(nil)),
				name:  "Test0",
				scope: ScopeSingleton,
				fieldDescriptors: []FieldDescriptor{
					{
						FieldIndex: 1,
						Name:       "f2",
						Typ:        reflect.TypeOf(int(0)),
						Unexported: true,
						Property: &PropertyFieldDescriptor{
							Name:    "f2v",
							Default: nil,
						},
						Bean: nil,
					},
					{
						FieldIndex: 2,
						Name:       "F3",
						Typ:        reflect.TypeOf(""),
						Unexported: false,
						Property: &PropertyFieldDescriptor{
							Name:    "f3v",
							Default: pointer.StringPtr("f3vd"),
						},
						Bean: nil,
					},
					{
						FieldIndex: 3,
						Name:       "f4",
						Typ:        reflect.TypeOf((*testAutowireBean)(nil)),
						Unexported: true,
						Property:   nil,
						Bean: &BeanFieldDescriptor{
							Name: "?",
						},
					},
					{
						FieldIndex: 4,
						Name:       "F5",
						Typ:        reflect.TypeOf((*testAutowireBean)(nil)),
						Unexported: false,
						Property:   nil,
						Bean: &BeanFieldDescriptor{
							Name: "f5b",
						},
					},
				},
			},
			err: "",
		},
		{
			desp:   "not ptr",
			typ:    reflect.TypeOf(int(0)),
			opts:   nil,
			expect: nil,
			err:    "Cannot build bean definition from 'int'",
		},
		{
			desp:   "not ptr to struct",
			typ:    reflect.TypeOf(testCase{}),
			opts:   nil,
			expect: nil,
			err:    "Cannot build bean definition from 'beans.testCase'",
		},
		{
			desp:   "empty airmid tag",
			typ:    reflect.TypeOf((*testStruct1)(nil)),
			opts:   nil,
			expect: nil,
			err:    "Invalid tag ''",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := NewBeanDefinition(tc.typ, tc.opts...)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			actual.(*beanDefinitionHolder).constructor = nil

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
