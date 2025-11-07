package beans

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testBeanNameAware struct {
	name string
}

func (a *testBeanNameAware) SetBeanName(name string) {
	a.name = name
}

type testBeanFactoryAware struct {
	beanFactory BeanFactory
}

func (a *testBeanFactoryAware) SetBeanFactory(beanFactory BeanFactory) {
	a.beanFactory = beanFactory
}

func TestBeanAwarePostProcessBeforeInitialization(t *testing.T) {
	type testCase struct {
		desp   string
		obj    any
		name   string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "not aware",
			obj:    "1",
			name:   "n1",
			expect: "1",
		},
		{
			desp:   "bean name aware",
			obj:    &testBeanNameAware{name: "1"},
			name:   "n1",
			expect: &testBeanNameAware{name: "n1"},
		},
		{
			desp: "bean factory aware",
			obj:  &testBeanFactoryAware{},
			name: "n1",
			expect: &testBeanFactoryAware{
				beanFactory: NewBeanFactory(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := &BeanAwareProcessor{
				beanFactory: NewBeanFactory(),
			}
			_, err := p.PostProcessBeforeInitialization(tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}

func TestBeanAwarePostProcessAfterInitialization(t *testing.T) {
	type testCase struct {
		desp   string
		obj    any
		name   string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "not aware",
			obj:    "1",
			name:   "n1",
			expect: "1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := &BeanAwareProcessor{}
			_, err := p.PostProcessBeforeInitialization(tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}
