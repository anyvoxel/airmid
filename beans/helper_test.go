package beans

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsTypeMatched(t *testing.T) {
	type testCase struct {
		desp      string
		targetTyp reflect.Type
		beanTyp   reflect.Type
		expect    bool
	}
	testCases := []testCase{
		{
			desp:      "type assignable",
			targetTyp: reflect.TypeOf(testCase{}),
			beanTyp:   reflect.TypeOf(testCase{}),
			expect:    true,
		},
		{
			desp:      "interface assignable",
			targetTyp: reflect.TypeOf(any(testCase{})),
			beanTyp:   reflect.TypeOf(testCase{}),
			expect:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(IsTypeMatched(tc.targetTyp, tc.beanTyp)).To(Equal(tc.expect))
		})
	}
}

type testNamerI interface {
	Name() string
}

type testNamerI2 interface {
	Name2() string
}

type testNamer1 struct {
	n string
}

func (t *testNamer1) Name() string {
	return t.n
}

func (t *testNamer1) Name2() string {
	return t.n + "/2"
}

type proxyNamer struct {
	testNamerI
}

func (p *proxyNamer) OriginalObject() any {
	return p.testNamerI
}

type nonproxyNamer struct {
	testNamerI
}

func TestIndirectTo(t *testing.T) {
	t.Run("directly implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &testNamer1{
			n: "1",
		}
		actual := IndirectTo[testNamerI](v)
		g.Expect(actual).To(Equal(v))
	})

	t.Run("proxy directly implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &proxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI](v)
		g.Expect(actual).To(Equal(v))
	})

	t.Run("proxy non implement", func(t *testing.T) {
		g := NewWithT(t)

		v := &proxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI2](v)
		g.Expect(actual).To(Equal(v.testNamerI))
	})

	t.Run("non-proxy wrapper", func(t *testing.T) {
		g := NewWithT(t)

		v := &nonproxyNamer{
			&testNamer1{
				n: "1",
			},
		}
		actual := IndirectTo[testNamerI2](v)
		g.Expect(actual).To(BeNil())
	})
}

func TestGetBean(t *testing.T) {
	t.Run("didn't exists", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`No bean '1' registered: ObjectNotFound`))
		g.Expect(o).To(BeNil())
	})

	t.Run("type mismatch", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()
		bf.RegisterBeanDefinition("1", MustNewBeanDefinition(reflect.TypeOf((*proxyNamer)(nil))))

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`cannot convert \*beans\.proxyNamer to \*beans\.nonproxyNamer`))
		g.Expect(o).To(BeNil())
	})

	t.Run("normal", func(t *testing.T) {
		g := NewWithT(t)
		bf := NewBeanFactory()
		bf.RegisterBeanDefinition("1", MustNewBeanDefinition(reflect.TypeOf((*nonproxyNamer)(nil))))

		o, err := GetBean[*nonproxyNamer](bf, "1")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(o).ToNot(BeNil())
	})
}
