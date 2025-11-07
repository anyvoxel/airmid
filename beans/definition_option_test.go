package beans

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidate(t *testing.T) {
	type testCase struct {
		desp string
		o    *beanDefinitionOption
		err  string
	}
	testCases := []testCase{
		{
			desp: "default option test",
			o:    defaultBeanDefinitionOption(),
			err:  "",
		},
		{
			desp: "singleton scope",
			o: &beanDefinitionOption{
				scope: ScopeSingleton,
			},
			err: "",
		},
		{
			desp: "prototype scope",
			o: &beanDefinitionOption{
				scope: ScopePrototype,
			},
			err: "",
		},
		{
			desp: "invalid scope",
			o: &beanDefinitionOption{
				scope: "1",
			},
			err: "Unsupport scope '1'",
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

func TestWithBeanName(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.name).To(Equal(""))

	o := WithBeanName("1")
	o.Apply(opt)

	g.Expect(opt.name).To(Equal("1"))
}

func TestWithBeanScope(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.scope).To(Equal(ScopeSingleton))

	o := WithBeanScope(ScopePrototype)
	o.Apply(opt)

	g.Expect(opt.scope).To(Equal(ScopePrototype))
}

func TestWithLazyMode(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.lazy).To(BeFalse())

	o := WithLazyMode()
	o.Apply(opt)

	g.Expect(opt.lazy).To(BeTrue())
}

func TestWithPrimary(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt.primary).To(BeFalse())

	o := WithPrimary()
	o.Apply(opt)

	g.Expect(opt.primary).To(BeTrue())
}
