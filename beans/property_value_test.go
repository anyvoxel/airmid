package beans

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAddValue(t *testing.T) {
	g := NewWithT(t)
	p := NewPropertyValues().(*propertyValuesImpl)

	g.Expect(p.values).To(Equal(make(map[int]reflect.Value)))
	p.AddValue(1, reflect.ValueOf(g))
	g.Expect(p.values).To(Equal(map[int]reflect.Value{
		1: reflect.ValueOf(g),
	}))
}

func TestSetProperty(t *testing.T) {
	g := NewWithT(t)
	p := NewPropertyValues()

	type testBean struct {
		v1 int //nolint
		v2 string
		V3 []int
	}
	p.AddValue(1, reflect.ValueOf("v2"))
	p.AddValue(2, reflect.ValueOf([]int{0, 1}))
	b := &testBean{}
	p.SetProperty(reflect.ValueOf(b).Elem(), []FieldDescriptor{
		{
			FieldIndex: 0,
			Name:       "v1",
			Typ:        reflect.TypeOf(int(0)),
			Unexported: true,
		},
		{
			FieldIndex: 1,
			Name:       "v2",
			Typ:        reflect.TypeOf(""),
			Unexported: true,
		},
		{
			FieldIndex: 2,
			Name:       "V3",
			Typ:        reflect.TypeOf([]int{}),
			Unexported: false,
		},
	})
	g.Expect(b).To(Equal(&testBean{
		v2: "v2",
		V3: []int{0, 1},
	}))
}
